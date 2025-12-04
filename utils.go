package Intouchpay

import "net/http"

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
	HTTPClient      *http.Client
}

type FailedRequestResponse struct {
	Success      bool   `json:"success"`
	ResponseCode string `json:"responsecode"`
	Message      string `json:"message"`
}

type RequestPaymentParams struct {
	Amount               uint   `json:"amount"` // Amount as a positive integer with no decimals
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

type RequestPaymentResponse struct {
	Status               string `json:"status"`
	RequestTransactionId string `json:"requesttransactionid"`
	Success              bool   `json:"success"`
	ResponseCode         string `json:"responsecode"`
	TransactionId        string `json:"transactionid"`
	Message              string `json:"message"`
}

type RequestPaymentBody struct {
	Username             string `json:"username"`
	Timestamp            string `json:"timestamp"`
	Amount               string `json:"amount"`
	Password             string `json:"password"`
	MobilePhoneNo        string `json:"mobilephoneno"`
	RequestTransactionId string `json:"requesttransactionid"`
	AccountNo            string `json:"accountno"`
	CallbackURL          string `json:"callbackurl,omitempty"`
}

type BalanceResponse struct {
	Balance      string `json:"balance"`
	Success      bool   `json:"success"`
	ResponseCode string `json:"responsecode,omitempty"`
	Message      string `json:"message,omitempty"`
}

type RequestDepositParams struct {
	Amount               uint   `json:"amount"` // Amount as a positive integer with no decimals
	WithdrawCharge       int    `json:"withdrawcharge"` // Set to 1 to include Withdraw Charges in amount sent to subscriber
	Reason               string `json:"reason"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

type RequestDepositResponse struct {
	RequestTransactionId string `json:"requesttransactionid"`
	ReferenceId          string `json:"referenceid,omitempty"` // Only returned if successful
	ResponseCode         string `json:"responsecode"`
	Success              bool   `json:"success"`
}

type GetTransactionStatusParams struct {
	RequestTransactionId string `json:"requesttransactionid"`
	TransactionId        string `json:"transactionid"`
}

type GetTransactionStatusResponse struct {
	Success      bool   `json:"success"`
	ResponseCode string `json:"responsecode"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message"`
}
