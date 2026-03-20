package Intouchpay_test

import (
	"net/http"
	"testing"
	"time"

	Intouchpay "github.com/samueltuyizere/go-intouchpay"
	"github.com/stretchr/testify/assert"
)

// MockAuthenticator implements Authenticator for testing
type MockAuthenticator struct {
	Creds Intouchpay.Credentials
}

func (m *MockAuthenticator) Authenticate() Intouchpay.Credentials {
	return m.Creds
}

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

	// Fixed: Now expects 60 seconds (DefaultTimeout)
	if client.HTTPClient.Timeout != Intouchpay.DefaultTimeout {
		t.Errorf("Expected timeout of %v, but got %v", Intouchpay.DefaultTimeout, client.HTTPClient.Timeout)
	}
}

// Sets the HTTP client timeout to DefaultTimeout (60 seconds)
func TestNewClientSetsHTTPClientTimeout(t *testing.T) {
	username := "testuser"
	accountNumber := "1234567890"
	partnerPassword := "password"
	callbackUrl := "https://example.com/callback"
	sid := 12345

	client := Intouchpay.NewClient(username, accountNumber, partnerPassword, callbackUrl, sid)

	assert.Equal(t, Intouchpay.DefaultTimeout, client.HTTPClient.Timeout)
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

// Test NewClientWithAuth with mock authenticator
func TestNewClientWithAuth(t *testing.T) {
	mockAuth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	client := Intouchpay.NewClientWithAuth(mockAuth)

	assert.NotNil(t, client)
	assert.Equal(t, Intouchpay.DefaultTimeout, client.HTTPClient.Timeout)
}

// Test NewClientWithAuth with options
func TestNewClientWithAuthWithOptions(t *testing.T) {
	mockAuth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	customTimeout := 30 * time.Second
	client := Intouchpay.NewClientWithAuth(mockAuth, Intouchpay.WithTimeout(customTimeout))

	assert.NotNil(t, client)
	assert.Equal(t, customTimeout, client.HTTPClient.Timeout)
}

// Test NewClientWithOptions
func TestNewClientWithOptions(t *testing.T) {
	client := Intouchpay.NewClientWithOptions(
		"testuser",
		"1234567890",
		"password",
		Intouchpay.WithCallbackURL("https://example.com/callback"),
		Intouchpay.WithSid(12345),
		Intouchpay.WithTimeout(30*time.Second),
	)

	assert.NotNil(t, client)
	assert.Equal(t, "testuser", client.Username)
	assert.Equal(t, "1234567890", client.AccountNo)
	assert.Equal(t, "password", client.PartnerPassword)
	assert.Equal(t, "https://example.com/callback", client.CallbackURL)
	assert.Equal(t, 12345, client.Sid)
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout)
}

// Test WithHTTPClient option
func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 10 * time.Second}
	mockAuth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	client := Intouchpay.NewClientWithAuth(mockAuth, Intouchpay.WithHTTPClient(customClient))

	assert.Equal(t, customClient, client.HTTPClient)
	assert.Equal(t, 10*time.Second, client.HTTPClient.Timeout)
}

// Test WithAuthenticator option
func TestWithAuthenticator(t *testing.T) {
	mockAuth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	client := Intouchpay.NewClientWithOptions("user", "acc", "pass", Intouchpay.WithAuthenticator(mockAuth))

	// Verify the authenticator is used by calling it
	creds := client.GetAuthCredentials()
	assert.Equal(t, "test_user", creds.Username)
	assert.Equal(t, "20260320120000", creds.Timestamp)
	assert.Equal(t, "test_hash", creds.Password)
}
