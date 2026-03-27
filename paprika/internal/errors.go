package internal

import "fmt"

// HTTPError represents an unexpected HTTP response from the API.
type HTTPError struct {
	StatusCode int
	Body       string
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// FormatHTTPError returns an error that includes the HTTP StatusCode and body.
func FormatHTTPError(statusCode int, body string) error {
	return fmt.Errorf("HTTP %d: %s", statusCode, body)
}
