package inpu

import (
	"io"
	"net/http"
	"slices"
)

type StatusMatcher interface {
	Match(statusCode int) bool
	Priority() int
}
type ResponseHandler func(r *http.Response) error

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

func (s *statusChecker) Priority() int {
	return s.priority
}

var StatusAny = newStatusChecker(func(_ int) bool {
	return true
}, 10)

func StatusAnyExcept(statusCode int) StatusMatcher {
	return newStatusChecker(func(actualStatus int) bool {
		return statusCode != actualStatus
	}, 8)
}

func StatusAnyExceptOneOf(statusCodes ...int) StatusMatcher {
	return newStatusChecker(func(statusCode int) bool {
		return !slices.Contains(statusCodes, statusCode)
	}, 9)
}

var StatusIsSuccess = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}, 3)

var StatusIsInformational = newStatusChecker(func(statusCode int) bool {
	return statusCode < 200
}, 3)

var StatusIsRedirection = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}, 3)

var StatusIsClientError = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}, 3)

var StatusIsServerError = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 500
}, 3)

func StatusIsOneOf(statusCodes ...int) StatusMatcher {
	return newStatusChecker(func(statusCode int) bool {
		return slices.Contains(statusCodes, statusCode)
	}, 2)
}

func StatusIs(expectedStatus int) StatusMatcher {
	return newStatusChecker(func(actualStatus int) bool {
		return expectedStatus == actualStatus
	}, 1)
}

func DrainBodyAndClose(body io.ReadCloser) error {
	defer body.Close()
	// Limit drain to prevent memory issues with huge responses
	_, err := io.CopyN(io.Discard, body, 1<<20) // 1MB limit
	if err != nil {
		return err
	}

	return nil
}
