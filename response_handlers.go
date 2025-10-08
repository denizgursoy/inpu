package inpu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func UnmarshalJson(t any) ResponseHandler {
	return func(r *http.Response) error {
		if t == nil {
			return ErrMarshalToNil
		}

		if reflect.ValueOf(t).Kind() != reflect.Ptr {
			return ErrNotPointerParameter
		}

		return json.NewDecoder(r.Body).Decode(t)
	}
}

func ReturnError(err error) ResponseHandler {
	return func(_ *http.Response) error {
		return err
	}
}

func ReturnDefaultError(r *http.Response) error {
	return fmt.Errorf("called [%s] %s and got %d", r.Request.Method, r.Request.URL.String(), r.StatusCode)
}
