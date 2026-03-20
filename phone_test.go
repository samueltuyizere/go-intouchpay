package Intouchpay_test

import (
	"testing"

	Intouchpay "github.com/samueltuyizere/go-intouchpay"
	"github.com/stretchr/testify/assert"
)

// TestSanitizePhoneNumberValidMTN tests valid MTN number
func TestSanitizePhoneNumberValidMTN(t *testing.T) {
	// MTN prefix: 078, 079
	phone, err := Intouchpay.SanitizePhoneNumber("0781234567")
	assert.NoError(t, err)
	assert.Equal(t, "250781234567", phone)
}

// TestSanitizePhoneNumberValidAirtelTigo tests valid Airtel/Tigo number
func TestSanitizePhoneNumberValidAirtelTigo(t *testing.T) {
	// Airtel/Tigo prefix: 072, 073
	phone, err := Intouchpay.SanitizePhoneNumber("0721234567")
	assert.NoError(t, err)
	assert.Equal(t, "250721234567", phone)
}

// TestSanitizePhoneNumberAlreadyPrefixed tests number already with 250 prefix
func TestSanitizePhoneNumberAlreadyPrefixed(t *testing.T) {
	phone, err := Intouchpay.SanitizePhoneNumber("250781234567")
	assert.NoError(t, err)
	assert.Equal(t, "250781234567", phone)
}

// TestSanitizePhoneNumberInvalid tests invalid phone number
func TestSanitizePhoneNumberInvalid(t *testing.T) {
	_, err := Intouchpay.SanitizePhoneNumber("123456789")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid phone number")
}

// TestSanitizePhoneNumberEmpty tests empty phone number
func TestSanitizePhoneNumberEmpty(t *testing.T) {
	_, err := Intouchpay.SanitizePhoneNumber("")
	assert.Error(t, err)
}

// TestSanitizePhoneNumberTooShort tests phone number that's too short
func TestSanitizePhoneNumberTooShort(t *testing.T) {
	_, err := Intouchpay.SanitizePhoneNumber("078123")
	assert.Error(t, err)
}

// TestPhoneValidatorSanitizePhoneNumber tests the PhoneValidator struct method
func TestPhoneValidatorSanitizePhoneNumber(t *testing.T) {
	validator := Intouchpay.NewPhoneValidator()

	phone, err := validator.SanitizePhoneNumber("0781234567")
	assert.NoError(t, err)
	assert.Equal(t, "250781234567", phone)
}

// TestPhoneValidatorInvalidNumber tests PhoneValidator with invalid number
func TestPhoneValidatorInvalidNumber(t *testing.T) {
	validator := Intouchpay.NewPhoneValidator()

	_, err := validator.SanitizePhoneNumber("invalid")
	assert.Error(t, err)
	assert.True(t, Intouchpay.IsValidationError(err))
}

// TestNewPhoneValidator creates a new PhoneValidator
func TestNewPhoneValidator(t *testing.T) {
	validator := Intouchpay.NewPhoneValidator()
	assert.NotNil(t, validator)
}
