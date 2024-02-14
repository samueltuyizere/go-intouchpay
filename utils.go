package intouchpay

import "time"

const (
	RequestPaymentEndpoint = "/requestpayment/"
	RequestDepositEndpoint = "/requestdeposit/"
	GetBalanceEndpoint     = "/getbalance/"
	BaseUrl                = "https://www.intouchpay.co.rw/api"
)

type FailedRequestResponse struct {
	Success      bool   `json:"success"`
	ResponseCode string `json:"responsecode"`
	Message      string `json:"message"`
}

type RequestPaymentParams struct {
	UserName             string    `json:"username"`
	Amount               string    `json:"amount"`
	Password             string    `json:"password"`
	MobilePhone          string    `json:"mobilephone"`
	RequestTransactionId string    `json:"requesttransactionid"`
	Timestamp            time.Time `json:"timestamp"`
}

type RequestPaymentResponse struct {
	Status               string `json:"status"`
	RequestTransactionId string `json:"requesttransactionid"`
	Success              bool   `json:"success"`
	ResponseCode         string `json:"responsecode"`
	TransactionId        int    `json:"transactionid"`
	Message              string `json:"message"`
}

type RequestBalanceParams struct {
	UserName  string `json:"username"`
	Timestamp string `json:"timestamp"`
	Password  string `json:"password"`
}

type BalanceResponse struct {
	Balance string `json:"balance"`
	Succes  bool   `json:"success"`
}

type RequestDepositParams struct {
	Username             string `json:"username"`
	Timestamp            string `json:"timestamp"`
	Aount                int    `json:"amount"`
	WithdrawCharge       int    `json:"withdrawcharge"`
	Reason               string `json:"reason"`
	Sid                  string `json:"sid"`
	Password             string `json:"password"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

type RequestDepositResponse struct {
	RequestTransactionId string `json:"requesttransactionid"`
	ReferenceId          string `json:"referenceid"`
	ResponseCode         string `json:"responsecode"`
	Succes               int    `json:"success"`
}
