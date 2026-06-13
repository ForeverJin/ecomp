package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
)

func (a *App) handleLCSCForm(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "立创商城订单导入",
		ActivePage: "lcsc",
	}
	a.render(w, r, "lcsc_import", data)
}

func (a *App) handleLCSCUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		setFlash(w, "danger", "文件上传失败: "+err.Error())
		http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		setFlash(w, "danger", "请选择文件")
		http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	var items []LCSCItem

	switch ext {
	case ".xls":
		// Read file into memory for xls library
		data, err := io.ReadAll(file)
		if err != nil {
			setFlash(w, "danger", "读取文件失败")
			http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
			return
		}
		items, err = parseXLS(data)
		if err != nil {
			setFlash(w, "danger", "解析XLS失败: "+err.Error())
			http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
			return
		}
	case ".xlsx":
		items, err = parseXLSX(file)
		if err != nil {
			setFlash(w, "danger", "解析XLSX失败: "+err.Error())
			http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
			return
		}
	default:
		setFlash(w, "danger", "不支持的文件格式，请上传 .xls 或 .xlsx 文件")
		http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
		return
	}

	if len(items) == 0 {
		setFlash(w, "danger", "未找到商品明细，请确认文件格式正确")
		http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
		return
	}

	// Match against existing components
	for i := range items {
		comp := a.matchComponent(&items[i])
		if comp != nil {
			items[i].MatchedID = comp.ID
			items[i].MatchType = "existing"
			items[i].CurrentQty = comp.Quantity
		} else {
			items[i].MatchType = "new"
		}
	}

	// Encode items for preview form (signed hidden field)
	itemsJSON, _ := json.Marshal(items)
	itemsB64 := base64.StdEncoding.EncodeToString(itemsJSON)

	data := &TemplateData{
		Title:      "确认导入 - 立创订单",
		ActivePage: "lcsc",
		Data: map[string]interface{}{
			"Items":         items,
			"ItemsData":     itemsB64,
			"TotalItems":    len(items),
			"ExistingCount": countExisting(items),
			"NewCount":      countNew(items),
		},
	}
	a.render(w, r, "lcsc_preview", data)
}

func (a *App) handleLCSCConfirm(w http.ResponseWriter, r *http.Request) {
	if !a.verifyCSRF(r) {
		http.Error(w, "CSRF验证失败", http.StatusForbidden)
		return
	}

	// Decode items from hidden field
	itemsB64 := r.FormValue("items_data")
	itemsJSON, err := base64.StdEncoding.DecodeString(itemsB64)
	if err != nil {
		setFlash(w, "danger", "数据解码失败")
		http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
		return
	}

	var items []LCSCItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		setFlash(w, "danger", "数据解析失败")
		http.Redirect(w, r, "/transactions/lcsc", http.StatusFound)
		return
	}

	// Process selected items
	imported := 0
	for _, item := range items {
		// Check if this item was selected
		if r.FormValue("item_"+strconv.Itoa(item.Index)) != "on" {
			continue
		}

		var componentID int
		if item.MatchedID > 0 {
			// Existing component - just update quantity
			componentID = item.MatchedID
		} else {
			// New component - create it
			c := &Component{
				Name:     item.Name,
				Model:    item.Model,
				Package:  item.Package,
				Quantity: 0, // Will be updated by transaction
				Location: "",
				Notes:    "立创编号: " + item.LCSCNumber,
			}
			id, err := a.db.CreateComponent(c)
			if err != nil {
				log.Printf("创建元器件失败 (%s): %v", item.Name, err)
				continue
			}
			componentID = id
		}

		// Create inbound transaction
		tx := &Transaction{
			ComponentID: componentID,
			Type:        "inbound",
			Quantity:    item.Quantity,
			Operator:    "",
			Reason:      "立创订单导入",
			Notes:       "商品编号: " + item.LCSCNumber,
		}
		if err := a.db.CreateTransaction(tx); err != nil {
			log.Printf("创建交易记录失败: %v", err)
			continue
		}

		// Update stock
		if err := a.db.UpdateComponentQuantity(componentID, item.Quantity); err != nil {
			log.Printf("更新库存失败: %v", err)
			continue
		}
		imported++
	}

	setFlash(w, "success", fmt.Sprintf("成功导入 %d 个元器件", imported))
	http.Redirect(w, r, "/components", http.StatusFound)
}

// --- LCSC Excel Parsing ---

func parseXLS(data []byte) ([]LCSCItem, error) {
	reader := bytes.NewReader(data)
	wb, err := xls.OpenReader(reader, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("打开XLS文件失败: %w", err)
	}

	sheet := wb.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("未找到工作表")
	}

	return parseLCSCRows(func(maxRow int) [][]string {
		var rows [][]string
		for i := 0; i <= maxRow; i++ {
			row := sheet.Row(i)
			if row == nil {
				continue
			}
			var cells []string
			// Read up to 20 columns
			for j := 0; j < 20; j++ {
				cells = append(cells, row.Col(j))
			}
			rows = append(rows, cells)
		}
		return rows
	}, int(sheet.MaxRow))
}

func parseXLSX(file io.Reader) ([]LCSCItem, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("打开XLSX文件失败: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %w", err)
	}

	return parseLCSCRows(func(maxRow int) [][]string {
		return rows
	}, len(rows))
}

type rowReader func(maxRow int) [][]string

func parseLCSCRows(readRows rowReader, maxRow int) ([]LCSCItem, error) {
	rows := readRows(maxRow)

	// Find header row containing "序号" and "商品编号"
	headerRow := -1
	colIndex := -1    // 序号
	colLCSC := -1     // 商品编号
	colModel := -1    // 厂家型号
	colPackage := -1  // 封装
	colName := -1     // 商品名称
	colQty := -1      // 订购数量

	for i, row := range rows {
		hasSeq := false
		hasLCSC := false
		for j, cell := range row {
			c := strings.TrimSpace(cell)
			switch {
			case c == "序号":
				hasSeq = true
				colIndex = j
			case c == "商品编号":
				hasLCSC = true
				colLCSC = j
			case c == "厂家型号" || c == "型号":
				colModel = j
			case c == "封装" || c == "封装/规格":
				colPackage = j
			case c == "商品名称":
				colName = j
			case c == "订购数量":
				colQty = j
			}
		}
		if hasSeq && hasLCSC {
			headerRow = i
			break
		}
	}

	if headerRow == -1 {
		return nil, fmt.Errorf("未找到包含\"序号\"和\"商品编号\"的表头行")
	}

	var items []LCSCItem
	for i := headerRow + 1; i < len(rows); i++ {
		row := rows[i]

		// Skip empty rows
		getCell := func(col int) string {
			if col < 0 || col >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[col])
		}

		idx := getCell(colIndex)
		if idx == "" {
			continue
		}

		name := getCell(colName)
		model := getCell(colModel)
		if name == "" && model == "" {
			continue
		}

		// Parse quantity - may have "个" suffix
		qtyStr := getCell(colQty)
		qtyStr = strings.TrimSuffix(qtyStr, "个")
		qtyStr = strings.TrimSpace(qtyStr)
		qty, _ := strconv.Atoi(qtyStr)
		if qty <= 0 {
			// Try to extract number from string
			var n int
			fmt.Sscanf(qtyStr, "%d", &n)
			if n > 0 {
				qty = n
			} else {
				continue
			}
		}

		items = append(items, LCSCItem{
			Index:      len(items),
			LCSCNumber: getCell(colLCSC),
			Name:       name,
			Model:      model,
			Package:    getCell(colPackage),
			Quantity:   qty,
		})
	}

	return items, nil
}

// matchComponent tries to find an existing component matching the LCSC item.
func (a *App) matchComponent(item *LCSCItem) *Component {
	// Priority 1: exact model match
	if item.Model != "" {
		comp, _ := a.db.FindComponentByModel(item.Model)
		if comp != nil {
			return comp
		}
	}
	// Priority 2: name + package match
	if item.Name != "" {
		comp, _ := a.db.FindComponentByNamePackage(item.Name, item.Package)
		if comp != nil {
			return comp
		}
	}
	return nil
}

func countExisting(items []LCSCItem) int {
	n := 0
	for _, item := range items {
		if item.MatchType == "existing" {
			n++
		}
	}
	return n
}

func countNew(items []LCSCItem) int {
	n := 0
	for _, item := range items {
		if item.MatchType == "new" {
			n++
		}
	}
	return n
}

// lcscItemJSON is used by template to render items_data as safe HTML attribute.
func init() {
	// Register additional template function for LCSC
	templateFuncs()["safeAttr"] = func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	}
}
