package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"
)

// --- CSRF Protection (Double Submit Cookie) ---

func (a *App) csrfToken(r *http.Request) string {
	if c, err := r.Cookie("csrf_token"); err == nil && c.Value != "" {
		return c.Value
	}
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func setCSRFCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400 * 30,
	})
}

func (a *App) verifyCSRF(r *http.Request) bool {
	cookieToken := ""
	if c, err := r.Cookie("csrf_token"); err == nil {
		cookieToken = c.Value
	}
	formToken := r.FormValue("csrf_token")
	return cookieToken != "" && formToken != "" && cookieToken == formToken
}

// --- Flash Messages ---

func setFlash(w http.ResponseWriter, flashType, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    flashType + ":" + msg,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60,
	})
}

func getFlash(w http.ResponseWriter, r *http.Request) (string, string) {
	c, err := r.Cookie("flash")
	if err != nil {
		return "", ""
	}
	parts := strings.SplitN(c.Value, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	// Clear the flash cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return parts[0], parts[1]
}

// --- Template Rendering ---

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"csrfField": func(token string) template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input type="hidden" name="csrf_token" value="%s">`,
				template.HTMLEscapeString(token),
			))
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"selected": func(val, current interface{}) template.HTMLAttr {
			if val == current {
				return template.HTMLAttr("selected")
			}
			return ""
		},
		"qtyClass": func(qty int) string {
			if qty == 0 {
				return "bg-danger"
			}
			if qty <= 5 {
				return "bg-warning text-dark"
			}
			return "bg-success"
		},
		"qtyLabel": func(qty int) string {
			if qty == 0 {
				return "缺货"
			}
			if qty <= 5 {
				return "低库存"
			}
			return ""
		},
		"typeLabel": TypeLabel,
		"typeBadge": func(t string) string {
			if t == "inbound" {
				return "bg-success"
			}
			return "bg-danger"
		},
		"sub":  func(a, b int) int { return a - b },
		"add":  func(a, b int) int { return a + b },
		"inc":  func(i int) int { return i + 1 },
	}
}

// render executes a pre-built page template (layout + content) with common data.
func (a *App) render(w http.ResponseWriter, r *http.Request, name string, data *TemplateData) {
	// CSRF token
	token := a.csrfToken(r)
	setCSRFCookie(w, token)
	data.CSRFToken = token
	data.Configured = a.configured

	// Flash
	if ft, fm := getFlash(w, r); ft != "" {
		data.FlashType = ft
		data.Flash = fm
	}

	// Look up pre-built template (layout + page content)
	t := a.pages[name]
	if t == nil {
		log.Printf("Template not found: %s", name)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Template error (%s): %v", name, err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
	}
}

// --- Middleware ---

// requireConfig redirects to /settings if database is not configured.
func (a *App) requireConfig(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.configured || a.db == nil {
			http.Redirect(w, r, "/settings", http.StatusFound)
			return
		}
		handler(w, r)
	}
}

// staticFileServer returns an http.Handler for embedded static files.
func staticFileServer(staticFS fs.FS) http.Handler {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	return http.FileServer(http.FS(sub))
}
