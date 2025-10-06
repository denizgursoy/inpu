package inpu

import (
	"encoding/json"
	"net/http"
	"reflect"
	"slices"
)

type statusMatcher func(statusCode int) bool
type processor func(r *http.Response) error

func StatusAny(statusCode int) bool {
	return true
}

func StatusAnyExcept(statusCode int) func(actualStatus int) bool {
	return func(actualStatus int) bool {
		return statusCode != actualStatus
	}
}

func StatusAnyExceptOneOf(statusCodes ...int) func(actualStatus int) bool {
	return func(statusCode int) bool {
		return !slices.Contains(statusCodes, statusCode)
	}
}

func StatusIsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
func StatusIsInformational(statusCode int) bool {
	return statusCode < 200
}

func StatusIsRedirection(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

func StatusIsClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func StatusIsServerError(statusCode int) bool {
	return statusCode >= 500
}

func StatusIsOneOf(statusCodes ...int) func(statusCode int) bool {
	return func(statusCode int) bool {
		return slices.Contains(statusCodes, statusCode)
	}
}

func StatusIs(expectedStatus int) func(statusCode int) bool {
	return func(actualStatus int) bool {
		return expectedStatus == actualStatus
	}
}

func UnmarshalJson(t any) func(r *http.Response) error {
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

func ReturnError(err error) func(r *http.Response) error {
	return func(_ *http.Response) error {
		return err
	}
}
