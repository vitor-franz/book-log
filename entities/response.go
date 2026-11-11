package entities

// DefaultResponse acknowledges a completed operation.
type DefaultResponse struct {
	// Success indicates whether the operation completed successfully.
	Success bool `json:"success"`
}

// NewDefaultResponse returns a response with Success set to true.
func NewDefaultResponse() DefaultResponse {
	return DefaultResponse{Success: true}
}
