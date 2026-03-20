package Intouchpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPClient defines the interface for making HTTP requests to the IntouchPay API
type HTTPClient interface {
	// Do sends a POST request to the given endpoint with the provided body
	Do(endpoint string, body interface{}) (*map[string]interface{}, error)
}

// defaultHTTPClient implements HTTPClient using net/http
type defaultHTTPClient struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClient creates a new HTTP client with the provided configuration
func NewHTTPClient(httpClient *http.Client, baseURL string) HTTPClient {
	return &defaultHTTPClient{
		client:  httpClient,
		baseURL: baseURL,
	}
}

// Do sends a POST request to the given endpoint with the provided body
func (c *defaultHTTPClient) Do(endpoint string, body interface{}) (*map[string]interface{}, error) {
	var response *map[string]interface{}
	requestUrl := c.baseURL + endpoint

	jsonData, err := json.Marshal(body)
	if err != nil {
		return response, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return response, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("IntouchPay API error: %d\n %s\n %v", resp.StatusCode, resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return response, newAPIError(resp.StatusCode, resp.Status, *response)
	}

	return response, nil
}
