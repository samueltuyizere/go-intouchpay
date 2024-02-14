package intouchpay_test

import (
	"net/http"
	"testing"
	"time"

	intouchpay "github.com/samueltuyizere/go-intouchpay"
)

// Should create a new instance of Client with the provided username, account number, and partner password
func TestNewClientWithValidInputs(t *testing.T) {
	username := "testuser"
	accountNumber := "123456789"
	partnerPassword := "password"

	client := intouchpay.NewClient(username, accountNumber, partnerPassword)

	if client.Username != username {
		t.Errorf("Expected username to be %s, but got %s", username, client.Username)
	}
	if client.AccountNo != accountNumber {
		t.Errorf("Expected account number to be %s, but got %s", accountNumber, client.AccountNo)
	}
	if client.PartnerPassword != partnerPassword {
		t.Errorf("Expected partner password to be %s, but got %s", partnerPassword, client.PartnerPassword)
	}
	if client.HTTPClient.Timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5 seconds, but got %s", client.HTTPClient.Timeout)
	}
}

// Should set the HTTPClient timeout to 5 seconds
func TestNewClientWithDefaultTimeout(t *testing.T) {
	username := "testuser"
	accountNumber := "123456789"
	partnerPassword := "password"

	client := intouchpay.NewClient(username, accountNumber, partnerPassword)

	if client.HTTPClient.Timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5 seconds, but got %s", client.HTTPClient.Timeout)
	}
}

// Should handle empty username, account number, and partner password strings
func TestNewClientWithEmptyInputs(t *testing.T) {
	username := ""
	accountNumber := ""
	partnerPassword := ""

	client := intouchpay.NewClient(username, accountNumber, partnerPassword)

	if client.Username != username {
		t.Errorf("Expected username to be %s, but got %s", username, client.Username)
	}
	if client.AccountNo != accountNumber {
		t.Errorf("Expected account number to be %s, but got %s", accountNumber, client.AccountNo)
	}
	if client.PartnerPassword != partnerPassword {
		t.Errorf("Expected partner password to be %s, but got %s", partnerPassword, client.PartnerPassword)
	}
}

// Should handle invalid or empty HTTPClient
func TestNewClientWithInvalidHTTPClient(t *testing.T) {
	username := "testuser"
	accountNumber := "123456789"
	partnerPassword := "password"

	client := &intouchpay.Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		HTTPClient:      nil,
	}

	if client.HTTPClient != nil {
		t.Errorf("Expected HTTPClient to be nil, but got %v", client.HTTPClient)
	}
}

// Should handle invalid or empty timeout value
func TestNewClientWithInvalidTimeout(t *testing.T) {
	username := "testuser"
	accountNumber := "123456789"
	partnerPassword := "password"

	client := &intouchpay.Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		HTTPClient: &http.Client{
			Timeout: 0,
		},
	}

	if client.HTTPClient.Timeout != 0 {
		t.Errorf("Expected timeout to be 0, but got %s", client.HTTPClient.Timeout)
	}
}

// Should handle special characters and whitespace in username, account number, and partner password strings
func TestNewClientWithSpecialCharacters(t *testing.T) {
	username := "test@user"
	accountNumber := "123 456 789"
	partnerPassword := "pass word"

	client := intouchpay.NewClient(username, accountNumber, partnerPassword)

	if client.Username != username {
		t.Errorf("Expected username to be %s, but got %s", username, client.Username)
	}
	if client.AccountNo != accountNumber {
		t.Errorf("Expected account number to be %s, but got %s", accountNumber, client.AccountNo)
	}
	if client.PartnerPassword != partnerPassword {
		t.Errorf("Expected partner password to be %s, but got %s", partnerPassword, client.PartnerPassword)
	}
}
