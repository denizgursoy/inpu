package inpu

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

// ResponseHandler is the function you pass along with a status matcher in the OnReplyIf call.
// It can be
type ResponseHandler func(response *http.Response) error

// ThenUnmarshalJsonTo marshals the body to the pointer provided in the targetAsPointer argument.
// It checks if the type is pointer as well.
// It does not close the body because body is closed after this function is called by the caller
// Usage:
// OnReplyIf(StatusAny, ThenUnmarshalJsonTo(&items))
func ThenUnmarshalJsonTo(targetAsPointer any) ResponseHandler {
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

// ThenUnmarshalJsonAndReturnError marshals the body to the pointer with ThenUnmarshalJsonTo and returns provided error.
// Usage:
// OnReplyIf(StatusAny, ThenUnmarshalJsonAndReturnError(&items, errors.New("request failed")))
func ThenUnmarshalJsonAndReturnError(targetAsPointer any, err error) ResponseHandler {
	return func(response *http.Response) error {
		if marshalError := ThenUnmarshalJsonTo(targetAsPointer)(response); err != nil {
			return errors.Join(marshalError, err)
		}

		return err
	}
}

// ThenReturnError returns the provided error directly in case of status is matched
// Usage:
// OnReplyIf(StatusAny, ThenReturnError(errors.New("something happened")))
func ThenReturnError(err error) ResponseHandler {
	return func(_ *http.Response) error {
		return err
	}
}

// ThenReturnDefaultError returns an error that contains the request method, requests URL and the status code
// Usage:
// OnReplyIf(StatusAny, ThenReturnDefaultError)
func ThenReturnDefaultError(r *http.Response) error {
	return &DefaultError{
		res: r,
	}
}

// ThenDoNothing returns nil error
// Usage:
// OnReplyIf(StatusAny, ThenDoNothing)
func ThenDoNothing(_ *http.Response) error {
	return nil
}
