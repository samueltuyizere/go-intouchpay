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
	// Remove any existing country code prefix
	cleanedNumber := phoneNumber
	if len(phoneNumber) > 3 && phoneNumber[:3] == "250" {
		return cleanedNumber, nil
	}

	// Validate the number without country code
	isValidMtn := validate_rw_phone_numbers.ValidateMtn(cleanedNumber)
	isValidAirtelTigo := validate_rw_phone_numbers.ValidateAirtelTigo(cleanedNumber)
	if !isValidMtn && !isValidAirtelTigo {
		return "", newValidationError("mobilePhone", "invalid phone number format")
	}

	// Return with country code prefix
	newNumber := fmt.Sprintf("25%s", cleanedNumber)
	return newNumber, nil
}

// SanitizePhoneNumber is a package-level function for backward compatibility
// It validates and formats a Rwandan phone number with the "250" country code prefix
func SanitizePhoneNumber(phoneNumber string) (string, error) {
	// Remove any existing country code prefix
	cleanedNumber := phoneNumber
	if len(phoneNumber) > 3 && phoneNumber[:3] == "250" {
		return cleanedNumber, nil
	}

	// Validate the number without country code
	isValidMtn := validate_rw_phone_numbers.ValidateMtn(cleanedNumber)
	isValidAirtelTigo := validate_rw_phone_numbers.ValidateAirtelTigo(cleanedNumber)
	if !isValidMtn && !isValidAirtelTigo {
		return "", errors.New("invalid phone number")
	}

	// Return with country code prefix
	newNumber := fmt.Sprintf("25%s", cleanedNumber)
	return newNumber, nil
}
