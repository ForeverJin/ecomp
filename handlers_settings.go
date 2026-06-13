package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *App) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := &TemplateData{
			Title:      "数据库设置",
			ActivePage: "settings",
			Data: map[string]interface{}{
				"Config": a.config,
			},
		}
		a.render(w, r, "settings", data)
		return
	}

	// POST: save settings
	cfg := &Config{
		DBHost:     r.FormValue("db_host"),
		DBPort:     3306,
		DBUser:     r.FormValue("db_user"),
		DBPassword: r.FormValue("db_password"),
		DBName:     r.FormValue("db_name"),
	}
	fmt.Sscanf(r.FormValue("db_port"), "%d", &cfg.DBPort)
	if cfg.DBPort == 0 {
		cfg.DBPort = 3306
	}

	// Validate
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		setFlash(w, "danger", "请填写所有必填字段")
		http.Redirect(w, r, "/settings", http.StatusFound)
		return
	}

	// Test connection
	if err := TestConnection(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName); err != nil {
		setFlash(w, "danger", "连接测试失败: "+err.Error())
		http.Redirect(w, r, "/settings", http.StatusFound)
		return
	}

	// Save config
	if err := SaveConfig(a.dataDir, cfg); err != nil {
		setFlash(w, "danger", "保存配置失败: "+err.Error())
		http.Redirect(w, r, "/settings", http.StatusFound)
		return
	}

	// Connect to MySQL
	db, err := NewDB(cfg)
	if err != nil {
		setFlash(w, "danger", "连接数据库失败: "+err.Error())
		http.Redirect(w, r, "/settings", http.StatusFound)
		return
	}

	// Create tables
	if err := db.CreateTables(); err != nil {
		db.Close()
		setFlash(w, "danger", "创建数据表失败: "+err.Error())
		http.Redirect(w, r, "/settings", http.StatusFound)
		return
	}

	// Close old connection if any
	if a.db != nil {
		a.db.Close()
	}

	a.db = db
	a.config = cfg
	a.configured = true

	setFlash(w, "success", "数据库配置已保存并连接成功！")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) handleSettingsTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		DBHost     string `json:"db_host"`
		DBPort     int    `json:"db_port"`
		DBUser     string `json:"db_user"`
		DBPassword string `json:"db_password"`
		DBName     string `json:"db_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, false, "无效的请求数据")
		return
	}

	if req.DBPort == 0 {
		req.DBPort = 3306
	}

	err := TestConnection(req.DBHost, req.DBPort, req.DBUser, req.DBPassword, req.DBName)
	if err != nil {
		jsonResponse(w, http.StatusOK, false, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, true, "连接成功！")
}

func jsonResponse(w http.ResponseWriter, status int, success bool, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": message,
	})
}
