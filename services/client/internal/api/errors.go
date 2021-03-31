package api

import (
	"fmt"
)

// HTTPError wraps a `net/http` error that occurs during a request.
type HTTPError struct {
	WrappedError error
}

// Error prints the HTTPError as a string.
func (e *HTTPError) Error() string {
	return e.WrappedError.Error()
}

// IntoHTTPError takes an error and wraps it with a *HTTPError.
func IntoHTTPError(err error) *HTTPError {
	return &HTTPError{
		err,
	}
}

// APIError is returned by an API client when an error occurs.
type APIError struct {
	ResponseBody string // The body of the response
	Message      string // Error message
	Reference    string // An error reference, eg. "deserialization_failure"
	WrappedError error  // The originating error, if applicable (otherwise `nil`)
}

// Error prints the APIError as a string.
func (e *APIError) Error() string {
	return fmt.Sprintf("api error %s\n\t%s\n\tgot body: %s", e.Reference, e.Message, e.ResponseBody)
}
