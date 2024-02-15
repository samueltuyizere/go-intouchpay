package Intouchpay

import "net/http"

const (
	RequestPaymentEndpoint = "/requestpayment/"
	RequestDepositEndpoint = "/requestdeposit/"
	GetBalanceEndpoint     = "/getbalance/"
	BaseUrl                = "https://www.intouchpay.co.rw/api"
)

// Client represents an IntouchPay client configured with authentication details
type Client struct {
	Username        string
	AccountNo       string
	PartnerPassword string
	CallbackURL     string
	Sid             int // can only be 0 or 1
	HTTPClient      *http.Client
}

type FailedRequestResponse struct {
	Success      bool   `json:"success"`
	ResponseCode string `json:"responsecode"`
	Message      string `json:"message"`
}

type RequestPaymentParams struct {
	Amount               string `json:"amount"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
	Reason               string `json:"reason"` // optional, the reason for the payment being made
}

type RequestPaymentResponse struct {
	Status               string `json:"status"`
	RequestTransactionId string `json:"requesttransactionid"`
	Success              bool   `json:"success"`
	ResponseCode         string `json:"responsecode"`
	TransactionId        int    `json:"transactionid"`
	Message              string `json:"message"`
}

type BalanceResponse struct {
	Balance string `json:"balance"`
	Succes  bool   `json:"success"`
}

type RequestDepositParams struct {
	Amount               int    `json:"amount"`
	WithdrawCharge       int    `json:"withdrawcharge"`
	Reason               string `json:"reason"`
	Sid                  string `json:"sid"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

type RequestDepositResponse struct {
	RequestTransactionId string `json:"requesttransactionid"`
	ReferenceId          string `json:"referenceid"`
	ResponseCode         string `json:"responsecode"`
	Success              int    `json:"success"`
}
