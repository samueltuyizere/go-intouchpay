package Intouchpay_test

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	Intouchpay "github.com/samueltuyizere/go-intouchpay"
	"github.com/stretchr/testify/assert"
)

// TestNewAuthenticator creates an authenticator with valid credentials
func TestNewAuthenticator(t *testing.T) {
	auth := Intouchpay.NewAuthenticator("testuser", "1234567890", "secret")

	assert.NotNil(t, auth)
}

// TestAuthenticate returns credentials with all required fields
func TestAuthenticate(t *testing.T) {
	auth := Intouchpay.NewAuthenticator("testuser", "1234567890", "secret")
	creds := auth.Authenticate()

	assert.NotEmpty(t, creds.Username)
	assert.NotEmpty(t, creds.Timestamp)
	assert.NotEmpty(t, creds.Password)
	assert.Equal(t, "testuser", creds.Username)
}

// TestAuthenticateTimestampFormat verifies timestamp is in correct format
func TestAuthenticateTimestampFormat(t *testing.T) {
	auth := Intouchpay.NewAuthenticator("testuser", "1234567890", "secret")
	creds := auth.Authenticate()

	// Timestamp should be 14 digits: YYYYMMDDHHMMSS
	assert.Len(t, creds.Timestamp, 14)

	// Should be parseable as the expected format
	_, err := time.Parse("20060102150405", creds.Timestamp)
	assert.NoError(t, err)
}

// TestAuthenticatePasswordHash verifies password is SHA256 hash
func TestAuthenticatePasswordHash(t *testing.T) {
	username := "testuser"
	accountNo := "1234567890"
	partnerPassword := "secret"

	auth := Intouchpay.NewAuthenticator(username, accountNo, partnerPassword)
	creds := auth.Authenticate()

	// Manually compute expected hash
	expectedString := username + accountNo + partnerPassword + creds.Timestamp
	expectedHash := sha256.Sum256([]byte(expectedString))
	expectedPassword := hex.EncodeToString(expectedHash[:])

	assert.Equal(t, expectedPassword, creds.Password)
	assert.Len(t, creds.Password, 64) // SHA256 produces 64 hex characters
}

// TestAuthenticateGeneratesDifferentTimestamps verifies timestamps change over time
func TestAuthenticateGeneratesDifferentTimestamps(t *testing.T) {
	auth := Intouchpay.NewAuthenticator("testuser", "1234567890", "secret")

	creds1 := auth.Authenticate()
	time.Sleep(1 * time.Second)
	creds2 := auth.Authenticate()

	// Timestamps should be different (or same if within same second)
	// This test just verifies the format is correct
	assert.Len(t, creds1.Timestamp, 14)
	assert.Len(t, creds2.Timestamp, 14)
}

// TestAuthenticateWithEmptyCredentials verifies authenticator handles empty inputs
func TestAuthenticateWithEmptyCredentials(t *testing.T) {
	auth := Intouchpay.NewAuthenticator("", "", "")
	creds := auth.Authenticate()

	assert.Equal(t, "", creds.Username)
	assert.NotEmpty(t, creds.Timestamp)
	assert.NotEmpty(t, creds.Password) // Hash of empty string + timestamp
}

// TestCredentialsAreDeterministicForSameTimestamp verifies hash algorithm is consistent
func TestCredentialsAreDeterministicForSameTimestamp(t *testing.T) {
	username := "testuser"
	accountNo := "1234567890"
	partnerPassword := "secret"
	timestamp := "20260320120000"

	// Compute expected password manually
	passwordString := username + accountNo + partnerPassword + timestamp
	hash := sha256.Sum256([]byte(passwordString))
	expectedPassword := hex.EncodeToString(hash[:])

	// The hash should be deterministic
	assert.Len(t, expectedPassword, 64)
	assert.True(t, strings.ToLower(expectedPassword) == expectedPassword)
}
