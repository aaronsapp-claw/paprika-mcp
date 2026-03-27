package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultBaseURL is the default API base URL embedded from the OpenAPI spec.
const DefaultBaseURL = "https://www.paprikaapp.com/api/v2/sync"


// Client holds the configuration for making authenticated HTTP requests.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient constructs a Client with the given baseURL and token.
// When baseURL is empty, DefaultBaseURL is used.
func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// execute injects authentication headers into req, sends it, reads the response
// body, and returns an error for non-2xx status codes. Both Do and DoMultipart
// delegate to this method so auth injection is defined exactly once.
func (c *Client) execute(req *http.Request) ([]byte, error) {
	// Inject authentication credentials based on security schemes.
	// Bearer takes precedence over Basic when both schemes are present.
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Return a descriptive error for non-2xx responses.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyStr := string(responseBody)
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "... (truncated)"
		}
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, bodyStr)
	}

	return responseBody, nil
}

// Do executes an HTTP request against the API.
//
// method is the HTTP verb (GET, POST, etc.).
// path is the URL path template (e.g., "/users/{userId}").
// pathParams maps placeholder names to their runtime values for path interpolation.
// queryParams maps query parameter names to values appended to the URL.
// body is an optional request body; pass nil for requests without a body.
//
// Path parameter substitution uses strings.NewReplacer to replace {param}
// placeholders with the corresponding values from pathParams.
func (c *Client) Do(method, path string, pathParams map[string]string, queryParams map[string]string, body interface{}) ([]byte, error) {
	// Interpolate {param} placeholders in the path template.
	pairs := make([]string, 0, len(pathParams)*2)
	for key, value := range pathParams {
		pairs = append(pairs, "{"+key+"}", value)
	}
	interpolatedPath := strings.NewReplacer(pairs...).Replace(path)

	requestURL := strings.TrimRight(c.BaseURL, "/") + interpolatedPath

	// Append query parameters.
	if len(queryParams) > 0 {
		separator := "?"
		for key, value := range queryParams {
			requestURL += separator + key + "=" + url.QueryEscape(value)
			separator = "&"
		}
	}

	// Encode the request body as JSON when provided.
	var bodyReader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
		bodyReader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.execute(req)
}

// DoMultipart executes a multipart/form-data HTTP request against the API.
//
// method is the HTTP verb (typically POST or PUT).
// path is the URL path template (e.g., "/documents/upload").
// pathParams maps placeholder names to their runtime values for path interpolation.
// queryParams maps query parameter names to values appended to the URL.
// body is an io.Reader containing the pre-built multipart body (built by the caller
// using mime/multipart.Writer).
// contentType is the full Content-Type header value including the boundary parameter,
// obtained from multipart.Writer.FormDataContentType().
func (c *Client) DoMultipart(method, path string, pathParams map[string]string, queryParams map[string]string, body io.Reader, contentType string) ([]byte, error) {
	// Interpolate {param} placeholders in the path template.
	pairs := make([]string, 0, len(pathParams)*2)
	for key, value := range pathParams {
		pairs = append(pairs, "{"+key+"}", value)
	}
	interpolatedPath := strings.NewReplacer(pairs...).Replace(path)

	requestURL := strings.TrimRight(c.BaseURL, "/") + interpolatedPath

	// Append query parameters.
	if len(queryParams) > 0 {
		separator := "?"
		for key, value := range queryParams {
			requestURL += separator + key + "=" + url.QueryEscape(value)
			separator = "&"
		}
	}

	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set Content-Type to multipart/form-data with the writer's boundary.
	req.Header.Set("Content-Type", contentType)

	return c.execute(req)
}
