package devtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RequestMethod represents an HTTP request method
type RequestMethod string

const (
	GET     RequestMethod = "GET"
	POST    RequestMethod = "POST"
	PUT     RequestMethod = "PUT"
	DELETE  RequestMethod = "DELETE"
	PATCH   RequestMethod = "PATCH"
	OPTIONS RequestMethod = "OPTIONS"
	HEAD    RequestMethod = "HEAD"
)

// APIRequest represents an API request
type APIRequest struct {
	URL     string            `json:"url"`
	Method  RequestMethod     `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Timeout int               `json:"timeout"` // in seconds
}

// APIResponse represents an API response
type APIResponse struct {
	StatusCode int               `json:"statusCode"`
	Status     string            `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Duration   int64             `json:"duration"` // in milliseconds
	Error      string            `json:"error,omitempty"`
}

// APITester provides functionality to test API endpoints
type APITester struct {
	client *http.Client
}

// NewAPITester creates a new APITester
func NewAPITester() *APITester {
	return &APITester{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendRequest sends an API request and returns the response
func (at *APITester) SendRequest(req APIRequest) APIResponse {
	// Set default timeout if not specified
	timeout := 30
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	// Create HTTP client with the specified timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(string(req.Method), req.URL, bytes.NewBufferString(req.Body))
	if err != nil {
		return APIResponse{
			StatusCode: 0,
			Status:     "Error",
			Headers:    make(map[string]string),
			Body:       "",
			Duration:   0,
			Error:      fmt.Sprintf("Error creating request: %v", err),
		}
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set default Content-Type if not specified and body is not empty
	if req.Body != "" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Send request and measure duration
	startTime := time.Now()
	resp, err := client.Do(httpReq)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		return APIResponse{
			StatusCode: 0,
			Status:     "Error",
			Headers:    make(map[string]string),
			Body:       "",
			Duration:   duration,
			Error:      fmt.Sprintf("Error sending request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    convertHeaders(resp.Header),
			Body:       "",
			Duration:   duration,
			Error:      fmt.Sprintf("Error reading response body: %v", err),
		}
	}

	// Format JSON response if possible
	var formattedBody string
	if isJSON(string(body)) {
		var jsonObj interface{}
		if err := json.Unmarshal(body, &jsonObj); err == nil {
			formattedJSON, err := json.MarshalIndent(jsonObj, "", "  ")
			if err == nil {
				formattedBody = string(formattedJSON)
			} else {
				formattedBody = string(body)
			}
		} else {
			formattedBody = string(body)
		}
	} else {
		formattedBody = string(body)
	}

	// Return response
	return APIResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    convertHeaders(resp.Header),
		Body:       formattedBody,
		Duration:   duration,
		Error:      "",
	}
}

// convertHeaders converts http.Header to map[string]string
func convertHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// isJSON checks if a string is valid JSON
func isJSON(str string) bool {
	var js interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

// GetSavedRequests returns a list of saved API requests
// In a real implementation, these would be loaded from a file or database
func (at *APITester) GetSavedRequests() []APIRequest {
	return []APIRequest{
		{
			URL:     "https://jsonplaceholder.typicode.com/posts/1",
			Method:  GET,
			Headers: map[string]string{"Accept": "application/json"},
			Body:    "",
			Timeout: 30,
		},
		{
			URL:     "https://jsonplaceholder.typicode.com/posts",
			Method:  POST,
			Headers: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
			Body:    `{"title": "foo", "body": "bar", "userId": 1}`,
			Timeout: 30,
		},
		{
			URL:     "https://jsonplaceholder.typicode.com/posts/1",
			Method:  PUT,
			Headers: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
			Body:    `{"id": 1, "title": "foo", "body": "bar", "userId": 1}`,
			Timeout: 30,
		},
	}
}
