package intouchpay

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents an IntouchPay client configured with authentication details
type Client struct {
	Username        string
	AccountNo       string
	PartnerPassword string
	HTTPClient      *http.Client
}

// NewClient creates a new IntouchPay client
func NewClient(username, accountNumber, partnerPassword string) *Client {
	return &Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// generatePassword calculates the SHA256 hash of the password string
func (c *Client) generatePassword(timestamp string) string {
	data := c.Username + c.AccountNo + c.PartnerPassword + timestamp
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// doRequest sends an HTTP request and handles the response
func (c *Client) doRequest(method, endpoint string, data interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(data)
	bodyReader := bytes.NewReader(jsonBody)
	body := io.NopCloser(bodyReader)
	requestUrl, _ := url.Parse(BaseUrl + endpoint)
	req := http.Request{
		Method: method,
		URL:    requestUrl,
		Body:   body,
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(&req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("IntouchPay API error: %d %s", resp.StatusCode, resp.Status))
	}
	return resp, nil
}

// RequestPayment initiates a payment request
func (c *Client) RequestPayment(params *RequestPaymentParams) (*RequestPaymentResponse, error) {
	var cResp *RequestPaymentResponse
	resp, err := c.doRequest(http.MethodPost, RequestPaymentEndpoint, params)
	if err != nil {
		return cResp, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return cResp, err
	}
	return cResp, nil
}

// RequestDeposit initiates a deposit request
func (c *Client) RequestDeposit(params *RequestDepositParams) (*RequestDepositResponse, error) {
	var cResp *RequestDepositResponse
	resp, err := c.doRequest(http.MethodPost, RequestDepositEndpoint, params)
	if err != nil {
		return cResp, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return cResp, err
	}
	return cResp, nil
}

// GetBalance queries account balance
func (c *Client) GetBalance(params *RequestBalanceParams) (*BalanceResponse, error) {
	var cResp *BalanceResponse
	resp, err := c.doRequest(http.MethodGet, GetBalanceEndpoint, params)
	if err != nil {
		return cResp, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return cResp, err
	}
	return cResp, nil
}
