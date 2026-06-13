package main

import (
	"net/http"
	"strconv"
)

func (a *App) handleInbound(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		componentID := r.URL.Query().Get("component_id")
		components, _ := a.db.ListComponents("")
		if components == nil {
			components = []Component{}
		}

		data := &TemplateData{
			Title:      "入库",
			ActivePage: "inbound",
			Data: map[string]interface{}{
				"Components":     components,
				"SelectedCompID": componentID,
			},
		}
		a.render(w, r, "inbound", data)
		return
	}

	// POST
	if !a.verifyCSRF(r) {
		http.Error(w, "CSRF验证失败", http.StatusForbidden)
		return
	}

	componentID, _ := strconv.Atoi(r.FormValue("component_id"))
	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	operator := r.FormValue("operator")
	reason := r.FormValue("reason")
	notes := r.FormValue("notes")

	if componentID == 0 || qty <= 0 {
		setFlash(w, "danger", "请选择元器件并输入有效数量")
		http.Redirect(w, r, "/transactions/inbound", http.StatusFound)
		return
	}

	tx := &Transaction{
		ComponentID: componentID,
		Type:        "inbound",
		Quantity:    qty,
		Operator:    operator,
		Reason:      reason,
		Notes:       notes,
	}
	if err := a.db.CreateTransaction(tx); err != nil {
		setFlash(w, "danger", "入库失败: "+err.Error())
		http.Redirect(w, r, "/transactions/inbound", http.StatusFound)
		return
	}
	if err := a.db.UpdateComponentQuantity(componentID, qty); err != nil {
		setFlash(w, "danger", "更新库存失败: "+err.Error())
		http.Redirect(w, r, "/transactions/inbound", http.StatusFound)
		return
	}

	setFlash(w, "success", "入库成功！")
	http.Redirect(w, r, "/transactions/inbound", http.StatusFound)
}

func (a *App) handleOutbound(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		componentID := r.URL.Query().Get("component_id")
		components, _ := a.db.ListComponents("")
		if components == nil {
			components = []Component{}
		}

		data := &TemplateData{
			Title:      "出库",
			ActivePage: "outbound",
			Data: map[string]interface{}{
				"Components":     components,
				"SelectedCompID": componentID,
			},
		}
		a.render(w, r, "outbound", data)
		return
	}

	// POST
	if !a.verifyCSRF(r) {
		http.Error(w, "CSRF验证失败", http.StatusForbidden)
		return
	}

	componentID, _ := strconv.Atoi(r.FormValue("component_id"))
	qty, _ := strconv.Atoi(r.FormValue("quantity"))
	operator := r.FormValue("operator")
	reason := r.FormValue("reason")
	notes := r.FormValue("notes")

	if componentID == 0 || qty <= 0 {
		setFlash(w, "danger", "请选择元器件并输入有效数量")
		http.Redirect(w, r, "/transactions/outbound", http.StatusFound)
		return
	}

	// Check stock
	currentQty, err := a.db.GetComponentQuantity(componentID)
	if err != nil {
		setFlash(w, "danger", "查询库存失败")
		http.Redirect(w, r, "/transactions/outbound", http.StatusFound)
		return
	}
	if qty > currentQty {
		setFlash(w, "danger", "出库数量超过当前库存("+strconv.Itoa(currentQty)+")")
		http.Redirect(w, r, "/transactions/outbound", http.StatusFound)
		return
	}

	tx := &Transaction{
		ComponentID: componentID,
		Type:        "outbound",
		Quantity:    qty,
		Operator:    operator,
		Reason:      reason,
		Notes:       notes,
	}
	if err := a.db.CreateTransaction(tx); err != nil {
		setFlash(w, "danger", "出库失败: "+err.Error())
		http.Redirect(w, r, "/transactions/outbound", http.StatusFound)
		return
	}
	if err := a.db.UpdateComponentQuantity(componentID, -qty); err != nil {
		setFlash(w, "danger", "更新库存失败: "+err.Error())
		http.Redirect(w, r, "/transactions/outbound", http.StatusFound)
		return
	}

	setFlash(w, "success", "出库成功！")
	http.Redirect(w, r, "/transactions/outbound", http.StatusFound)
}

func (a *App) handleHistory(w http.ResponseWriter, r *http.Request) {
	filter := TransactionFilter{
		Type:        r.URL.Query().Get("type"),
		ComponentID: 0,
		DateFrom:    r.URL.Query().Get("date_from"),
		DateTo:      r.URL.Query().Get("date_to"),
	}
	if cid := r.URL.Query().Get("component_id"); cid != "" {
		filter.ComponentID, _ = strconv.Atoi(cid)
	}

	transactions, _ := a.db.ListTransactions(filter)
	if transactions == nil {
		transactions = []TransactionWithComponent{}
	}

	components, _ := a.db.ListComponents("")
	if components == nil {
		components = []Component{}
	}

	data := &TemplateData{
		Title:      "出入库记录",
		ActivePage: "history",
		Data: map[string]interface{}{
			"Transactions": transactions,
			"Components":   components,
			"Filter":       filter,
		},
	}
	a.render(w, r, "history", data)
}
