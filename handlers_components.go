package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (a *App) handleComponentsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("q")
	components, err := a.db.ListComponents(search)
	if err != nil {
		components = []Component{}
	}
	if components == nil {
		components = []Component{}
	}

	data := &TemplateData{
		Title:      "元器件列表",
		ActivePage: "components",
		Data: map[string]interface{}{
			"Components": components,
			"Search":     search,
		},
	}
	a.render(w, r, "component_list", data)
}

func (a *App) handleComponentAddForm(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Title:      "新增元器件",
		ActivePage: "components",
		Data: map[string]interface{}{
			"IsEdit":    false,
			"Component": &Component{},
		},
	}
	a.render(w, r, "component_form", data)
}

func (a *App) handleComponentAdd(w http.ResponseWriter, r *http.Request) {
	if !a.verifyCSRF(r) {
		http.Error(w, "CSRF验证失败", http.StatusForbidden)
		return
	}

	name := r.FormValue("name")
	model := r.FormValue("model")
	pkg := r.FormValue("package")
	location := r.FormValue("location")
	notes := r.FormValue("notes")
	qty, _ := strconv.Atoi(r.FormValue("quantity"))

	if name == "" || model == "" {
		setFlash(w, "danger", "名称和型号为必填项")
		http.Redirect(w, r, "/components/add", http.StatusFound)
		return
	}

	// Check duplicate
	dup, _ := a.db.IsDuplicateComponent(name, model, pkg, 0)
	if dup {
		setFlash(w, "danger", "已存在相同名称、型号、封装的元器件")
		http.Redirect(w, r, "/components/add", http.StatusFound)
		return
	}

	c := &Component{
		Name:     name,
		Model:    model,
		Package:  pkg,
		Quantity: qty,
		Location: location,
		Notes:    notes,
	}
	id, err := a.db.CreateComponent(c)
	if err != nil {
		setFlash(w, "danger", "创建失败: "+err.Error())
		http.Redirect(w, r, "/components/add", http.StatusFound)
		return
	}

	setFlash(w, "success", "元器件创建成功")
	http.Redirect(w, r, "/components/"+strconv.Itoa(id), http.StatusFound)
}

func (a *App) handleComponentDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	comp, err := a.db.GetComponent(id)
	if err != nil {
		http.Error(w, "元器件不存在", http.StatusNotFound)
		return
	}

	transactions, _ := a.db.GetComponentTransactions(id)
	if transactions == nil {
		transactions = []Transaction{}
	}

	data := &TemplateData{
		Title:      comp.Name,
		ActivePage: "components",
		Data: map[string]interface{}{
			"Component":    comp,
			"Transactions": transactions,
		},
	}
	a.render(w, r, "component_detail", data)
}

func (a *App) handleComponentEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	comp, err := a.db.GetComponent(id)
	if err != nil {
		http.Error(w, "元器件不存在", http.StatusNotFound)
		return
	}

	data := &TemplateData{
		Title:      "编辑元器件",
		ActivePage: "components",
		Data: map[string]interface{}{
			"IsEdit":    true,
			"Component": comp,
		},
	}
	a.render(w, r, "component_form", data)
}

func (a *App) handleComponentEdit(w http.ResponseWriter, r *http.Request) {
	if !a.verifyCSRF(r) {
		http.Error(w, "CSRF验证失败", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	model := r.FormValue("model")
	pkg := r.FormValue("package")
	location := r.FormValue("location")
	notes := r.FormValue("notes")

	if name == "" || model == "" {
		setFlash(w, "danger", "名称和型号为必填项")
		http.Redirect(w, r, "/components/"+strconv.Itoa(id)+"/edit", http.StatusFound)
		return
	}

	// Check duplicate (excluding self)
	dup, _ := a.db.IsDuplicateComponent(name, model, pkg, id)
	if dup {
		setFlash(w, "danger", "已存在相同名称、型号、封装的元器件")
		http.Redirect(w, r, "/components/"+strconv.Itoa(id)+"/edit", http.StatusFound)
		return
	}

	c := &Component{
		ID:       id,
		Name:     name,
		Model:    model,
		Package:  pkg,
		Location: location,
		Notes:    notes,
	}
	if err := a.db.UpdateComponent(c); err != nil {
		setFlash(w, "danger", "更新失败: "+err.Error())
		http.Redirect(w, r, "/components/"+strconv.Itoa(id)+"/edit", http.StatusFound)
		return
	}

	setFlash(w, "success", "元器件更新成功")
	http.Redirect(w, r, "/components/"+strconv.Itoa(id), http.StatusFound)
}

func (a *App) handleComponentDelete(w http.ResponseWriter, r *http.Request) {
	if !a.verifyCSRF(r) {
		http.Error(w, "CSRF验证失败", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	hasTx, _ := a.db.HasTransactions(id)
	if hasTx {
		setFlash(w, "danger", "该元器件有出入库记录，无法删除")
		http.Redirect(w, r, "/components/"+strconv.Itoa(id), http.StatusFound)
		return
	}

	if err := a.db.DeleteComponent(id); err != nil {
		setFlash(w, "danger", "删除失败: "+err.Error())
		http.Redirect(w, r, "/components/"+strconv.Itoa(id), http.StatusFound)
		return
	}

	setFlash(w, "success", "元器件已删除")
	http.Redirect(w, r, "/components", http.StatusFound)
}

func (a *App) handleComponentsSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	components, err := a.db.SearchComponents(q)
	if err != nil {
		components = []Component{}
	}
	if components == nil {
		components = []Component{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(components)
}
