package repositories

import "fmt"

// ErrorCode represents different types of repository errors
type ErrorCode int

const (
	ErrCodeUserAlreadyExists ErrorCode = iota
	ErrCodeUserNotFound
	ErrCodeInvalidCredentials
	ErrCodeDatabaseError
	ErrCodeHashingError
	ErrCodeTokenGenerationError
	ErrCodeInvalidInput
	ErrCodeTOTPGenerationError
	ErrCodeInvalidToken
	ErrCodeTokenExpired
)

// RepositoryError represents a custom error type for repository operations
type RepositoryError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// Is allows error comparison
func (e *RepositoryError) Is(target error) bool {
	t, ok := target.(*RepositoryError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(code ErrorCode, message string, err error) *RepositoryError {
	return &RepositoryError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrUserAlreadyExists = &RepositoryError{
		Code:    ErrCodeUserAlreadyExists,
		Message: "user already exists",
	}

	ErrUserNotFound = &RepositoryError{
		Code:    ErrCodeUserNotFound,
		Message: "user not found",
	}

	ErrInvalidCredentials = &RepositoryError{
		Code:    ErrCodeInvalidCredentials,
		Message: "invalid credentials",
	}
)
