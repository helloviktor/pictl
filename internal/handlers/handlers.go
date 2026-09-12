package handlers

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/helloviktor/pictl/internal/models"
)

// SystemService is the subset of system operations Handlers depends on.
type SystemService interface {
	CPUUsage() (float64, error)
	MemoryUsage() (float64, error)
	DiskUsage() (float64, error)
	CPUTemperature() (float64, error)
	AvailableUpdates() (int, error)
	UpdateSystem() error
	RestartSystem() error
	ShutdownSystem() error
}

// Handlers handles HTTP requests
type Handlers struct {
	sysService SystemService
	staticFS   embed.FS
}

// NewHandlers creates a new Handlers instance
func NewHandlers(sysService SystemService, staticFS embed.FS) *Handlers {
	return &Handlers{
		sysService: sysService,
		staticFS:   staticFS,
	}
}

// ServeHTML serves the main HTML file
func (h *Handlers) ServeHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data, err := h.staticFS.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}
}

// ServeStatic serves static files (CSS, JS, etc.)
func (h *Handlers) ServeStatic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a sub-filesystem for static files
		fsys, err := fs.Sub(h.staticFS, "static")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Strip the /static/ prefix and serve files
		handler := http.StripPrefix("/static/", http.FileServer(http.FS(fsys)))
		handler.ServeHTTP(w, r)
	}
}

// GetCPUUsage returns the current CPU usage percentage
func (h *Handlers) GetCPUUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value, err := h.sysService.CPUUsage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// GetMemoryUsage returns the current memory usage percentage
func (h *Handlers) GetMemoryUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value, err := h.sysService.MemoryUsage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// GetDiskUsage returns the current disk usage percentage
func (h *Handlers) GetDiskUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value, err := h.sysService.DiskUsage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// GetCPUTemperature returns the current CPU temperature in Celsius
func (h *Handlers) GetCPUTemperature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value, err := h.sysService.CPUTemperature()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// GetAvailableUpdates returns the number of available package updates
func (h *Handlers) GetAvailableUpdates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value, err := h.sysService.AvailableUpdates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// UpdateSystem handles system update request
func (h *Handlers) UpdateSystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.sysService.UpdateSystem()
	resp := models.UpdateResponse{
		Success: err == nil,
	}
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp.Message = "System update started"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RestartSystem handles system restart request
func (h *Handlers) RestartSystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.sysService.RestartSystem()
	resp := models.UpdateResponse{
		Success: err == nil,
	}
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp.Message = "System restart initiated"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ShutdownSystem handles system shutdown request
func (h *Handlers) ShutdownSystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.sysService.ShutdownSystem()
	resp := models.UpdateResponse{
		Success: err == nil,
	}
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp.Message = "System shutdown initiated"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
