// Package Intouchpay provides a Go client for the IntouchPay Payments Gateway API.
// It supports mobile money payments, deposits, balance inquiries, and transaction status checks.
package Intouchpay

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Credentials represents authentication data for a single API request
type Credentials struct {
	Username  string
	Timestamp string
	Password  string
}

// Authenticator generates credentials for API requests
type Authenticator interface {
	// Authenticate generates credentials for a single API call
	Authenticate() Credentials
}

// sha256Auth implements Authenticator using SHA256 hashing
type sha256Auth struct {
	username        string
	accountNo       string
	partnerPassword string
}

// NewAuthenticator creates the default authenticator with standard SHA256 hashing
func NewAuthenticator(username, accountNo, partnerPassword string) Authenticator {
	return &sha256Auth{
		username:        username,
		accountNo:       accountNo,
		partnerPassword: partnerPassword,
	}
}

// Authenticate generates credentials for a single API call
func (a *sha256Auth) Authenticate() Credentials {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	passwordString := a.username + a.accountNo + a.partnerPassword + timestamp
	hash := sha256.Sum256([]byte(passwordString))
	password := hex.EncodeToString(hash[:])

	return Credentials{
		Username:  a.username,
		Timestamp: timestamp,
		Password:  password,
	}
}
