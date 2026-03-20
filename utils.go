package Intouchpay

import "net/http"

// API endpoint constants
const (
	RequestPaymentEndpoint       string = "/requestpayment/"
	RequestDepositEndpoint       string = "/requestdeposit/"
	GetBalanceEndpoint           string = "/getbalance/"
	GetTransactionStatusEndpoint string = "/gettransactionstatus/"
	BaseUrl                      string = "https://www.intouchpay.co.rw/api"
)

// Client represents an IntouchPay client configured with authentication details
type Client struct {
	Username        string // User name assigned to your account
	AccountNo       string
	PartnerPassword string
	CallbackURL     string
	Sid             int // Service ID. Set to 1 For Bulk Payments, can only be 0 or 1
	HTTPClient      *http.Client // Kept for backward compatibility
	auth            Authenticator
	httpClient      HTTPClient // Internal HTTP client interface
}

// FailedRequestResponse represents a failed API response
type FailedRequestResponse struct {
	Success      bool   `json:"success"`
	ResponseCode string `json:"responsecode"`
	Message      string `json:"message"`
}

// RequestPaymentParams represents parameters for RequestPayment
type RequestPaymentParams struct {
	Amount               uint   `json:"amount"` // Amount as a positive integer with no decimals
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

// RequestPaymentResponse represents the response from RequestPayment
type RequestPaymentResponse struct {
	Status               string `json:"status"`
	RequestTransactionId string `json:"requesttransactionid"`
	Success              bool   `json:"success"`
	ResponseCode         string `json:"responsecode"`
	TransactionId        string `json:"transactionid"`
	Message              string `json:"message"`
}

// RequestPaymentBody represents the request body for RequestPayment
type RequestPaymentBody struct {
	Username             string `json:"username"`
	Timestamp            string `json:"timestamp"`
	Amount               uint   `json:"amount"`
	Password             string `json:"password"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
	AccountNo            string `json:"accountno"`
	CallbackURL          string `json:"callbackurl,omitempty"`
}

// RequestDepositBody represents the request body for RequestDeposit
type RequestDepositBody struct {
	Username             string `json:"username"`
	Timestamp            string `json:"timestamp"`
	Amount               uint   `json:"amount"`
	WithdrawCharge       int    `json:"withdrawcharge"`
	Reason               string `json:"reason"`
	Sid                  int    `json:"sid"`
	Password             string `json:"password"`
	MobilePhoneNo        string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
	AccountNo            string `json:"accountno"`
}

// GetBalanceBody represents the request body for GetBalance
type GetBalanceBody struct {
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
	AccountNo string `json:"accountno"`
	Password  string `json:"password"`
}

// GetTransactionStatusBody represents the request body for GetTransactionStatus
type GetTransactionStatusBody struct {
	Username             string `json:"username"`
	Timestamp            string `json:"timestamp"`
	RequestTransactionId string `json:"requesttransactionid"`
	TransactionId        string `json:"transactionid"`
	Password             string `json:"password"`
}

// BalanceResponse represents the response from GetBalance
type BalanceResponse struct {
	Balance      float64 `json:"balance"`
	Success      bool    `json:"success"`
	ResponseCode int     `json:"responsecode,omitempty"`
	Message      string  `json:"message,omitempty"`
}

// RequestDepositParams represents parameters for RequestDeposit
type RequestDepositParams struct {
	Amount               uint   `json:"amount"`         // Amount as a positive integer with no decimals
	WithdrawCharge       int    `json:"withdrawcharge"` // Set to 1 to include Withdraw Charges in amount sent to subscriber
	Reason               string `json:"reason"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

// RequestDepositResponse represents the response from RequestDeposit
type RequestDepositResponse struct {
	RequestTransactionId string `json:"requesttransactionid"`
	ReferenceId          string `json:"referenceid,omitempty"` // Only returned if successful
	ResponseCode         string `json:"responsecode"`
	Success              bool   `json:"success"`
}

// GetTransactionStatusParams represents parameters for GetTransactionStatus
type GetTransactionStatusParams struct {
	RequestTransactionId string `json:"requesttransactionid"`
	TransactionId        string `json:"transactionid"`
}

// GetTransactionStatusResponse represents the response from GetTransactionStatus
type GetTransactionStatusResponse struct {
	Success      bool   `json:"success"`
	ResponseCode int    `json:"responsecode"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message"`
}
