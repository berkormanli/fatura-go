package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/berkormanli/fatura-go/errors"
)

type Client struct {
	HttpClient *http.Client
	BaseURL    string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Request performs a GET or POST request
// In PHP: constructor does the request.
// In Go, we'll have a method.
func (c *Client) Request(endpoint string, params map[string]string, post bool) (map[string]interface{}, error) {
	// Params are form-urlencoded
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	
	var req *http.Request
	var err error
	
	if post {
		req, err = http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		// GET
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		u.RawQuery = data.Encode()
		req, err = http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}
	}
	
	// Execute
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, errors.NewBadResponseError(err.Error(), params, nil)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewBadResponseError("Failed to read response body", params, nil)
	}
	
	var responseMap map[string]interface{}
	// Attempt JSON decode
	// API might return JSON?
	// PHP: json_decode($contents, true)
	
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
			// If not JSON, it might be an error page or simple string?
			// The PHP code assumes JSON. If not JSON, $response is null/false.
			// And it checks if !$this->response => throws ApiException.
			// So strict JSON is expected.
		}
	}
	
	if responseMap == nil {
		return nil, errors.NewBadResponseError("Invalid JSON response", params, string(bodyBytes))
	}
	
	// Error checks logic from PHP:
	// if (!$this->response || isset($this->response['error']) || !empty($this->response['data']['hata']))
	
	if _, hasError := responseMap["error"]; hasError {
		return responseMap, errors.NewApiError("İstek başarısız oldu.", params, responseMap)
	}
	
	if dataMap, ok := responseMap["data"].(map[string]interface{}); ok {
		if val, exists := dataMap["hata"]; exists && val != nil && val != "" {
			return responseMap, errors.NewApiError("İstek başarısız oldu.", params, responseMap)
		}
	}
	
	return responseMap, nil
}

// RequestJSON performs a POST request with JSON body
func (c *Client) RequestJSON(endpoint string, payload interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(jsonBytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, errors.NewBadResponseError(err.Error(), nil, nil)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewBadResponseError("Failed to read response body", nil, nil)
	}
	
	var responseMap map[string]interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
			return nil, errors.NewBadResponseError("Invalid JSON response", nil, string(bodyBytes))
		}
	}
	
	if responseMap == nil {
		return nil, errors.NewBadResponseError("Empty response", nil, string(bodyBytes))
	}
	
	return responseMap, nil
}
