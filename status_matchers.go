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

// 2xx Success
// StatusIsOk checks if the response status is 200 OK.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsOk, func(r *http.Response) error{}) -> only matches when the status code is 200
var StatusIsOk = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusOK
}, 1)

// StatusIsCreated checks if the response status is 201 Created.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsCreated, func(r *http.Response) error{}) -> only matches when the status code is 201
var StatusIsCreated = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusCreated
}, 1)

// StatusIsNoContent checks if the response status is 204 No Content.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsNoContent, func(r *http.Response) error{}) -> only matches when the status code is 204
var StatusIsNoContent = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNoContent
}, 1)

// 3xx Redirection

// StatusIsMovedPermanently checks if the response status is 301 Moved Permanently.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsMovedPermanently, func(r *http.Response) error{}) -> only matches when the status code is 301
var StatusIsMovedPermanently = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently
}, 1)

// StatusIsFound checks if the response status is 302 Found.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsFound, func(r *http.Response) error{}) -> only matches when the status code is 302
var StatusIsFound = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusFound
}, 1)

// StatusIsNotModified checks if the response status is 304 Not Modified.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsNotModified, func(r *http.Response) error{}) -> only matches when the status code is 304
var StatusIsNotModified = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotModified
}, 1)

// 4xx Client Errors

// StatusIsBadRequest checks if the response status is 400 Bad Request.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsBadRequest, func(r *http.Response) error{}) -> only matches when the status code is 400
var StatusIsBadRequest = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusBadRequest
}, 1)

// StatusIsUnauthorized checks if the response status is 401 Unauthorized.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsUnauthorized, func(r *http.Response) error{}) -> only matches when the status code is 401
var StatusIsUnauthorized = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUnauthorized
}, 1)

// StatusIsForbidden checks if the response status is 403 Forbidden.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsForbidden, func(r *http.Response) error{}) -> only matches when the status code is 403
var StatusIsForbidden = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusForbidden
}, 1)

// StatusIsNotFound checks if the response status is 404 Not Found.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsNotFound, func(r *http.Response) error{}) -> only matches when the status code is 404
var StatusIsNotFound = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotFound
}, 1)

// StatusIsMethodNotAllowed checks if the response status is 405 Method Not Allowed.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsMethodNotAllowed, func(r *http.Response) error{}) -> only matches when the status code is 405
var StatusIsMethodNotAllowed = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMethodNotAllowed
}, 1)

// StatusIsTooManyRequests checks if the response status is 429 Too Many Requests.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsTooManyRequests, func(r *http.Response) error{}) -> only matches when the status code is 429
var StatusIsTooManyRequests = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests
}, 1)

// 5xx Server Errors

// StatusIsInternalServerError checks if the response status is 500 Internal Server Error.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsInternalServerError, func(r *http.Response) error{}) -> only matches when the status code is 500
var StatusIsInternalServerError = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusInternalServerError
}, 1)

// StatusIsBadGateway checks if the response status is 502 Bad Gateway.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsBadGateway, func(r *http.Response) error{}) -> only matches when the status code is 502
var StatusIsBadGateway = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusBadGateway
}, 1)

// StatusIsServiceUnavailable checks if the response status is 503 Service Unavailable.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsServiceUnavailable, func(r *http.Response) error{}) -> only matches when the status code is 503
var StatusIsServiceUnavailable = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusServiceUnavailable
}, 1)

// StatusIsGatewayTimeout checks if the response status is 504 Gateway Timeout.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReply(StatusIsGatewayTimeout, func(r *http.Response) error{}) -> only matches when the status code is 504
var StatusIsGatewayTimeout = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusGatewayTimeout
}, 1)

func DrainBodyAndClose(body io.ReadCloser) error {
	defer body.Close()
	// Limit drain to prevent memory issues with huge responses
	_, err := io.Copy(io.Discard, io.LimitReader(body, 1<<20)) // 1MB limit
	if err != nil {
		return err
	}

	return nil
}
