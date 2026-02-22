package dto

type ErrorResponse struct {
	Code         string        `json:"code"`
	Message      string        `json:"message"`
	ErrorDetails []ErrorDetail `json:"error_details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
