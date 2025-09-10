package application

type ErrorCode string

const (
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	ErrInternal     ErrorCode = "INTERNAL"
)

// AppError is a struct for application-level errors
type AppError struct {
	code    ErrorCode
	message string
}

// Error returns ErrorCode of AppError
func (e *AppError) Code() ErrorCode {
	return e.code
}

// Code returns error message of AppError
func (e *AppError) Error() string {
	return e.message
}

// NewAppError creates a new instance of AppError
func NewAppError(code ErrorCode, message string) *AppError {
	return &AppError{code: code, message: message}
}
