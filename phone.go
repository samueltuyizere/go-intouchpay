package Intouchpay

import (
	"errors"
	"fmt"

	"github.com/samueltuyizere/validate_rw_phone_numbers"
)

// PhoneValidator validates and sanitizes Rwandan phone numbers
type PhoneValidator struct{}

// NewPhoneValidator creates a new PhoneValidator
func NewPhoneValidator() *PhoneValidator {
	return &PhoneValidator{}
}

// SanitizePhoneNumber validates and formats a Rwandan phone number
// It returns the number with the "250" country code prefix
func (p *PhoneValidator) SanitizePhoneNumber(phoneNumber string) (string, error) {
	// Remove any existing country code prefix to validate the base number
	cleanedNumber := phoneNumber
	if len(phoneNumber) > 3 && phoneNumber[:3] == "250" {
		cleanedNumber = phoneNumber[3:] // Strip prefix to validate
	}

	// The validation library expects local format with leading "0"
	// If the number doesn't start with "0", add it for validation
	validationNumber := cleanedNumber
	if len(cleanedNumber) > 0 && cleanedNumber[0] != '0' {
		validationNumber = "0" + cleanedNumber
	}

	// Validate the number without country code
	isValidMtn := validate_rw_phone_numbers.ValidateMtn(validationNumber)
	isValidAirtelTigo := validate_rw_phone_numbers.ValidateAirtelTigo(validationNumber)
	if !isValidMtn && !isValidAirtelTigo {
		return "", newValidationError("mobilePhone", "invalid phone number format")
	}

	// Return with country code prefix
	newNumber := fmt.Sprintf("25%s", validationNumber)
	return newNumber, nil
}

// SanitizePhoneNumber is a package-level function for backward compatibility
// It validates and formats a Rwandan phone number with the "250" country code prefix
func SanitizePhoneNumber(phoneNumber string) (string, error) {
	// Remove any existing country code prefix to validate the base number
	cleanedNumber := phoneNumber
	if len(phoneNumber) > 3 && phoneNumber[:3] == "250" {
		cleanedNumber = phoneNumber[3:] // Strip prefix to validate
	}

	// The validation library expects local format with leading "0"
	// If the number doesn't start with "0", add it for validation
	validationNumber := cleanedNumber
	if len(cleanedNumber) > 0 && cleanedNumber[0] != '0' {
		validationNumber = "0" + cleanedNumber
	}

	// Validate the number without country code
	isValidMtn := validate_rw_phone_numbers.ValidateMtn(validationNumber)
	isValidAirtelTigo := validate_rw_phone_numbers.ValidateAirtelTigo(validationNumber)
	if !isValidMtn && !isValidAirtelTigo {
		return "", errors.New("invalid phone number")
	}

	// Return with country code prefix
	newNumber := fmt.Sprintf("25%s", validationNumber)
	return newNumber, nil
}
