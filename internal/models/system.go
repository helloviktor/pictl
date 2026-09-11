package models

// UpdateResponse contains the response for update operations
type UpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
