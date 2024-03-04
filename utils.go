package Intouchpay

import "net/http"

const (
	RequestPaymentEndpoint string = "/requestpayment/"
	RequestDepositEndpoint string = "/requestdeposit/"
	GetBalanceEndpoint     string = "/getbalance"
	BaseUrl                string = "https://www.intouchpay.co.rw/api"
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
	Amount               int    `json:"amount"`
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
	Amount               int    `json:"amount"`
	Password             string `json:"password"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
	CallbackURL          string `json:"callbackurl"`
}

type BalanceResponse struct {
	Balance string `json:"balance"`
	Succes  bool   `json:"success"`
}

type RequestDepositParams struct {
	Amount               int    `json:"amount"`
	WithdrawCharge       int    `json:"withdrawcharge"` // Set to 1 to include Withdraw Charges in amount sent to subscriber
	Reason               string `json:"reason"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

type RequestDepositResponse struct {
	RequestTransactionId string `json:"requesttransactionid"`
	ReferenceId          string `json:"referenceid"`
	ResponseCode         string `json:"responsecode"`
	Success              bool   `json:"success"`
}

type RequestData interface {
	*RequestPaymentBody
}
