package Intouchpay

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/samueltuyizere/validate_rw_phone_numbers"
)

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
	Amount               uint   `json:"amount"`
	Password             string `json:"password"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
	AccountNo            string `json:"accountno"`
	CallbackURL          string `json:"callbackurl,omitempty"`
}

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

type GetBalanceBody struct {
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
	AccountNo string `json:"accountno"`
	Password  string `json:"password"`
}

type GetTransactionStatusBody struct {
	Username             string `json:"username"`
	Timestamp            string `json:"timestamp"`
	RequestTransactionId string `json:"requesttransactionid"`
	TransactionId        string `json:"transactionid"`
	Password             string `json:"password"`
}

type BalanceResponse struct {
	Balance      float64 `json:"balance"`
	Success      bool    `json:"success"`
	ResponseCode int     `json:"responsecode,omitempty"`
	Message      string  `json:"message,omitempty"`
}

type RequestDepositParams struct {
	Amount               uint   `json:"amount"`         // Amount as a positive integer with no decimals
	WithdrawCharge       int    `json:"withdrawcharge"` // Set to 1 to include Withdraw Charges in amount sent to subscriber
	Reason               string `json:"reason"`
	MobilePhone          string `json:"mobilephone"`
	RequestTransactionId string `json:"requesttransactionid"`
}

type RequestDepositResponse struct {
	RequestTransactionId string `json:"requesttransactionid"`
	ReferenceId          string `json:"referenceid,omitempty"` // Only returned if successful
	ResponseCode         int    `json:"responsecode"`
	Success              bool   `json:"success"`
}

type GetTransactionStatusParams struct {
	RequestTransactionId string `json:"requesttransactionid"`
	TransactionId        string `json:"transactionid"`
}

type GetTransactionStatusResponse struct {
	Success      bool   `json:"success"`
	ResponseCode int    `json:"responsecode"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message"`
}

func SanitizePhoneNumber(phoneNumber string) (string, error) {
	// Remove any existing country code prefix
	cleanedNumber := phoneNumber
	if len(phoneNumber) > 3 && phoneNumber[:3] == "250" {
		return cleanedNumber, nil
	}

	// Validate the number without country code
	isValidMtn := validate_rw_phone_numbers.ValidateMtn(cleanedNumber)
	isValidAirtelTigo := validate_rw_phone_numbers.ValidateAirtelTigo(cleanedNumber)
	if !isValidMtn && !isValidAirtelTigo {
		return "", errors.New("invalid phone number")
	}

	// Return with country code prefix
	newNumber := fmt.Sprintf("25%s", cleanedNumber)
	return newNumber, nil
}
