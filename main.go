package Intouchpay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
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
		HTTPClient:      &http.Client{Timeout: 60 * time.Second}, // API responds within 60 seconds per documentation
	}
}

// generatePassword calculates the SHA256 hash of the password string
func (c *Client) generatePassword(timestamp string) string {
	passwordString := c.Username + c.AccountNo + c.PartnerPassword + timestamp
	hash := sha256.New()
	hash.Write([]byte(passwordString))
	hashInBytes := hash.Sum(nil)
	password := hex.EncodeToString(hashInBytes)
	return password
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
		return response, fmt.Errorf("IntouchPay API error: %d\n %s\n %v", resp.StatusCode, resp.Status, *response)
	}
	return response, nil
}

// RequestPayment initiates a payment request
func (c *Client) RequestPayment(params *RequestPaymentParams) (*RequestPaymentResponse, error) {
	phoneNumber, err := SanitizePhoneNumber(params.MobilePhone)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	requestBody := RequestPaymentBody{
		Username:             c.Username,
		Timestamp:            timestamp,
		Amount:               params.Amount,
		Password:             password,
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
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	requestBody := RequestDepositBody{
		Username:             c.Username,
		Timestamp:            timestamp,
		Amount:               params.Amount,
		WithdrawCharge:       params.WithdrawCharge,
		Reason:               params.Reason,
		Sid:                  c.Sid,
		Password:             password,
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
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	requestBody := GetBalanceBody{
		Username:  c.Username,
		Timestamp: timestamp,
		AccountNo: c.AccountNo,
		Password:  password,
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

// GetTransactionStatus queries the status of a transaction
func (c *Client) GetTransactionStatus(params *GetTransactionStatusParams) (*GetTransactionStatusResponse, error) {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	requestBody := GetTransactionStatusBody{
		Username:             c.Username,
		Timestamp:            timestamp,
		RequestTransactionId: params.RequestTransactionId,
		TransactionId:        params.TransactionId,
		Password:             password,
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
