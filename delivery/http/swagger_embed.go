package http

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed swagger/openapi.yaml swagger/index.html
var swaggerFS embed.FS

// ServeAPIDocs serves the API docs UI and spec at /api-docs.
// Register with: mux.Handle("/api-docs/", http.StripPrefix("/api-docs", ...))
// and mux.HandleFunc("/api-docs", apiDocsRedirect) so /api-docs redirects to /api-docs/
func ServeAPIDocs(mux *http.ServeMux) {
	sub, _ := fs.Sub(swaggerFS, "swagger")
	docsFS := http.FileServer(http.FS(sub))
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs", docsFS))
	mux.HandleFunc("/api-docs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api-docs" {
			http.NotFound(w, r)
			return
		}
		// Redirect /api-docs to /api-docs/ so relative paths in index.html work
		http.Redirect(w, r, "/api-docs/", http.StatusMovedPermanently)
	})
}
