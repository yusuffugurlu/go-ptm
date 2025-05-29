package errors

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewUnauthorizedError(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}

func NewForbiddenError(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}
