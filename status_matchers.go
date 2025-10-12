package inpu

import (
	"io"
	"slices"
)

type StatusMatcher interface {
	Match(statusCode int) bool
	Priority() int
}

func newStatusChecker(matcher func(statusCode int) bool, priority int) *statusChecker {
	return &statusChecker{
		matcher:  matcher,
		priority: priority,
	}
}

type statusChecker struct {
	matcher  func(statusCode int) bool
	priority int
}

func (s *statusChecker) Match(statusCode int) bool {
	return s.matcher(statusCode)
}

// Priority is the order in which the matchers are used.
// The less priority is the higher precedence on the matching.
// Current priorities are:
// StatusIs -> 1
// StatusIsOneOf -> 2
// StatusIsInformational, StatusIsSuccess, StatusIsRedirection, StatusIsClientError, StatusIsServerError -> 3
// StatusAnyExcept -> 8
// StatusAnyExceptOneOf -> 9
// StatusAny -> 10
func (s *statusChecker) Priority() int {
	return s.priority
}

// StatusAny matches any status code. It can be used as fallback in case previous OnReply matches is not called
// It has the least priority 10, and it is checked after StatusAnyExceptOneOf
// Usage:
// OnReply(StatusAny,func(r *http.Response) error{})
var StatusAny = newStatusChecker(func(_ int) bool {
	return true
}, 10)

// StatusAnyExceptOneOf matches any status code except those provided.
// It has the priority 9, and it is checked after StatusAnyExcept.
// Usage:
// OnReply(StatusAnyExceptOneOf(http.StatusOK,http.StatusCreated),func(r *http.Response) error{})
func StatusAnyExceptOneOf(statusCodes ...int) StatusMatcher {
	return newStatusChecker(func(statusCode int) bool {
		return !slices.Contains(statusCodes, statusCode)
	}, 9)
}

// StatusAnyExcept matches any status code except the one provided.
// It has the priority 8, and it is checked after StatusIsInformational, StatusIsSuccess, StatusIsRedirection, StatusIsClientError, StatusIsServerError.
// Usage:
// OnReply(StatusAnyExcept(http.StatusOK),func(r *http.Response) error{})
func StatusAnyExcept(statusCode int) StatusMatcher {
	return newStatusChecker(func(actualStatus int) bool {
		return statusCode != actualStatus
	}, 8)
}

// StatusIsSuccess checks if the response status code is between [200,300).
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReply(StatusIsSuccess,func(r *http.Response) error{})
var StatusIsSuccess = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}, 3)

// StatusIsInformational checks if the response status code is less than 200.
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReply(StatusIsInformational,func(r *http.Response) error{})
var StatusIsInformational = newStatusChecker(func(statusCode int) bool {
	return statusCode < 200
}, 3)

// StatusIsRedirection checks if the response status code is between [300,400).
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReply(StatusIsRedirection,func(r *http.Response) error{})
var StatusIsRedirection = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}, 3)

// StatusIsClientError checks if the response status code is between [400,500).
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReply(StatusIsClientError,func(r *http.Response) error{})
var StatusIsClientError = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}, 3)

// StatusIsServerError checks if the response status is greater or equal than 500.
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReply(StatusIsServerError,func(r *http.Response) error{})
var StatusIsServerError = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 500
}, 3)

// StatusIsOneOf checks if the response status is one of the provided codes.
// It has the priority 2, and it is checked after StatusIs.
// Usage:
// OnReply(StatusIsOneOf(http.StatusOK,http.StatusCreated),func(r *http.Response) error{}) -> only matches when the status code is either 200,201
func StatusIsOneOf(statusCodes ...int) StatusMatcher {
	return newStatusChecker(func(statusCode int) bool {
		return slices.Contains(statusCodes, statusCode)
	}, 2)
}

// StatusIs checks if the response status is the provided status code.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIs(http.StatusOK),func(r *http.Response) error{}) -> only matches when the status code is 200
func StatusIs(expectedStatus int) StatusMatcher {
	return newStatusChecker(func(actualStatus int) bool {
		return expectedStatus == actualStatus
	}, 1)
}

func DrainBodyAndClose(body io.ReadCloser) error {
	defer body.Close()
	// Limit drain to prevent memory issues with huge responses
	_, err := io.Copy(io.Discard, io.LimitReader(body, 1<<20)) // 1MB limit
	if err != nil {
		return err
	}

	return nil
}
