package Intouchpay_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	Intouchpay "github.com/samueltuyizere/go-intouchpay"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClient implements APIRequester for testing
type MockHTTPClient struct {
	Response *map[string]interface{}
	Error    error
	Called   bool
}

func (m *MockHTTPClient) Do(_ string, _ interface{}) (*map[string]interface{}, error) {
	m.Called = true
	return m.Response, m.Error
}

// TestNewHTTPClient creates a new HTTP client
func TestNewHTTPClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	client := Intouchpay.NewHTTPClient(httpClient, "https://example.com/api")

	assert.NotNil(t, client)
}

// TestHTTPClientDoSuccess tests successful HTTP request
func TestHTTPClientDoSuccess(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		response := map[string]interface{}{
			"success": true,
			"message": "OK",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client := Intouchpay.NewHTTPClient(httpClient, server.URL)

	body := map[string]string{"test": "value"}
	resp, err := client.Do("/test", body)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, true, (*resp)["success"])
}

// TestHTTPClientDoError tests HTTP request with error response
func TestHTTPClientDoError(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"message": "Bad request",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client := Intouchpay.NewHTTPClient(httpClient, server.URL)

	body := map[string]string{"test": "value"}
	resp, err := client.Do("/test", body)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.True(t, Intouchpay.IsAPIError(err))
}

// TestMockHTTPClient tests using mock HTTP client for testing
func TestMockHTTPClient(t *testing.T) {
	mockResp := &map[string]interface{}{
		"status":               "success",
		"requesttransactionid": "TX123",
		"success":              true,
	}
	mockClient := &MockHTTPClient{Response: mockResp}

	auth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	client := Intouchpay.NewClientWithHTTPClient(auth, mockClient)

	params := &Intouchpay.RequestPaymentParams{
		Amount:               1000,
		MobilePhone:          "0781234567",
		RequestTransactionID: "TX123",
	}

	resp, err := client.RequestPayment(params)

	assert.NoError(t, err)
	assert.True(t, mockClient.Called)
	assert.NotNil(t, resp)
}

// TestRequestPaymentWithMockHTTP tests request construction with mock
func TestRequestPaymentWithMockHTTP(t *testing.T) {
	mockResp := &map[string]interface{}{
		"status":               "pending",
		"requesttransactionid": "TX456",
		"success":              true,
		"transactionid":        "TR789",
	}
	mockClient := &MockHTTPClient{Response: mockResp}

	auth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	client := Intouchpay.NewClientWithHTTPClient(auth, mockClient, Intouchpay.WithCallbackURL("https://example.com/callback"))
	client.AccountNo = "ACC123"

	params := &Intouchpay.RequestPaymentParams{
		Amount:               5000,
		MobilePhone:          "0781234567",
		RequestTransactionID: "TX456",
	}

	resp, err := client.RequestPayment(params)

	assert.NoError(t, err)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, "TX456", resp.RequestTransactionID)
	assert.True(t, resp.Success)
}

// TestGetBalanceWithMockHTTP tests GetBalance with mock
func TestGetBalanceWithMockHTTP(t *testing.T) {
	mockResp := &map[string]interface{}{
		"balance": 10000.50,
		"success": true,
	}
	mockClient := &MockHTTPClient{Response: mockResp}

	auth := &MockAuthenticator{
		Creds: Intouchpay.Credentials{
			Username:  "test_user",
			Timestamp: "20260320120000",
			Password:  "test_hash",
		},
	}

	client := Intouchpay.NewClientWithHTTPClient(auth, mockClient)
	client.AccountNo = "ACC123"

	resp, err := client.GetBalance()

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 10000.50, resp.Balance)
}

// TestWithHTTPClientInterfaceOption tests the WithHTTPClientInterface option
func TestWithHTTPClientInterfaceOption(t *testing.T) {
	mockResp := &map[string]interface{}{
		"balance": 1234.56,
		"success": true,
	}
	mockClient := &MockHTTPClient{Response: mockResp}
	auth := Intouchpay.NewAuthenticator("user", "acc", "pass")

	client := Intouchpay.NewClientWithAuth(auth, Intouchpay.WithHTTPClientInterface(mockClient))

	assert.NotNil(t, client)
	// The mock client should be set (verified by using the client)
	resp, err := client.GetBalance()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, mockClient.Called)
}
