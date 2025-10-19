package inpu

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// ResponseHandler is the function you pass along with a status matcher in the OnReply call.
// It can be
type ResponseHandler func(response *http.Response) error

// UnmarshalJson marshals the body to the pointer provided in the targetAsPointer argument.
// It checks if the type is pointer as well.
// It does not close the body because body is closed after this function is called by the caller
// Usage:
// OnReply(StatusAny, UnmarshalJson(&items))
func UnmarshalJson(targetAsPointer any) ResponseHandler {
	return func(r *http.Response) error {
		if targetAsPointer == nil {
			return ErrMarshalToNil
		}

		if reflect.ValueOf(targetAsPointer).Kind() != reflect.Ptr {
			return ErrNotPointerParameter
		}

		return json.NewDecoder(r.Body).Decode(targetAsPointer)
	}
}

// ReturnError returns the provided error directly in case of status is matched
// Usage:
// OnReply(StatusAny, ReturnError(errors.New("something happened")))
func ReturnError(err error) ResponseHandler {
	return func(_ *http.Response) error {
		return err
	}
}

// ReturnDefaultError returns an error that contains the request method, requests URL and the status code
// Usage:
// OnReply(StatusAny, ReturnDefaultError)
func ReturnDefaultError(r *http.Response) error {
	return &DefaultError{
		Method:     r.Request.Method,
		URL:        r.Request.URL.Redacted(),
		StatusCode: r.StatusCode,
	}
}

// DoNothing returns nil error
// Usage:
// OnReply(StatusAny, DoNothing)
func DoNothing(_ *http.Response) error {
	return nil
}
