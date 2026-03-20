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
- ✅ **Testable architecture** - Mockable authentication and HTTP layers
- ✅ **Option pattern** - Flexible client configuration
- ✅ **Typed errors** - APIError and ValidationError types

## AI Agent Integration

This package includes special documentation for AI coding assistants (Claude, GPT-4, GitHub Copilot, Cursor, etc.).

### For AI Agents

If you are an AI coding agent tasked with integrating this package, please read the detailed integration guide:

📄 **[AI Agent Integration Guide](docs/AI_AGENT_INTEGRATION.md)**

This guide contains:
- Complete integration checklist
- Required credentials and how to request them
- Code examples for all API operations
- Callback/webhook implementation
- Error handling patterns
- Testing with mocks
- Complete working example

### Prompt for Your AI Coding Assistant

Copy and paste this prompt to your AI coding assistant to integrate go-intouchpay into your project:

```
I need you to integrate the go-intouchpay payment gateway package into my Go project.

Please follow the instructions in this document:
https://raw.githubusercontent.com/samueltuyizere/go-intouchpay/v2/docs/AI_AGENT_INTEGRATION.md

The package documentation is at:
https://github.com/samueltuyizere/go-intouchpay

Before you start, I'll provide my IntouchPay credentials:
- Username: [YOUR_USERNAME]
- Account Number: [YOUR_ACCOUNT_NUMBER]  
- Partner Password: [YOUR_PARTNER_PASSWORD]
- Callback URL: [YOUR_CALLBACK_URL]
- Service ID (SID): [0 or 1]

Please:
1. Install the package
2. Create a client with my credentials
3. Implement [payment request / deposit / balance / all operations]
4. Set up the callback webhook handler
5. Add proper error handling
6. Use environment variables for credentials (not hardcoded values)
```

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

#### Basic Usage

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

#### With Options Pattern

For more flexible configuration, use the options pattern:

```go
client := Intouchpay.NewClientWithOptions(
    "your_username",
    "your_account_number",
    "your_partner_password",
    Intouchpay.WithCallbackURL("https://yourdomain.com/callback"),
    Intouchpay.WithSid(0),
    Intouchpay.WithTimeout(30 * time.Second),
)
```

Available options:
- `WithTimeout(duration)` - Set HTTP client timeout
- `WithHTTPClient(*http.Client)` - Use a custom HTTP client
- `WithCallbackURL(url)` - Set callback URL
- `WithSid(sid)` - Set service ID
- `WithAuthenticator(auth)` - Use a custom authenticator (for testing)

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

## Testing

### Mock Authentication

For testing without real credentials, use a mock authenticator:

```go
// Create a mock authenticator
type MockAuthenticator struct {
    Creds Intouchpay.Credentials
}

func (m *MockAuthenticator) Authenticate() Intouchpay.Credentials {
    return m.Creds
}

// Use it in tests
mockAuth := &MockAuthenticator{
    Creds: Intouchpay.Credentials{
        Username:  "test_user",
        Timestamp: "20260320120000",
        Password:  "test_hash",
    },
}

client := Intouchpay.NewClientWithAuth(mockAuth)
```

### Mock HTTP Client

For testing without network calls, use a mock HTTP client:

```go
// Create a mock HTTP client
type MockHTTPClient struct {
    Response *map[string]interface{}
    Error    error
}

func (m *MockHTTPClient) Do(endpoint string, body interface{}) (*map[string]interface{}, error) {
    return m.Response, m.Error
}

// Use it in tests
mockHTTP := &MockHTTPClient{
    Response: &map[string]interface{}{
        "success": true,
        "balance": 10000.50,
    },
}

client := Intouchpay.NewClientWithHTTPClient(mockAuth, mockHTTP)
```

### Running Tests

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## Error Handling

### Typed Errors

The package provides typed errors for better error handling:

```go
response, err := client.RequestPayment(params)
if err != nil {
    // Check for API errors
    if Intouchpay.IsAPIError(err) {
        apiErr := err.(*Intouchpay.APIError)
        log.Printf("API Error: %d - %s", apiErr.StatusCode, apiErr.Status)
    }
    
    // Check for validation errors
    if Intouchpay.IsValidationError(err) {
        valErr := err.(*Intouchpay.ValidationError)
        log.Printf("Validation Error: %s - %s", valErr.Field, valErr.Message)
    }
    
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

### Error Types

- **APIError** - Returned when the API returns a non-200 status code
  - `StatusCode` - HTTP status code
  - `Status` - HTTP status text
  - `Response` - Parsed API response
  - `Message` - Error message

- **ValidationError** - Returned for client-side validation failures
  - `Field` - The field that failed validation
  - `Message` - Error message

## Architecture

This package follows deep module design principles for testability:

```
┌─────────────────────────────────────────────────────┐
│                     Client                          │
│  ┌──────────────┐  ┌──────────────┐                │
│  │ Authenticator│  │ APIRequester │                │
│  │  (interface) │  │  (interface) │                │
│  └──────────────┘  └──────────────┘                │
│         ↓                  ↓                        │
│  ┌──────────────┐  ┌──────────────┐                │
│  │  sha256Auth  │  │ defaultHTTPClient │            │
│  │ (production) │  │ (production) │                │
│  └──────────────┘  └──────────────┘                │
└─────────────────────────────────────────────────────┘
```

### Interfaces

- **Authenticator** - Generates authentication credentials
  - `Authenticate() Credentials`

- **APIRequester** - Makes HTTP requests to the API
  - `Do(endpoint string, body interface{}) (*map[string]interface{}, error)`

### Testing-Friendly Constructors

```go
// With mock authenticator
client := Intouchpay.NewClientWithAuth(mockAuth)

// With mock HTTP client
client := Intouchpay.NewClientWithHTTPClient(mockAuth, mockHTTP)

// With options
client := Intouchpay.NewClientWithOptions(
    "username", "account", "password",
    Intouchpay.WithTimeout(30*time.Second),
    Intouchpay.WithAuthenticator(mockAuth),
    Intouchpay.WithHTTPClientInterface(mockHTTP),
)
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

### Default Timeout

The default HTTP client timeout is **60 seconds**. You can customize this with:

```go
client := Intouchpay.NewClientWithOptions(
    "username", "account", "password",
    Intouchpay.WithTimeout(30 * time.Second),
)
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
        Amount:               1000,
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
        Amount:               5000,
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

## Migration Guide

### Upgrading from v0.1.2 to v2.0.0

**No code changes required!** This release is fully backward compatible.

#### Existing Code (Still Works)

```go
// This continues to work exactly as before
client := Intouchpay.NewClient(
    "your_username",
    "your_account_number",
    "your_partner_password",
    "https://yourdomain.com/callback",
    0,
)
```

#### New Features Available (Optional)

**1. Flexible Configuration with Options**

```go
// NEW: Use options for more control
client := Intouchpay.NewClientWithOptions(
    "your_username",
    "your_account_number",
    "your_partner_password",
    Intouchpay.WithCallbackURL("https://yourdomain.com/callback"),
    Intouchpay.WithSid(0),
    Intouchpay.WithTimeout(30 * time.Second),
)
```

**2. Better Error Handling**

```go
// NEW: Typed errors for clearer error handling
response, err := client.RequestPayment(params)
if err != nil {
    if Intouchpay.IsAPIError(err) {
        apiErr := err.(*Intouchpay.APIError)
        log.Printf("API returned status %d", apiErr.StatusCode)
    }
    if Intouchpay.IsValidationError(err) {
        valErr := err.(*Intouchpay.ValidationError)
        log.Printf("Invalid field: %s", valErr.Field)
    }
    return
}
```

**3. Testing with Mocks**

```go
// NEW: Mock authentication and HTTP for tests
mockAuth := &MockAuthenticator{Creds: Intouchpay.Credentials{
    Username: "test", Timestamp: "20260320120000", Password: "hash",
}}
mockHTTP := &MockHTTPClient{Response: &map[string]interface{}{"success": true}}

client := Intouchpay.NewClientWithHTTPClient(mockAuth, mockHTTP)
```

#### Timeout Change

The default HTTP timeout changed from undefined to **60 seconds**. This should not affect most users, but if you need a different timeout:

```go
client := Intouchpay.NewClientWithOptions(
    "username", "account", "password",
    Intouchpay.WithTimeout(30 * time.Second),
)
```

## Changelog

### v2.0.0

**Architecture Improvements:**
- Extracted authentication into `Authenticator` interface for testability
- Extracted HTTP layer into `APIRequester` interface for mocking
- Added typed errors: `APIError` and `ValidationError`
- Added `PhoneValidator` struct for phone number validation
- Added option pattern for flexible client configuration
- Fixed timeout inconsistency (now 60 seconds default)

**New Constructors:**
- `NewClientWithOptions()` - Flexible configuration with options
- `NewClientWithAuth()` - Custom authenticator for testing
- `NewClientWithHTTPClient()` - Custom HTTP client for testing

**New Options:**
- `WithTimeout()` - Set HTTP timeout
- `WithHTTPClient()` - Use custom HTTP client
- `WithCallbackURL()` - Set callback URL
- `WithSid()` - Set service ID
- `WithAuthenticator()` - Use custom authenticator
- `WithHTTPClientInterface()` - Use mock HTTP client

**Testing:**
- Added comprehensive test suite (41 tests)
- Full coverage of authentication, HTTP layer, phone validation, errors

### v0.1.2

Initial release with basic API functionality.
