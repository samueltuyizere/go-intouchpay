package Intouchpay

import (
	"net/http"
	"time"
)

// DefaultTimeout is the default HTTP client timeout
const DefaultTimeout = 60 * time.Second

// Option configures a Client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.HTTPClient = &http.Client{Timeout: timeout}
	}
}

// WithCallbackURL sets the callback URL
func WithCallbackURL(url string) Option {
	return func(c *Client) {
		c.CallbackURL = url
	}
}

// WithSid sets the service ID
func WithSid(sid int) Option {
	return func(c *Client) {
		c.Sid = sid
	}
}

// WithAuthenticator sets a custom authenticator
func WithAuthenticator(auth Authenticator) Option {
	return func(c *Client) {
		c.auth = auth
	}
}
