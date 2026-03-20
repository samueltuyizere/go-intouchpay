package Intouchpay

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the IntouchPay API
type APIError struct {
	StatusCode int
	Status     string
	Response   map[string]interface{}
	Message    string
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("IntouchPay API error: %d %s - %v", e.StatusCode, e.Status, e.Response)
}

// newAPIError creates an APIError from a failed response
func newAPIError(statusCode int, status string, response map[string]interface{}) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Status:     status,
		Response:   response,
		Message:    fmt.Sprintf("API request failed with status %d", statusCode),
	}
}

// NewAPIErrorForTest is exposed for testing purposes only
func NewAPIErrorForTest(statusCode int, status string, response map[string]interface{}) *APIError {
	return newAPIError(statusCode, status, response)
}

// ValidationError represents a client-side validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// newValidationError creates a ValidationError
func newValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ParseAPIError attempts to parse an error response from the API
func ParseAPIError(statusCode int, status string, body json.RawMessage) error {
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return &APIError{
			StatusCode: statusCode,
			Status:     status,
			Message:    "failed to parse error response",
		}
	}
	return newAPIError(statusCode, status, response)
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// MarshalError represents an error during JSON marshaling/unmarshaling
type MarshalError struct {
	Context string
	Err     error
}

// Error implements the error interface
func (e *MarshalError) Error() string {
	return fmt.Sprintf("marshal error: %s: %v", e.Context, e.Err)
}

// Unwrap returns the underlying error
func (e *MarshalError) Unwrap() error {
	return e.Err
}

// NewMarshalError creates a MarshalError
func NewMarshalError(context string, err error) *MarshalError {
	return &MarshalError{
		Context: context,
		Err:     err,
	}
}

// IsMarshalError checks if an error is a MarshalError
func IsMarshalError(err error) bool {
	_, ok := err.(*MarshalError)
	return ok
}

// HTTPStatus returns the HTTP status code if the error is an APIError
func HTTPStatus(err error) int {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode
	}
	return http.StatusInternalServerError
}
