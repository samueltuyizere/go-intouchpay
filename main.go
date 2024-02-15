package Intouchpay

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// NewClient creates a new IntouchPay client
func NewClient(username, accountNumber, partnerPassword, callbackUrl string, sid int) *Client {
	return &Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		CallbackURL:     callbackUrl,
		Sid:             sid,
		HTTPClient:      &http.Client{},
	}
}

// generatePassword calculates the SHA256 hash of the password string
func (c *Client) generatePassword() string {
	timestamp := time.Now().String()
	data := c.Username + c.AccountNo + c.PartnerPassword + timestamp
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// doRequest sends an HTTP request and handles the response
func (c *Client) doRequest(method, endpoint string, data interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(data)
	body := bytes.NewReader(jsonBody)
	requestUrl, _ := url.Parse(BaseUrl + endpoint)
	req, _ := http.NewRequest(method, requestUrl.RequestURI(), body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IntouchPay API error: %d %s", resp.StatusCode, resp.Status)
	}
	return resp, nil
}

// RequestPayment initiates a payment request
func (c *Client) RequestPayment(params *RequestPaymentParams) (*RequestPaymentResponse, error) {
	password := c.generatePassword()
	var cResp *RequestPaymentResponse
	body := map[string]interface{}{
		"username":             c.Username,
		"timestamp":            time.Now(),
		"amount":               params.Amount,
		"password":             password,
		"mobilephone":          params.MobilePhone,
		"requesttransactionid": params.RequestTransactionId,
		"callbackurl":          c.CallbackURL,
	}
	resp, err := c.doRequest(http.MethodPost, RequestPaymentEndpoint, body)
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
	password := c.generatePassword()
	var cResp *RequestDepositResponse
	body := map[string]interface{}{
		"username":             c.Username,
		"timestamp":            time.Now(),
		"amount":               params.Amount,
		"withdrawcharge":       0,
		"reason":               params.Reason,
		"sid":                  c.Sid,
		"password":             password,
		"mobilephone":          params.MobilePhone,
		"requesttransactionid": params.RequestTransactionId,
	}
	resp, err := c.doRequest(http.MethodPost, RequestDepositEndpoint, body)
	if err != nil {
		return cResp, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return cResp, err
	}
	return cResp, nil
}

// GetBalance queries account balance
func (c *Client) GetBalance() (*BalanceResponse, error) {
	timestamp := time.Now()
	var cResp *BalanceResponse
	resp, err := c.doRequest(http.MethodGet, GetBalanceEndpoint, timestamp)
	if err != nil {
		return cResp, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return cResp, err
	}
	return cResp, nil
}
