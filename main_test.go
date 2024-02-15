package Intouchpay_test

import (
	"testing"
	"time"

	Intouchpay "github.com/samueltuyizere/go-intouchpay"
	"github.com/stretchr/testify/assert"
)

// Returns a new instance of Client with the provided parameters

// Creates a new client with valid input parameters
func TestNewClientWithValidInputParameters(t *testing.T) {
	const (
		username        = "testuser"
		accountNumber   = "1234567890"
		partnerPassword = "password"
		callbackUrl     = "https://example.com/callback"
		sid             = 12345
	)

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	if client.Username != username {
		t.Errorf("Expected username %s, but got %s", username, client.Username)
	}

	if client.AccountNo != accountNumber {
		t.Errorf("Expected account number %s, but got %s", accountNumber, client.AccountNo)
	}

	if client.PartnerPassword != partnerPassword {
		t.Errorf("Expected password %s, but got %s", partnerPassword, client.PartnerPassword)
	}

	if client.CallbackURL != callbackUrl {
		t.Errorf("Expected callback URL %s, but got %s", callbackUrl, client.CallbackURL)
	}

	if client.Sid != sid {
		t.Errorf("Expected SID %d, but got %d", sid, client.Sid)
	}

	if client.HTTPClient.Timeout != 5*time.Second {
		t.Errorf("Expected timeout of %s, but got %s", 5*time.Second, client.HTTPClient.Timeout)
	}
}

// Sets the HTTP client timeout to 5 seconds
func TestNewClientSetsHTTPClientTimeout(t *testing.T) {
	username := "testuser"
	accountNumber := "1234567890"
	partnerPassword := "password"
	callbackUrl := "https://example.com/callback"
	sid := 12345

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	assert.Equal(t, 5*time.Second, client.HTTPClient.Timeout)
}

// Returns a pointer to a new client instance
func TestNewClientReturnsPointerToNewInstance(t *testing.T) {
	username := "testuser"
	accountNumber := "1234567890"
	partnerPassword := "password"
	callbackUrl := "https://example.com/callback"
	sid := 12345

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	assert.NotNil(t, client)
}

// Returns a pointer to a new client instance with empty username
func TestNewClientReturnsPointerToNewInstanceWithEmptyUsername(t *testing.T) {
	username := ""
	accountNumber := "1234567890"
	partnerPassword := "password"
	callbackUrl := "https://example.com/callback"
	sid := 12345

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	assert.NotNil(t, client)
	assert.Equal(t, username, client.Username)
}

// Returns a pointer to a new client instance with empty account number
func TestNewClientReturnsPointerToNewInstanceWithEmptyAccountNumber(t *testing.T) {
	username := "testuser"
	accountNumber := ""
	partnerPassword := "password"
	callbackUrl := "https://example.com/callback"
	sid := 12345

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	assert.NotNil(t, client)
	assert.Equal(t, accountNumber, client.AccountNo)
}

// Returns a pointer to a new client instance with empty partner password
func TestNewClientReturnsPointerToNewInstanceWithEmptyPartnerPassword(t *testing.T) {
	username := "testuser"
	accountNumber := "1234567890"
	partnerPassword := ""
	callbackUrl := "https://example.com/callback"
	sid := 12345

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	assert.NotNil(t, client)
	assert.Equal(t, partnerPassword, client.PartnerPassword)
}
