# AI Agent Integration Guide

This document provides instructions for AI coding agents (Claude, GPT-4, GitHub Copilot, etc.) to integrate the `go-intouchpay` package into a Go project.

## Package Overview

**Package:** `github.com/samueltuyizere/go-intouchpay`  
**Purpose:** IntouchPay Payments Gateway integration for mobile money payments in Rwanda  
**Go Version:** 1.16+  
**Architecture:** Deep module design with interfaces for testability

## Integration Checklist

When a user asks you to integrate this package, follow these steps:

### 1. Install the Package

```bash
go get github.com/samueltuyizere/go-intouchpay
```

### 2. Gather Required Credentials

Ask the user for the following IntouchPay credentials:

| Credential | Description | Required |
|------------|-------------|----------|
| Username | IntouchPay account username | ✅ Yes |
| Account Number | IntouchPay account number | ✅ Yes |
| Partner Password | Partner password for authentication | ✅ Yes |
| Callback URL | Webhook URL for transaction status updates | Recommended |
| Service ID (SID) | `0` for standard payments, `1` for bulk payments | ✅ Yes |

### 3. Choose Integration Pattern

#### Option A: Simple Integration (Recommended for most cases)

```go
package main

import (
    "fmt"
    "log"

    Intouchpay "github.com/samueltuyizere/go-intouchpay"
)

func main() {
    client := Intouchpay.NewClient(
        "USERNAME",        // IntouchPay username
        "ACCOUNT_NUMBER",  // IntouchPay account number
        "PARTNER_PASSWORD", // Partner password
        "https://yourdomain.com/callback", // Callback URL
        0,                 // Service ID (0 or 1)
    )

    // Use the client...
}
```

#### Option B: Options Pattern (For custom configuration)

```go
import "time"

client := Intouchpay.NewClientWithOptions(
    "USERNAME",
    "ACCOUNT_NUMBER",
    "PARTNER_PASSWORD",
    Intouchpay.WithCallbackURL("https://yourdomain.com/callback"),
    Intouchpay.WithSid(0),
    Intouchpay.WithTimeout(30 * time.Second),
)
```

### 4. Implement Required Operations

#### Request Payment (Receive money from subscriber)

```go
params := &Intouchpay.RequestPaymentParams{
    Amount:               1000,              // Amount in RWF
    MobilePhone:          "250788888888",    // Rwandan phone number
    RequestTransactionId: "unique_txn_id",   // Your unique transaction ID
}

response, err := client.RequestPayment(params)
if err != nil {
    log.Fatal(err)
}

if response.Success {
    fmt.Printf("Transaction ID: %s\n", response.TransactionId)
    fmt.Printf("Status: %s\n", response.Status)
}
```

#### Request Deposit (Send money to subscriber)

```go
params := &Intouchpay.RequestDepositParams{
    Amount:               5000,
    WithdrawCharge:       0,                    // 0 or 1
    Reason:               "Payment for services",
    MobilePhone:          "250788888888",
    RequestTransactionId: "unique_deposit_id",
}

response, err := client.RequestDeposit(params)
if err != nil {
    log.Fatal(err)
}

if response.Success {
    fmt.Printf("Reference ID: %s\n", response.ReferenceId)
}
```

#### Get Balance

```go
balance, err := client.GetBalance()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Balance: %s\n", balance.Balance)
```

#### Get Transaction Status

```go
params := &Intouchpay.GetTransactionStatusParams{
    RequestTransactionId: "your_txn_id",
    TransactionId:        "intouchpay_txn_id", // Optional
}

status, err := client.GetTransactionStatus(params)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status)
```

### 5. Implement Callback Handler (Webhook)

IntouchPay sends POST requests to your callback URL with transaction status updates.

```go
package main

import (
    "encoding/json"
    "net/http"
)

type CallbackPayload struct {
    JSONPayload struct {
        RequestTransactionId string `json:"requesttransactionid"`
        TransactionId        string `json:"transactionid"`
        ResponseCode         string `json:"responsecode"`
        Status               string `json:"status"`
        StatusDesc           string `json:"statusdesc"`
        ReferenceNo          string `json:"referenceno"`
    } `json:"jsonpayload"`
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var payload CallbackPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Process the callback based on status
    // "Successfully" = transaction completed
    // "Failed" = transaction failed

    // Respond to IntouchPay
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":    "success",
        "success":    true,
        "request_id": payload.JSONPayload.RequestTransactionId,
    })
}

func main() {
    http.HandleFunc("/callback", callbackHandler)
    http.ListenAndServe(":8080", nil)
}
```

### 6. Handle Errors Properly

```go
import Intouchpay "github.com/samueltuyizere/go-intouchpay"

response, err := client.RequestPayment(params)
if err != nil {
    // Check for typed errors
    if Intouchpay.IsAPIError(err) {
        apiErr := err.(*Intouchpay.APIError)
        log.Printf("API Error: %d - %s", apiErr.StatusCode, apiErr.Status)
    }
    
    if Intouchpay.IsValidationError(err) {
        valErr := err.(*Intouchpay.ValidationError)
        log.Printf("Validation Error: %s - %s", valErr.Field, valErr.Message)
    }
    
    return
}

// Handle API-level errors
if !response.Success {
    switch response.ResponseCode {
    case "2400":
        log.Println("Duplicate transaction ID")
    case "1005":
        log.Println("Insufficient funds")
    case "1002":
        log.Println("Mobile number not registered")
    default:
        log.Printf("Payment failed: %s", response.Message)
    }
    return
}
```

## Phone Number Format

The package handles Rwandan phone number sanitization:

- Input: `0788888888`, `+250788888888`, `250788888888`
- All formats are normalized to: `250788888888`

Valid prefixes: `25072`, `25073`, `25078`, `25079`

## Common Response Codes

| Code | Description |
|------|-------------|
| 1000 | Pending (waiting for user confirmation) |
| 01   | Success (payment transaction) |
| 2001 | Success (deposit transaction) |
| 2400 | Duplicate transaction ID |
| 1005 | Insufficient funds |
| 1002 | Mobile number not registered |

## Testing Integration

For testing without real credentials:

```go
// Create mock authenticator
type MockAuthenticator struct {
    Creds Intouchpay.Credentials
}

func (m *MockAuthenticator) Authenticate() Intouchpay.Credentials {
    return m.Creds
}

// Create mock HTTP client
type MockHTTPClient struct {
    Response *map[string]interface{}
    Error    error
}

func (m *MockHTTPClient) Do(endpoint string, body interface{}) (*map[string]interface{}, error) {
    return m.Response, m.Error
}

// Use mocks
mockAuth := &MockAuthenticator{
    Creds: Intouchpay.Credentials{
        Username:  "test_user",
        Timestamp: "20260320120000",
        Password:  "test_hash",
    },
}

mockHTTP := &MockHTTPClient{
    Response: &map[string]interface{}{
        "success": true,
        "balance": 10000.50,
    },
}

client := Intouchpay.NewClientWithHTTPClient(mockAuth, mockHTTP)
```

## Environment Variables (Recommended)

Store credentials securely using environment variables:

```go
import (
    "os"
    "time"
    
    Intouchpay "github.com/samueltuyizere/go-intouchpay"
)

func main() {
    client := Intouchpay.NewClientWithOptions(
        os.Getenv("INTOUCHPAY_USERNAME"),
        os.Getenv("INTOUCHPAY_ACCOUNT_NUMBER"),
        os.Getenv("INTOUCHPAY_PARTNER_PASSWORD"),
        Intouchpay.WithCallbackURL(os.Getenv("INTOUCHPAY_CALLBACK_URL")),
        Intouchpay.WithSid(0),
        Intouchpay.WithTimeout(60 * time.Second),
    )
}
```

## Complete Integration Example

Here's a complete example integrating all features:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    Intouchpay "github.com/samueltuyizere/go-intouchpay"
)

var client *Intouchpay.Client

func main() {
    // Initialize client
    client = Intouchpay.NewClientWithOptions(
        os.Getenv("INTOUCHPAY_USERNAME"),
        os.Getenv("INTOUCHPAY_ACCOUNT_NUMBER"),
        os.Getenv("INTOUCHPAY_PARTNER_PASSWORD"),
        Intouchpay.WithCallbackURL(os.Getenv("INTOUCHPAY_CALLBACK_URL")),
        Intouchpay.WithSid(0),
        Intouchpay.WithTimeout(60 * time.Second),
    )

    // Set up routes
    http.HandleFunc("/payment", requestPaymentHandler)
    http.HandleFunc("/deposit", requestDepositHandler)
    http.HandleFunc("/balance", getBalanceHandler)
    http.HandleFunc("/callback", callbackHandler)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func requestPaymentHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Amount     int    `json:"amount"`
        Phone      string `json:"phone"`
        TxnID      string `json:"transaction_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    params := &Intouchpay.RequestPaymentParams{
        Amount:               req.Amount,
        MobilePhone:          req.Phone,
        RequestTransactionId: req.TxnID,
    }

    response, err := client.RequestPayment(params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func requestDepositHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Amount     int    `json:"amount"`
        Phone      string `json:"phone"`
        TxnID      string `json:"transaction_id"`
        Reason     string `json:"reason"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    params := &Intouchpay.RequestDepositParams{
        Amount:               req.Amount,
        MobilePhone:          req.Phone,
        RequestTransactionId: req.TxnID,
        Reason:               req.Reason,
        WithdrawCharge:       0,
    }

    response, err := client.RequestDeposit(params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
    balance, err := client.GetBalance()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(balance)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var payload struct {
        JSONPayload struct {
            RequestTransactionId string `json:"requesttransactionid"`
            TransactionId        string `json:"transactionid"`
            ResponseCode         string `json:"responsecode"`
            Status               string `json:"status"`
            StatusDesc           string `json:"statusdesc"`
        } `json:"jsonpayload"`
    }

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Log and process the callback
    log.Printf("Callback received: TransactionID=%s, Status=%s, ResponseCode=%s",
        payload.JSONPayload.TransactionId,
        payload.JSONPayload.Status,
        payload.JSONPayload.ResponseCode,
    )

    // Update your database with the transaction status
    // ...

    // Respond to IntouchPay
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":    "success",
        "success":    true,
        "request_id": payload.JSONPayload.RequestTransactionId,
    })
}
```

## Quick Reference

| Method | Description |
|--------|-------------|
| `NewClient()` | Simple constructor with all parameters |
| `NewClientWithOptions()` | Flexible constructor with options |
| `RequestPayment()` | Request payment from subscriber |
| `RequestDeposit()` | Send deposit to subscriber |
| `GetBalance()` | Query account balance |
| `GetTransactionStatus()` | Check transaction status |

## Support

- GitHub Issues: https://github.com/samueltuyizere/go-intouchpay/issues
- Main Documentation: See README.md in the repository root
