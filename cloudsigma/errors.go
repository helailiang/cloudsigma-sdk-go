package cloudsigma

import (
	"errors"
	"fmt"
)

const (
	ErrorType_validation  = "validation"
	ErrorType_notexist   = "notexist"
	ErrorType_backend     = "backend"
	ErrorType_permission  = "permission"
	ErrorType_database    = "database"
	ErrorType_concurrency = "concurrency"
	ErrorType_billing     = "billing"
	ErrorType_payment     = "payment"
)

// Errors used by the CloudSigma SDK.
var (
	// ErrEmptyPayloadNotAllowed is returned when a request body is empty
	// and does not contain a mandatory JSON payload.
	ErrEmptyPayloadNotAllowed = errors.New("cloudsigma: empty payload not allowed")

	// ErrEmptyArgument is returned when a mandatory function argument is empty.
	ErrEmptyArgument = errors.New("cloudsigma: argument cannot be empty")

	ErrServerNoExist = errors.New("server dont exist")
)

// An ErrorResponse reports one or more errors caused by an API request.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/errors.html
type ErrorResponse struct {
	Response *Response // HTTP response that caused this error.
	Errors   []Error
}

// Error represents a single error caused by an API request.
type Error struct {
	Message string `json:"error_message"`
	Point   string `json:"error_point"`
	Type    string `json:"error_type"`
}

// Error represents a string error message (may contain request id).
func (r *ErrorResponse) Error() string {
	if r.Response.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %+v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Response.RequestID, r.Errors)
	}
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Errors)
}
