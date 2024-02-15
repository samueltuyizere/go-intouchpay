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
		HTTPClient:      &http.Client{},
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

// doRequest sends an HTTP request and handles the response
func (c *Client) doRequest(method, endpoint string, data interface{}) (*interface{}, error) {
	var response *interface{}
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(data)
	requestUrl := BaseUrl + endpoint
	req, _ := http.NewRequest(method, requestUrl, body)
	req.Header.Set("Content-Type", "application/json")
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
	now := time.Now()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)
	var cResp *RequestPaymentResponse
	body := RequestPaymentBody{
		Username:             c.Username,
		Timestamp:            timestamp,
		Amount:               params.Amount,
		Password:             password,
		MobilePhone:          params.MobilePhone,
		RequestTransactionId: params.RequestTransactionId,
		CallbackURL:          c.CallbackURL,
	}
	resp, err := c.doRequest(http.MethodPost, RequestPaymentEndpoint, body)
	if err != nil {
		fmt.Printf("%v", err)
		return cResp, err
	}
	fmt.Printf("%v", &resp)
	return cResp, nil
}

// RequestDeposit initiates a deposit request
func (c *Client) RequestDeposit(params *RequestDepositParams) (*RequestDepositResponse, error) {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	password := c.generatePassword(timestamp)
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
	fmt.Printf("%v", resp)
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
	fmt.Printf("%v", resp)
	return cResp, nil
}
