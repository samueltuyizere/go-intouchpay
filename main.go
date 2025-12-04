package Intouchpay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// doRequest sends an HTTP form POST request and handles the response
func (c *Client) doRequest(endpoint string, formData url.Values) (*map[string]interface{}, error) {
	var response *map[string]interface{}
	requestUrl := BaseUrl + endpoint
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	if errr := json.NewDecoder(resp.Body).Decode(&response); errr != nil {
		return response, fmt.Errorf("IntouchPay API error: %d\n %s\n %v", resp.StatusCode, resp.Status, errr)
	}
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("IntouchPay API error: %d\n %s\n %v", resp.StatusCode, resp.Status, *response)
	}
	return response, nil
}

// RequestPayment initiates a payment request
func (c *Client) RequestPayment(params *RequestPaymentParams) (*RequestPaymentResponse, error) {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	formData := url.Values{}
	formData.Set("username", c.Username)
	formData.Set("timestamp", timestamp)
	formData.Set("amount", strconv.FormatUint(uint64(params.Amount), 10))
	formData.Set("password", password)
	formData.Set("mobilephoneno", params.MobilePhone)
	formData.Set("requesttransactionid", params.RequestTransactionId)
	formData.Set("accountno", c.AccountNo)
	if c.CallbackURL != "" {
		formData.Set("callbackurl", c.CallbackURL)
	}

	var cResp *RequestPaymentResponse
	resp, err := c.doRequest(RequestPaymentEndpoint, formData)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	errr := json.Unmarshal(respBytes, &cResp)
	if errr != nil {
		return cResp, errr
	}

	return cResp, nil
}

// RequestDeposit initiates a deposit request
func (c *Client) RequestDeposit(params *RequestDepositParams) (*RequestDepositResponse, error) {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	formData := url.Values{}
	formData.Set("username", c.Username)
	formData.Set("timestamp", timestamp)
	formData.Set("amount", strconv.FormatUint(uint64(params.Amount), 10))
	formData.Set("withdrawcharge", strconv.Itoa(params.WithdrawCharge))
	formData.Set("reason", params.Reason)
	formData.Set("sid", strconv.Itoa(c.Sid))
	formData.Set("password", password)
	formData.Set("mobilephoneno", params.MobilePhone)
	formData.Set("requesttransactionid", params.RequestTransactionId)
	formData.Set("accountno", c.AccountNo)

	var cResp *RequestDepositResponse
	resp, err := c.doRequest(RequestDepositEndpoint, formData)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	errr := json.Unmarshal(respBytes, &cResp)
	if errr != nil {
		return cResp, errr
	}

	return cResp, nil
}

// GetBalance queries account balance
func (c *Client) GetBalance() (*BalanceResponse, error) {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	formData := url.Values{}
	formData.Set("username", c.Username)
	formData.Set("timestamp", timestamp)
	formData.Set("accountno", c.AccountNo)
	formData.Set("password", password)

	var cResp *BalanceResponse
	resp, err := c.doRequest(GetBalanceEndpoint, formData)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	errr := json.Unmarshal(respBytes, &cResp)
	if errr != nil {
		return cResp, errr
	}

	return cResp, nil
}

// GetTransactionStatus queries the status of a transaction
func (c *Client) GetTransactionStatus(params *GetTransactionStatusParams) (*GetTransactionStatusResponse, error) {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)

	formData := url.Values{}
	formData.Set("username", c.Username)
	formData.Set("timestamp", timestamp)
	formData.Set("requesttransactionid", params.RequestTransactionId)
	formData.Set("transactionid", params.TransactionId)
	formData.Set("password", password)

	var cResp *GetTransactionStatusResponse
	resp, err := c.doRequest(GetTransactionStatusEndpoint, formData)
	if err != nil {
		return cResp, err
	}
	respBytes, _ := json.Marshal(resp)
	errr := json.Unmarshal(respBytes, &cResp)
	if errr != nil {
		return cResp, errr
	}

	return cResp, nil
}
