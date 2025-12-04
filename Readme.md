# go-intouchpay

A Go package for integrating with the IntouchPay Payments Gateway API. This package provides a simple and easy-to-use interface for processing mobile money payments, deposits, balance inquiries, and transaction status checks.

## Features

- ✅ **Request Payment** - Initiate payment requests to subscribers
- ✅ **Request Deposit** - Send deposits to subscribers
- ✅ **Get Balance** - Query account balance
- ✅ **Get Transaction Status** - Check transaction status
- ✅ Full API v1.2 compliance
- ✅ Automatic password generation using SHA256
- ✅ UTC timestamp handling
- ✅ Form-encoded POST requests (as per API specification)

## Installation

```bash
go get github.com/samueltuyizere/go-intouchpay
```

## Prerequisites

Before using this package, you need to have:

1. An IntouchPay account with:
   - Username
   - Account Number
   - Partner Password
2. A callback URL (webhook endpoint) for receiving transaction status updates (optional but recommended)
3. Service ID (SID) - Set to `1` for Bulk Payments, `0` otherwise

## Quick Start

### 1. Create a Client

```go
package main

import (
    "fmt"
    "log"

    Intouchpay "github.com/samueltuyizere/go-intouchpay"
)

func main() {
    // Initialize the IntouchPay client
    client := Intouchpay.NewClient(
        "your_username",           // Username assigned to your account
        "your_account_number",     // Account number
        "your_partner_password",   // Partner password
        "https://yourdomain.com/callback", // Callback URL (optional)
        0,                         // Service ID (0 or 1)
    )

    // Use the client to make API calls...
}
```

### 2. Request Payment (Receive Payment)

Request a payment from a subscriber. The transaction will be pending until the subscriber confirms it.

```go
// Prepare payment parameters
params := &Intouchpay.RequestPaymentParams{
    Amount:               1000,                     // Amount to be paid (positive integer, no decimals)
    MobilePhone:          "250788888888",            // Mobile phone number making the payment
    RequestTransactionId: "unique_txn_id_12345",     // Unique transaction ID from your system
}

// Make the payment request
response, err := client.RequestPayment(params)
if err != nil {
    log.Fatal("Payment request failed:", err)
}

// Check the response
if response.Success {
    fmt.Printf("Transaction Status: %s\n", response.Status)
    fmt.Printf("Transaction ID: %s\n", response.TransactionId)
    fmt.Printf("Response Code: %s\n", response.ResponseCode)
    fmt.Printf("Message: %s\n", response.Message)
} else {
    fmt.Printf("Payment request failed: %s\n", response.Message)
}
```

**Response Example:**

```json
{
  "status": "Pending",
  "requesttransactionid": "unique_txn_id_12345",
  "success": true,
  "responsecode": "1000",
  "transactionid": "1425",
  "message": "Transaction Pending"
}
```

**Note:** After the subscriber confirms the transaction, IntouchPay will send a POST request to your callback URL with the final transaction status.

### 3. Request Deposit (Send Payment)

Send a deposit to a subscriber. This is processed immediately.

```go
// Prepare deposit parameters
params := &Intouchpay.RequestDepositParams{
    Amount:               5000,                     // Amount to deposit (positive integer, no decimals)
    WithdrawCharge:       0,                          // Set to 1 to include withdraw charges in amount
    Reason:               "Payment for services",     // Reason for deposit
    MobilePhone:          "250788888888",             // Mobile phone number receiving the deposit
    RequestTransactionId: "unique_deposit_id_67890", // Unique transaction ID from your system
}

// Make the deposit request
response, err := client.RequestDeposit(params)
if err != nil {
    log.Fatal("Deposit request failed:", err)
}

// Check the response
if response.Success {
    fmt.Printf("Deposit successful!\n")
    fmt.Printf("Reference ID: %s\n", response.ReferenceId)
    fmt.Printf("Response Code: %s\n", response.ResponseCode)
} else {
    fmt.Printf("Deposit failed: %s\n", response.ResponseCode)
}
```

**Response Example (Success):**

```json
{
  "requesttransactionid": "unique_deposit_id_67890",
  "referenceid": "312333883",
  "responsecode": "2001",
  "success": true
}
```

### 4. Get Balance

Query your account balance.

```go
// Get account balance
balance, err := client.GetBalance()
if err != nil {
    log.Fatal("Balance inquiry failed:", err)
}

if balance.Success {
    fmt.Printf("Account Balance: %s\n", balance.Balance)
} else {
    fmt.Printf("Failed to get balance: %s\n", balance.Message)
}
```

**Response Example:**

```json
{
  "balance": "100000.0",
  "success": true
}
```

### 5. Get Transaction Status

Check the status of a previously initiated transaction.

```go
// Prepare status check parameters
params := &Intouchpay.GetTransactionStatusParams{
    RequestTransactionId: "unique_txn_id_12345", // Your transaction ID
    TransactionId:        "1425",                // IntouchPay transaction ID
}

// Get transaction status
status, err := client.GetTransactionStatus(params)
if err != nil {
    log.Fatal("Status check failed:", err)
}

if status.Success {
    fmt.Printf("Transaction Status: %s\n", status.Status)
    fmt.Printf("Message: %s\n", status.Message)
    fmt.Printf("Response Code: %s\n", status.ResponseCode)
} else {
    fmt.Printf("Status check failed: %s\n", status.Message)
}
```

**Response Example:**

```json
{
  "success": true,
  "responsecode": "01",
  "status": "Successfully",
  "message": "Transaction Successful"
}
```

## Handling Callbacks (Webhooks)

When you initiate a payment request, IntouchPay will send a POST request to your callback URL with the transaction status. You need to implement an endpoint to receive these callbacks.

### Example Callback Handler

```go
package main

import (
    "encoding/json"
    "fmt"
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

    // Process the callback
    fmt.Printf("Transaction ID: %s\n", payload.JSONPayload.TransactionId)
    fmt.Printf("Status: %s\n", payload.JSONPayload.Status)
    fmt.Printf("Response Code: %s\n", payload.JSONPayload.ResponseCode)

    // Respond to IntouchPay
    response := map[string]interface{}{
        "message":   "success",
        "success":   true,
        "request_id": payload.JSONPayload.RequestTransactionId,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/callback", callbackHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Response Codes

### Payment Request Response Codes

| Code | Description                                       |
| ---- | ------------------------------------------------- |
| 1000 | Pending                                           |
| 01   | Successfully                                      |
| 0002 | Missing Username Information                      |
| 0003 | Missing Password Information                      |
| 0004 | Missing Date Information                          |
| 0005 | Invalid Password                                  |
| 0006 | User Does not have an intouchPay Account          |
| 0007 | No such user                                      |
| 0008 | Failed to Authenticate                            |
| 2100 | Amount should be greater than 0                   |
| 2200 | Amount below minimum                              |
| 2300 | Amount above maximum                              |
| 2400 | Duplicate Transaction ID                          |
| 2500 | Route Not Found                                   |
| 2600 | Operation Not Allowed                             |
| 2700 | Failed to Complete Transaction                    |
| 1005 | Failed Due to Insufficient Funds                  |
| 1002 | Mobile number not registered on mobile money      |
| 1008 | General Failure                                   |
| 1200 | Invalid Number                                    |
| 1100 | Number not supported on this Mobile money network |
| 1300 | Failed to Complete Transaction, Unknown Exception |

### Deposit Request Response Codes

| Code | Description                          |
| ---- | ------------------------------------ |
| 2001 | Request Successful                   |
| 1100 | Error in Request                     |
| 1101 | Service ID not Recognized            |
| 1102 | Invalid Mobile Phone Number          |
| 1103 | Payment Above Allowed Maximum        |
| 1104 | Payment Below Allowed Minimum        |
| 1105 | Network Not Supported                |
| 1106 | Operation Not Permitted              |
| 1107 | Payment Account Not Configured       |
| 1108 | Insufficient Account Balance         |
| 1110 | Duplicate Remit ID                   |
| 2102 | Subscriber Could not be Identified   |
| 2105 | Non Existent Mobile Account          |
| 2106 | Own Mobile Account Provided          |
| 2107 | Invalid Amount Format                |
| 2108 | Insufficient Funds on Source Account |
| 2109 | Daily Limit Exceeded                 |
| 2110 | Source Account Not Active            |
| 2111 | Mobile Account Not Active            |

### Transaction Status Response Codes

| Code | Description                                    |
| ---- | ---------------------------------------------- |
| 1000 | Transaction Pending                            |
| 01   | Transaction Successful for Payment Transaction |
| 2001 | Transaction Successful for Deposit Transaction |
| 3000 | Missing Transaction ID Information             |
| 3100 | Transaction Doesn't Exist                      |
| 3200 | Missing Request Transaction ID Information     |

## Error Handling

Always check for errors and response codes:

```go
response, err := client.RequestPayment(params)
if err != nil {
    // Handle network or API errors
    log.Printf("Error: %v", err)
    return
}

if !response.Success {
    // Handle API-level errors
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

// Success case
log.Println("Payment request successful")
```

## Configuration

### Service ID (SID)

The Service ID parameter determines the type of payment:

- `0` - Standard payments
- `1` - Bulk payments

### Withdraw Charge

When making a deposit, the `WithdrawCharge` parameter:

- `0` - Withdraw charges are NOT included in the amount sent
- `1` - Withdraw charges are included in the amount sent to the subscriber

### Timestamp Format

The package automatically generates timestamps in UTC format: `yyyymmddhhmmss` (e.g., `20161231115242`)

### Password Generation

The package automatically generates passwords using SHA256 encryption:

```text
password = SHA256(username + accountno + partnerpassword + timestamp)
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"

    Intouchpay "github.com/samueltuyizere/go-intouchpay"
)

func main() {
    // Initialize client
    client := Intouchpay.NewClient(
        "your_username",
        "your_account_number",
        "your_partner_password",
        "https://yourdomain.com/callback",
        0,
    )

    // Example 1: Request Payment
    paymentParams := &Intouchpay.RequestPaymentParams{
        Amount:               1000.0,
        MobilePhone:          "250788888888",
        RequestTransactionId: "txn_001",
    }

    paymentResp, err := client.RequestPayment(paymentParams)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Payment Response: %+v\n", paymentResp)

    // Example 2: Get Balance
    balance, err := client.GetBalance()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %s\n", balance.Balance)

    // Example 3: Request Deposit
    depositParams := &Intouchpay.RequestDepositParams{
        Amount:               5000.0,
        WithdrawCharge:       0,
        Reason:               "Payment for services",
        MobilePhone:          "250788888888",
        RequestTransactionId: "deposit_001",
    }

    depositResp, err := client.RequestDeposit(depositParams)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Deposit Response: %+v\n", depositResp)

    // Example 4: Get Transaction Status
    statusParams := &Intouchpay.GetTransactionStatusParams{
        RequestTransactionId: "txn_001",
        TransactionId:        paymentResp.TransactionId,
    }

    statusResp, err := client.GetTransactionStatus(statusParams)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %+v\n", statusResp)
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

## API Documentation

This package implements the IntouchPay API v1.2. For more details, refer to the official IntouchPay API documentation.

## License

See the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions:

- Open an issue on GitHub
- Contact IntouchPay support for API-related questions

## Version

Current version: **v0.1.2**
