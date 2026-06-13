package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*.html templates/**/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

// App holds global application state.
type App struct {
	config     *Config
	db         *DB
	pages      map[string]*template.Template
	configured bool
	dataDir    string
}

func main() {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/app/data"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Load config
	cfg, err := LoadConfig(dataDir)
	if err != nil {
		log.Printf("未找到配置文件，请先配置数据库连接")
		cfg = &Config{DBPort: 3306}
	}
	if cfg.DBPort == 0 {
		cfg.DBPort = 3306
	}

	// Pre-build page templates: each page = layout.html + specific page file
	// This allows {{define "content"}} in each page to be picked up by
	// {{template "content" .}} in layout.html without naming conflicts.
	pageFiles := map[string]string{
		"dashboard":        "templates/dashboard.html",
		"settings":         "templates/settings.html",
		"component_list":   "templates/components/list.html",
		"component_form":   "templates/components/form.html",
		"component_detail": "templates/components/detail.html",
		"inbound":          "templates/transactions/inbound.html",
		"outbound":         "templates/transactions/outbound.html",
		"history":          "templates/transactions/history.html",
		"lcsc_import":      "templates/transactions/lcsc_import.html",
		"lcsc_preview":     "templates/transactions/lcsc_preview.html",
	}
	pages := make(map[string]*template.Template, len(pageFiles))
	for name, file := range pageFiles {
		t := template.New("").Funcs(templateFuncs())
		template.Must(t.ParseFS(templatesFS, "templates/layout.html", file))
		pages[name] = t
	}

	app := &App{
		config:     cfg,
		pages:      pages,
		configured: cfg.IsConfigured(),
		dataDir:    dataDir,
	}

	// Connect to MySQL if configured
	if app.configured {
		db, err := NewDB(cfg)
		if err != nil {
			log.Printf("连接 MySQL 失败: %v，请检查配置", err)
			app.configured = false
		} else {
			app.db = db
			if err := db.CreateTables(); err != nil {
				log.Printf("创建数据表失败: %v", err)
			} else {
				log.Println("已连接 MySQL，数据表就绪")
			}
		}
	}

	// Routes
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer(staticFS)))

	// Dashboard
	mux.HandleFunc("/", app.handleDashboard)

	// Settings
	mux.HandleFunc("/settings", app.handleSettings)
	mux.HandleFunc("/settings/test", app.handleSettingsTest)

	// Components
	mux.HandleFunc("/components", app.requireConfig(app.handleComponentsList))
	mux.HandleFunc("GET /components/add", app.requireConfig(app.handleComponentAddForm))
	mux.HandleFunc("POST /components/add", app.requireConfig(app.handleComponentAdd))
	mux.HandleFunc("GET /components/search", app.requireConfig(app.handleComponentsSearch))
	mux.HandleFunc("GET /components/{id}", app.requireConfig(app.handleComponentDetail))
	mux.HandleFunc("GET /components/{id}/edit", app.requireConfig(app.handleComponentEditForm))
	mux.HandleFunc("POST /components/{id}/edit", app.requireConfig(app.handleComponentEdit))
	mux.HandleFunc("POST /components/{id}/delete", app.requireConfig(app.handleComponentDelete))

	// Transactions
	mux.HandleFunc("/transactions/inbound", app.requireConfig(app.handleInbound))
	mux.HandleFunc("/transactions/outbound", app.requireConfig(app.handleOutbound))
	mux.HandleFunc("/transactions/history", app.requireConfig(app.handleHistory))

	// LCSC import
	mux.HandleFunc("GET /transactions/lcsc", app.requireConfig(app.handleLCSCForm))
	mux.HandleFunc("POST /transactions/lcsc", app.requireConfig(app.handleLCSCUpload))
	mux.HandleFunc("POST /transactions/lcsc/confirm", app.requireConfig(app.handleLCSCConfirm))

	log.Printf("Ecomp 启动于 :%s (配置: %v)", port, app.configured)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
