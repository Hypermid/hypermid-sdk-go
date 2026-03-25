package hypermid

import "fmt"

// ApiMeta contains metadata returned with every API response.
type ApiMeta struct {
	RequestID string         `json:"requestId"`
	Timestamp int64          `json:"timestamp"`
	RateLimit *RateLimitInfo `json:"rateLimit,omitempty"`
}

// RateLimitInfo contains rate limit details from the API response.
type RateLimitInfo struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}

// ApiError represents the error object in an API response envelope.
type ApiError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HyperMidError is returned when the API responds with an error.
type HyperMidError struct {
	// Code is the API error code (e.g. "NO_ROUTE_FOUND", "RATE_LIMITED").
	Code string
	// Msg is the human-readable error message.
	Msg string
	// Status is the HTTP status code.
	Status int
	// Meta contains request metadata (requestId, timestamp, rateLimit).
	Meta ApiMeta
	// Details contains additional error details (lifiCode, toolErrors, etc.).
	Details map[string]interface{}
}

func (e *HyperMidError) Error() string {
	return fmt.Sprintf("hypermid: %s: %s (HTTP %d, requestId=%s)", e.Code, e.Msg, e.Status, e.Meta.RequestID)
}

// HyperMidTimeoutError is returned when a request exceeds the configured timeout.
type HyperMidTimeoutError struct {
	TimeoutMs int
}

func (e *HyperMidTimeoutError) Error() string {
	return fmt.Sprintf("hypermid: request timed out after %dms", e.TimeoutMs)
}

// HyperMidNetworkError is returned when a network-level error occurs.
type HyperMidNetworkError struct {
	Msg   string
	Cause error
}

func (e *HyperMidNetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("hypermid: network error: %s: %v", e.Msg, e.Cause)
	}
	return fmt.Sprintf("hypermid: network error: %s", e.Msg)
}

func (e *HyperMidNetworkError) Unwrap() error {
	return e.Cause
}

// HyperMidPollTimeoutError is returned when status polling exceeds the maximum wait time or attempts.
type HyperMidPollTimeoutError struct {
	Msg string
}

func (e *HyperMidPollTimeoutError) Error() string {
	return fmt.Sprintf("hypermid: %s", e.Msg)
}
