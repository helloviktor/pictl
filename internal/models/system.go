package models

// SystemInfo contains system information
type SystemInfo struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	DiskUsage    float64 `json:"disk_usage"`
	CPUTemp      float64 `json:"cpu_temp"`
	LastUpdate   string  `json:"last_update"`
	UpdatesAvail int     `json:"updates_available"`
}

// UpdateResponse contains the response for update operations
type UpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
