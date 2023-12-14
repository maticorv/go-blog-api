package errors

type AppError struct {
	Message    []string `json:"message"`
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
}

func NewAppError(message []string, errorMessage string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		Error:      errorMessage,
		StatusCode: statusCode,
	}
}
