package Intouchpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewClient creates a new IntouchPay client with the provided credentials.
// This is a convenience constructor that creates a default authenticator.
func NewClient(username, accountNumber, partnerPassword, callbackUrl string, sid int) *Client {
	auth := NewAuthenticator(username, accountNumber, partnerPassword)
	c := &Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		CallbackURL:     callbackUrl,
		Sid:             sid,
		auth:            auth,
		HTTPClient:      &http.Client{Timeout: DefaultTimeout},
	}
	return c
}

// NewClientWithAuth creates a new IntouchPay client with a custom authenticator.
// Use this for testing or when you need custom authentication behavior.
func NewClientWithAuth(auth Authenticator, opts ...Option) *Client {
	c := &Client{
		auth:       auth,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewClientWithOptions creates a new IntouchPay client with options.
// Use this for flexible configuration including custom HTTP client, timeout, etc.
func NewClientWithOptions(username, accountNumber, partnerPassword string, opts ...Option) *Client {
	auth := NewAuthenticator(username, accountNumber, partnerPassword)
	c := &Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		auth:            auth,
		HTTPClient:      &http.Client{Timeout: DefaultTimeout},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doRequest sends an HTTP JSON POST request and handles the response
func (c *Client) doRequest(endpoint string, requestBody interface{}) (*map[string]interface{}, error) {
	var response *map[string]interface{}
	requestUrl := BaseUrl + endpoint

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return response, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return response, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log the error if needed, but don't fail the request
			// The response body has already been read at this point
			fmt.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("IntouchPay API error: %d\n %s\n %v", resp.StatusCode, resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return response, newAPIError(resp.StatusCode, resp.Status, *response)
	}
	return response, nil
}

// RequestPayment initiates a payment request
func (c *Client) RequestPayment(params *RequestPaymentParams) (*RequestPaymentResponse, error) {
	phoneNumber, err := SanitizePhoneNumber(params.MobilePhone)
	if err != nil {
		return nil, err
	}

	creds := c.auth.Authenticate()
	requestBody := RequestPaymentBody{
		Username:             creds.Username,
		Timestamp:            creds.Timestamp,
		Amount:               params.Amount,
		Password:             creds.Password,
		MobilePhone:          phoneNumber,
		RequestTransactionId: params.RequestTransactionId,
		AccountNo:            c.AccountNo,
	}
	if c.CallbackURL != "" {
		requestBody.CallbackURL = c.CallbackURL
	}

	var cResp *RequestPaymentResponse
	resp, err := c.doRequest(RequestPaymentEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	err = json.Unmarshal(respBytes, &cResp)
	if err != nil {
		return cResp, err
	}

	return cResp, nil
}

// RequestDeposit initiates a deposit request
func (c *Client) RequestDeposit(params *RequestDepositParams) (*RequestDepositResponse, error) {
	phoneNumber, err := SanitizePhoneNumber(params.MobilePhone)
	if err != nil {
		return nil, err
	}

	creds := c.auth.Authenticate()
	requestBody := RequestDepositBody{
		Username:             creds.Username,
		Timestamp:            creds.Timestamp,
		Amount:               params.Amount,
		WithdrawCharge:       params.WithdrawCharge,
		Reason:               params.Reason,
		Sid:                  c.Sid,
		Password:             creds.Password,
		MobilePhoneNo:        phoneNumber,
		RequestTransactionId: params.RequestTransactionId,
		AccountNo:            c.AccountNo,
	}

	var cResp *RequestDepositResponse
	resp, err := c.doRequest(RequestDepositEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	err = json.Unmarshal(respBytes, &cResp)
	if err != nil {
		return cResp, err
	}

	return cResp, nil
}

// GetBalance queries account balance
func (c *Client) GetBalance() (*BalanceResponse, error) {
	creds := c.auth.Authenticate()
	requestBody := GetBalanceBody{
		Username:  creds.Username,
		Timestamp: creds.Timestamp,
		AccountNo: c.AccountNo,
		Password:  creds.Password,
	}

	var cResp *BalanceResponse
	resp, err := c.doRequest(GetBalanceEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	err = json.Unmarshal(respBytes, &cResp)
	if err != nil {
		return cResp, err
	}

	return cResp, nil
}

// GetAuthCredentials returns the current authentication credentials.
// This is primarily useful for testing.
func (c *Client) GetAuthCredentials() Credentials {
	return c.auth.Authenticate()
}

// GetTransactionStatus queries the status of a transaction
func (c *Client) GetTransactionStatus(params *GetTransactionStatusParams) (*GetTransactionStatusResponse, error) {
	creds := c.auth.Authenticate()
	requestBody := GetTransactionStatusBody{
		Username:             creds.Username,
		Timestamp:            creds.Timestamp,
		RequestTransactionId: params.RequestTransactionId,
		TransactionId:        params.TransactionId,
		Password:             creds.Password,
	}

	var cResp *GetTransactionStatusResponse
	resp, err := c.doRequest(GetTransactionStatusEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	err = json.Unmarshal(respBytes, &cResp)
	if err != nil {
		return cResp, err
	}

	return cResp, nil
}
