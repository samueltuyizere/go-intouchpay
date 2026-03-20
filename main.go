package Intouchpay

import (
	"encoding/json"
	"net/http"
)

// NewClient creates a new IntouchPay client with the provided credentials.
// This is a convenience constructor that creates a default authenticator.
func NewClient(username, accountNumber, partnerPassword, callbackURL string, sid int) *Client {
	auth := NewAuthenticator(username, accountNumber, partnerPassword)
	httpClient := &http.Client{Timeout: DefaultTimeout}
	c := &Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		CallbackURL:     callbackURL,
		Sid:             sid,
		auth:            auth,
		HTTPClient:      httpClient,
		httpClient:      NewHTTPClient(httpClient, BaseURL),
	}
	return c
}

// NewClientWithAuth creates a new IntouchPay client with a custom authenticator.
// Use this for testing or when you need custom authentication behavior.
func NewClientWithAuth(auth Authenticator, opts ...Option) *Client {
	httpClient := &http.Client{Timeout: DefaultTimeout}
	c := &Client{
		auth:       auth,
		HTTPClient: httpClient,
		httpClient: NewHTTPClient(httpClient, BaseURL),
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
	httpClient := &http.Client{Timeout: DefaultTimeout}
	c := &Client{
		Username:        username,
		AccountNo:       accountNumber,
		PartnerPassword: partnerPassword,
		auth:            auth,
		HTTPClient:      httpClient,
		httpClient:      NewHTTPClient(httpClient, BaseURL),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewClientWithHTTPClient creates a new IntouchPay client with a custom APIRequester interface.
// Use this for testing with a mock HTTP client.
func NewClientWithHTTPClient(auth Authenticator, httpClient APIRequester, opts ...Option) *Client {
	c := &Client{
		auth:       auth,
		httpClient: httpClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
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
		RequestTransactionID: params.RequestTransactionID,
		AccountNo:            c.AccountNo,
	}
	if c.CallbackURL != "" {
		requestBody.CallbackURL = c.CallbackURL
	}

	var cResp *RequestPaymentResponse
	resp, err := c.httpClient.Do(RequestPaymentEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, NewMarshalError("request payment response", err)
	}
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
		RequestTransactionID: params.RequestTransactionID,
		AccountNo:            c.AccountNo,
	}

	var cResp *RequestDepositResponse
	resp, err := c.httpClient.Do(RequestDepositEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, NewMarshalError("request deposit response", err)
	}
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
	resp, err := c.httpClient.Do(GetBalanceEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, NewMarshalError("balance response", err)
	}
	err = json.Unmarshal(respBytes, &cResp)
	if err != nil {
		return cResp, err
	}

	return cResp, nil
}

// GetTransactionStatus queries the status of a transaction
func (c *Client) GetTransactionStatus(params *GetTransactionStatusParams) (*GetTransactionStatusResponse, error) {
	creds := c.auth.Authenticate()
	requestBody := GetTransactionStatusBody{
		Username:             creds.Username,
		Timestamp:            creds.Timestamp,
		RequestTransactionID: params.RequestTransactionID,
		TransactionID:        params.TransactionID,
		Password:             creds.Password,
	}

	var cResp *GetTransactionStatusResponse
	resp, err := c.httpClient.Do(GetTransactionStatusEndpoint, requestBody)
	if err != nil {
		return cResp, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, NewMarshalError("transaction status response", err)
	}
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
