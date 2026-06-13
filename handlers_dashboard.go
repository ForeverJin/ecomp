package main

import (
	"net/http"
)

func (a *App) handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Redirect to settings if not configured
	if !a.configured || a.db == nil {
		http.Redirect(w, r, "/settings", http.StatusFound)
		return
	}

	// Only handle exact "/" path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	stats, err := a.db.GetStats()
	if err != nil {
		stats = &Stats{}
	}

	recent, _ := a.db.GetRecentTransactions(10)
	lowStock, _ := a.db.GetLowStockComponents()
	outStock, _ := a.db.GetOutStockComponents()
	components, _ := a.db.SearchComponents("")

	if recent == nil {
		recent = []TransactionWithComponent{}
	}
	if lowStock == nil {
		lowStock = []Component{}
	}
	if outStock == nil {
		outStock = []Component{}
	}
	if components == nil {
		components = []Component{}
	}

	data := &TemplateData{
		Title:      "仪表盘",
		ActivePage: "dashboard",
		Data: map[string]interface{}{
			"Stats":      stats,
			"Recent":     recent,
			"LowStock":   lowStock,
			"OutStock":   outStock,
			"Components": components,
		},
	}
	a.render(w, r, "dashboard", data)
}
