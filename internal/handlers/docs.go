package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"secretary/alpha/pkg/utils"

	"github.com/gorilla/mux"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) RegisterRoutes(r *mux.Router) {
	// Serve Swagger UI
	r.HandleFunc("/docs", h.SwaggerUI).Methods("GET")
	r.HandleFunc("/docs/", h.SwaggerUI).Methods("GET")

	// Serve the swagger.yaml file
	r.HandleFunc("/docs/swagger.yaml", h.SwaggerYAML).Methods("GET")

	// Serve static assets for Swagger UI (if needed)
	r.PathPrefix("/docs/static/").Handler(http.StripPrefix("/docs/static/", http.FileServer(http.Dir("docs/static/"))))
}

func (h *DocsHandler) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Read the swagger-ui.html file
	htmlPath := filepath.Join("docs", "swagger-ui.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		utils.InternalError(w, "Failed to load API documentation", "Could not read documentation file")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(htmlContent)
}

func (h *DocsHandler) SwaggerYAML(w http.ResponseWriter, r *http.Request) {
	// Read the swagger.yaml file
	yamlPath := filepath.Join("docs", "swagger.yaml")
	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		utils.InternalError(w, "Failed to load API specification", "Could not read swagger.yaml file")
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(yamlContent)
}
