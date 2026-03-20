package Intouchpay_test

import (
	"net/http"
	"testing"

	Intouchpay "github.com/samueltuyizere/go-intouchpay"
	"github.com/stretchr/testify/assert"
)

// TestAPIError tests APIError struct and Error method
func TestAPIError(t *testing.T) {
	err := &Intouchpay.APIError{
		StatusCode: 400,
		Status:     "Bad Request",
		Response:   map[string]interface{}{"message": "Invalid request"},
		Message:    "API request failed",
	}

	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "Bad Request")
	assert.Equal(t, 400, err.StatusCode)
}

// TestValidationError tests ValidationError struct and Error method
func TestValidationError(t *testing.T) {
	err := &Intouchpay.ValidationError{
		Field:   "mobilePhone",
		Message: "invalid format",
	}

	assert.Contains(t, err.Error(), "mobilePhone")
	assert.Contains(t, err.Error(), "invalid format")
}

// TestIsAPIError tests IsAPIError helper function
func TestIsAPIError(t *testing.T) {
	apiErr := &Intouchpay.APIError{StatusCode: 500}
	otherErr := &Intouchpay.ValidationError{Field: "test"}

	assert.True(t, Intouchpay.IsAPIError(apiErr))
	assert.False(t, Intouchpay.IsAPIError(otherErr))
}

// TestIsValidationError tests IsValidationError helper function
func TestIsValidationError(t *testing.T) {
	validationErr := &Intouchpay.ValidationError{Field: "test"}
	apiErr := &Intouchpay.APIError{StatusCode: 500}

	assert.True(t, Intouchpay.IsValidationError(validationErr))
	assert.False(t, Intouchpay.IsValidationError(apiErr))
}

// TestHTTPStatus tests HTTPStatus helper function
func TestHTTPStatus(t *testing.T) {
	apiErr := &Intouchpay.APIError{StatusCode: 400}
	validationErr := &Intouchpay.ValidationError{Field: "test"}

	assert.Equal(t, 400, Intouchpay.HTTPStatus(apiErr))
	assert.Equal(t, http.StatusInternalServerError, Intouchpay.HTTPStatus(validationErr))
}

// TestNewAPIError tests newAPIError constructor
func TestNewAPIError(t *testing.T) {
	response := map[string]interface{}{"error": "test"}
	err := Intouchpay.NewAPIErrorForTest(401, "Unauthorized", response)

	assert.NotNil(t, err)
	assert.Equal(t, 401, err.StatusCode)
	assert.Equal(t, "Unauthorized", err.Status)
}
