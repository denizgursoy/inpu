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
// StatusAnyExceptOneOf ->
// StatusAny -> 10
func (s *statusChecker) Priority() int {
	return s.priority
}

// StatusAny matches any status code. It can be used as fallback in case previous OnReplyIf matches is not called
// It has the least priority 10, and it is checked after StatusAnyExceptOneOf
// Usage:
// OnReplyIf(StatusAny,func(r *http.Response) error{})
var StatusAny = newStatusChecker(func(_ int) bool {
	return true
}, 10)

// StatusAnyExceptOneOf matches any status code except those provided.
// It has the priority 9, and it is checked after StatusAnyExcept.
// Usage:
// OnReplyIf(StatusAnyExceptOneOf(http.StatusOK,http.StatusCreated),func(r *http.Response) error{})
func StatusAnyExceptOneOf(statusCodes ...int) StatusMatcher {
	return newStatusChecker(func(statusCode int) bool {
		return !slices.Contains(statusCodes, statusCode)
	}, 9)
}

// StatusAnyExcept matches any status code except the one provided.
// It has the priority 8, and it is checked after StatusIsInformational, StatusIsSuccess, StatusIsRedirection, StatusIsClientError, StatusIsServerError.
// Usage:
// OnReplyIf(StatusAnyExcept(http.StatusOK),func(r *http.Response) error{})
func StatusAnyExcept(statusCode int) StatusMatcher {
	return newStatusChecker(func(actualStatus int) bool {
		return statusCode != actualStatus
	}, 8)
}

// StatusIsSuccess checks if the response status code is between [200,300).
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReplyIf(StatusIsSuccess,func(r *http.Response) error{})
var StatusIsSuccess = newStatusChecker(func(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}, 3)

// StatusIsInformational checks if the response status code is less than 200.
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReplyIf(StatusIsInformational,func(r *http.Response) error{})
var StatusIsInformational = newStatusChecker(func(statusCode int) bool {
	return statusCode < http.StatusOK
}, 3)

// StatusIsRedirection checks if the response status code is between [300,400).
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReplyIf(StatusIsRedirection,func(r *http.Response) error{})
var StatusIsRedirection = newStatusChecker(func(statusCode int) bool {
	return statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest
}, 3)

// StatusIsClientError checks if the response status code is between [400,500).
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReplyIf(StatusIsClientError,func(r *http.Response) error{})
var StatusIsClientError = newStatusChecker(func(statusCode int) bool {
	return statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError
}, 3)

// StatusIsServerError checks if the response status is greater or equal than 500.
// It has the priority 3, and it is checked after StatusIsOneOf.
// Usage:
// OnReplyIf(StatusIsServerError,func(r *http.Response) error{})
var StatusIsServerError = newStatusChecker(func(statusCode int) bool {
	return statusCode >= 500
}, 3)

// StatusIsOneOf checks if the response status is one of the provided codes.
// It has the priority 2, and it is checked after StatusIs.
// Usage:
// OnReplyIf(StatusIsOneOf(http.StatusOK,http.StatusCreated),func(r *http.Response) error{}) -> only matches when the status code is either 200,201
func StatusIsOneOf(statusCodes ...int) StatusMatcher {
	return newStatusChecker(func(statusCode int) bool {
		return slices.Contains(statusCodes, statusCode)
	}, 2)
}

// StatusIs checks if the response status is the provided status code.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIs(http.StatusOK),func(r *http.Response) error{}) -> only matches when the status code is 200
func StatusIs(expectedStatus int) StatusMatcher {
	return newStatusChecker(func(actualStatus int) bool {
		return expectedStatus == actualStatus
	}, 1)
}

// 1xx Informational

// StatusIsContinue checks if the response status is 100 Continue.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsContinue, func(r *http.Response) error{}) -> only matches when the status code is 100
var StatusIsContinue = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusContinue
}, 1)

// StatusIsSwitchingProtocols checks if the response status is 101 Switching Protocols.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsSwitchingProtocols, func(r *http.Response) error{}) -> only matches when the status code is 101
var StatusIsSwitchingProtocols = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusSwitchingProtocols
}, 1)

// StatusIsProcessing checks if the response status is 102 Processing.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsProcessing, func(r *http.Response) error{}) -> only matches when the status code is 102
var StatusIsProcessing = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusProcessing
}, 1)

// StatusIsEarlyHints checks if the response status is 103 Early Hints.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsEarlyHints, func(r *http.Response) error{}) -> only matches when the status code is 103
var StatusIsEarlyHints = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusEarlyHints
}, 1)

// 2xx Success

// StatusIsOk checks if the response status is 200 OK.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsOk, func(r *http.Response) error{}) -> only matches when the status code is 200
var StatusIsOk = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusOK
}, 1)

// StatusIsCreated checks if the response status is 201 Created.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsCreated, func(r *http.Response) error{}) -> only matches when the status code is 201
var StatusIsCreated = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusCreated
}, 1)

// StatusIsAccepted checks if the response status is 202 Accepted.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsAccepted, func(r *http.Response) error{}) -> only matches when the status code is 202
var StatusIsAccepted = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusAccepted
}, 1)

// StatusIsNonAuthoritativeInfo checks if the response status is 203 Non-Authoritative Information.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNonAuthoritativeInfo, func(r *http.Response) error{}) -> only matches when the status code is 203
var StatusIsNonAuthoritativeInfo = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNonAuthoritativeInfo
}, 1)

// StatusIsNoContent checks if the response status is 204 No Content.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNoContent, func(r *http.Response) error{}) -> only matches when the status code is 204
var StatusIsNoContent = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNoContent
}, 1)

// StatusIsResetContent checks if the response status is 205 Reset Content.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsResetContent, func(r *http.Response) error{}) -> only matches when the status code is 205
var StatusIsResetContent = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusResetContent
}, 1)

// StatusIsPartialContent checks if the response status is 206 Partial Content.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsPartialContent, func(r *http.Response) error{}) -> only matches when the status code is 206
var StatusIsPartialContent = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusPartialContent
}, 1)

// StatusIsMultiStatus checks if the response status is 207 Multi-Status.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsMultiStatus, func(r *http.Response) error{}) -> only matches when the status code is 207
var StatusIsMultiStatus = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMultiStatus
}, 1)

// StatusIsAlreadyReported checks if the response status is 208 Already Reported.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsAlreadyReported, func(r *http.Response) error{}) -> only matches when the status code is 208
var StatusIsAlreadyReported = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusAlreadyReported
}, 1)

// StatusIsIMUsed checks if the response status is 226 IM Used.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsIMUsed, func(r *http.Response) error{}) -> only matches when the status code is 226
var StatusIsIMUsed = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusIMUsed
}, 1)

// 3xx Redirection

// StatusIsMultipleChoices checks if the response status is 300 Multiple Choices.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsMultipleChoices, func(r *http.Response) error{}) -> only matches when the status code is 300
var StatusIsMultipleChoices = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMultipleChoices
}, 1)

// StatusIsMovedPermanently checks if the response status is 301 Moved Permanently.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsMovedPermanently, func(r *http.Response) error{}) -> only matches when the status code is 301
var StatusIsMovedPermanently = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently
}, 1)

// StatusIsFound checks if the response status is 302 Found.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsFound, func(r *http.Response) error{}) -> only matches when the status code is 302
var StatusIsFound = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusFound
}, 1)

// StatusIsSeeOther checks if the response status is 303 See Other.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsSeeOther, func(r *http.Response) error{}) -> only matches when the status code is 303
var StatusIsSeeOther = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusSeeOther
}, 1)

// StatusIsNotModified checks if the response status is 304 Not Modified.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNotModified, func(r *http.Response) error{}) -> only matches when the status code is 304
var StatusIsNotModified = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotModified
}, 1)

// StatusIsUseProxy checks if the response status is 305 Use Proxy.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsUseProxy, func(r *http.Response) error{}) -> only matches when the status code is 305
var StatusIsUseProxy = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUseProxy
}, 1)

// StatusIsTemporaryRedirect checks if the response status is 307 Temporary Redirect.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsTemporaryRedirect, func(r *http.Response) error{}) -> only matches when the status code is 307
var StatusIsTemporaryRedirect = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusTemporaryRedirect
}, 1)

// StatusIsPermanentRedirect checks if the response status is 308 Permanent Redirect.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsPermanentRedirect, func(r *http.Response) error{}) -> only matches when the status code is 308
var StatusIsPermanentRedirect = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusPermanentRedirect
}, 1)

// 4xx Client Errors

// StatusIsBadRequest checks if the response status is 400 Bad Request.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsBadRequest, func(r *http.Response) error{}) -> only matches when the status code is 400
var StatusIsBadRequest = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusBadRequest
}, 1)

// StatusIsUnauthorized checks if the response status is 401 Unauthorized.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsUnauthorized, func(r *http.Response) error{}) -> only matches when the status code is 401
var StatusIsUnauthorized = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUnauthorized
}, 1)

// StatusIsPaymentRequired checks if the response status is 402 Payment Required.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsPaymentRequired, func(r *http.Response) error{}) -> only matches when the status code is 402
var StatusIsPaymentRequired = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusPaymentRequired
}, 1)

// StatusIsForbidden checks if the response status is 403 Forbidden.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsForbidden, func(r *http.Response) error{}) -> only matches when the status code is 403
var StatusIsForbidden = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusForbidden
}, 1)

// StatusIsNotFound checks if the response status is 404 Not Found.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNotFound, func(r *http.Response) error{}) -> only matches when the status code is 404
var StatusIsNotFound = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotFound
}, 1)

// StatusIsMethodNotAllowed checks if the response status is 405 Method Not Allowed.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsMethodNotAllowed, func(r *http.Response) error{}) -> only matches when the status code is 405
var StatusIsMethodNotAllowed = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMethodNotAllowed
}, 1)

// StatusIsNotAcceptable checks if the response status is 406 Not Acceptable.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNotAcceptable, func(r *http.Response) error{}) -> only matches when the status code is 406
var StatusIsNotAcceptable = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotAcceptable
}, 1)

// StatusIsProxyAuthRequired checks if the response status is 407 Proxy Authentication Required.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsProxyAuthRequired, func(r *http.Response) error{}) -> only matches when the status code is 407
var StatusIsProxyAuthRequired = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusProxyAuthRequired
}, 1)

// StatusIsRequestTimeout checks if the response status is 408 Request Timeout.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsRequestTimeout, func(r *http.Response) error{}) -> only matches when the status code is 408
var StatusIsRequestTimeout = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusRequestTimeout
}, 1)

// StatusIsConflict checks if the response status is 409 Conflict.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsConflict, func(r *http.Response) error{}) -> only matches when the status code is 409
var StatusIsConflict = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusConflict
}, 1)

// StatusIsGone checks if the response status is 410 Gone.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsGone, func(r *http.Response) error{}) -> only matches when the status code is 410
var StatusIsGone = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusGone
}, 1)

// StatusIsLengthRequired checks if the response status is 411 Length Required.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsLengthRequired, func(r *http.Response) error{}) -> only matches when the status code is 411
var StatusIsLengthRequired = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusLengthRequired
}, 1)

// StatusIsPreconditionFailed checks if the response status is 412 Precondition Failed.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsPreconditionFailed, func(r *http.Response) error{}) -> only matches when the status code is 412
var StatusIsPreconditionFailed = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusPreconditionFailed
}, 1)

// StatusIsRequestEntityTooLarge checks if the response status is 413 Request Entity Too Large.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsRequestEntityTooLarge, func(r *http.Response) error{}) -> only matches when the status code is 413
var StatusIsRequestEntityTooLarge = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusRequestEntityTooLarge
}, 1)

// StatusIsRequestURITooLong checks if the response status is 414 Request URI Too Long.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsRequestURITooLong, func(r *http.Response) error{}) -> only matches when the status code is 414
var StatusIsRequestURITooLong = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusRequestURITooLong
}, 1)

// StatusIsUnsupportedMediaType checks if the response status is 415 Unsupported Media Type.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsUnsupportedMediaType, func(r *http.Response) error{}) -> only matches when the status code is 415
var StatusIsUnsupportedMediaType = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUnsupportedMediaType
}, 1)

// StatusIsRequestedRangeNotSatisfiable checks if the response status is 416 Requested Range Not Satisfiable.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsRequestedRangeNotSatisfiable, func(r *http.Response) error{}) -> only matches when the status code is 416
var StatusIsRequestedRangeNotSatisfiable = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusRequestedRangeNotSatisfiable
}, 1)

// StatusIsExpectationFailed checks if the response status is 417 Expectation Failed.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsExpectationFailed, func(r *http.Response) error{}) -> only matches when the status code is 417
var StatusIsExpectationFailed = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusExpectationFailed
}, 1)

// StatusIsTeapot checks if the response status is 418 I'm a teapot.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsTeapot, func(r *http.Response) error{}) -> only matches when the status code is 418
var StatusIsTeapot = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusTeapot
}, 1)

// StatusIsMisdirectedRequest checks if the response status is 421 Misdirected Request.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsMisdirectedRequest, func(r *http.Response) error{}) -> only matches when the status code is 421
var StatusIsMisdirectedRequest = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusMisdirectedRequest
}, 1)

// StatusIsUnprocessableEntity checks if the response status is 422 Unprocessable Entity.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsUnprocessableEntity, func(r *http.Response) error{}) -> only matches when the status code is 422
var StatusIsUnprocessableEntity = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUnprocessableEntity
}, 1)

// StatusIsLocked checks if the response status is 423 Locked.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsLocked, func(r *http.Response) error{}) -> only matches when the status code is 423
var StatusIsLocked = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusLocked
}, 1)

// StatusIsFailedDependency checks if the response status is 424 Failed Dependency.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsFailedDependency, func(r *http.Response) error{}) -> only matches when the status code is 424
var StatusIsFailedDependency = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusFailedDependency
}, 1)

// StatusIsTooEarly checks if the response status is 425 Too Early.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsTooEarly, func(r *http.Response) error{}) -> only matches when the status code is 425
var StatusIsTooEarly = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusTooEarly
}, 1)

// StatusIsUpgradeRequired checks if the response status is 426 Upgrade Required.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsUpgradeRequired, func(r *http.Response) error{}) -> only matches when the status code is 426
var StatusIsUpgradeRequired = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUpgradeRequired
}, 1)

// StatusIsPreconditionRequired checks if the response status is 428 Precondition Required.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsPreconditionRequired, func(r *http.Response) error{}) -> only matches when the status code is 428
var StatusIsPreconditionRequired = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusPreconditionRequired
}, 1)

// StatusIsTooManyRequests checks if the response status is 429 Too Many Requests.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsTooManyRequests, func(r *http.Response) error{}) -> only matches when the status code is 429
var StatusIsTooManyRequests = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests
}, 1)

// StatusIsRequestHeaderFieldsTooLarge checks if the response status is 431 Request Header Fields Too Large.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsRequestHeaderFieldsTooLarge, func(r *http.Response) error{}) -> only matches when the status code is 431
var StatusIsRequestHeaderFieldsTooLarge = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusRequestHeaderFieldsTooLarge
}, 1)

// StatusIsUnavailableForLegalReasons checks if the response status is 451 Unavailable For Legal Reasons.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsUnavailableForLegalReasons, func(r *http.Response) error{}) -> only matches when the status code is 451
var StatusIsUnavailableForLegalReasons = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusUnavailableForLegalReasons
}, 1)

// 5xx Server Errors

// StatusIsInternalServerError checks if the response status is 500 Internal Server Error.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsInternalServerError, func(r *http.Response) error{}) -> only matches when the status code is 500
var StatusIsInternalServerError = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusInternalServerError
}, 1)

// StatusIsNotImplemented checks if the response status is 501 Not Implemented.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNotImplemented, func(r *http.Response) error{}) -> only matches when the status code is 501
var StatusIsNotImplemented = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotImplemented
}, 1)

// StatusIsBadGateway checks if the response status is 502 Bad Gateway.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsBadGateway, func(r *http.Response) error{}) -> only matches when the status code is 502
var StatusIsBadGateway = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusBadGateway
}, 1)

// StatusIsServiceUnavailable checks if the response status is 503 Service Unavailable.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsServiceUnavailable, func(r *http.Response) error{}) -> only matches when the status code is 503
var StatusIsServiceUnavailable = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusServiceUnavailable
}, 1)

// StatusIsGatewayTimeout checks if the response status is 504 Gateway Timeout.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsGatewayTimeout, func(r *http.Response) error{}) -> only matches when the status code is 504
var StatusIsGatewayTimeout = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusGatewayTimeout
}, 1)

// StatusIsHTTPVersionNotSupported checks if the response status is 505 HTTP Version Not Supported.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsHTTPVersionNotSupported, func(r *http.Response) error{}) -> only matches when the status code is 505
var StatusIsHTTPVersionNotSupported = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusHTTPVersionNotSupported
}, 1)

// StatusIsVariantAlsoNegotiates checks if the response status is 506 Variant Also Negotiates.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsVariantAlsoNegotiates, func(r *http.Response) error{}) -> only matches when the status code is 506
var StatusIsVariantAlsoNegotiates = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusVariantAlsoNegotiates
}, 1)

// StatusIsInsufficientStorage checks if the response status is 507 Insufficient Storage.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsInsufficientStorage, func(r *http.Response) error{}) -> only matches when the status code is 507
var StatusIsInsufficientStorage = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusInsufficientStorage
}, 1)

// StatusIsLoopDetected checks if the response status is 508 Loop Detected.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsLoopDetected, func(r *http.Response) error{}) -> only matches when the status code is 508
var StatusIsLoopDetected = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusLoopDetected
}, 1)

// StatusIsNotExtended checks if the response status is 510 Not Extended.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNotExtended, func(r *http.Response) error{}) -> only matches when the status code is 510
var StatusIsNotExtended = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNotExtended
}, 1)

// StatusIsNetworkAuthenticationRequired checks if the response status is 511 Network Authentication Required.
// It has the priority 1, and it has the top priority.
// Usage:
// OnReplyIf(StatusIsNetworkAuthenticationRequired, func(r *http.Response) error{}) -> only matches when the status code is 511
var StatusIsNetworkAuthenticationRequired = newStatusChecker(func(statusCode int) bool {
	return statusCode == http.StatusNetworkAuthenticationRequired
}, 1)

var Not = func(matcher StatusMatcher) *statusChecker {
	return &statusChecker{
		matcher: func(statusCode int) bool {
			return !matcher.Match(statusCode)
		},
		priority: matcher.Priority(),
	}
}

func DrainBodyAndClose(response *http.Response) error {
	if response != nil && response.Body != nil {
		ctx := response.Request.Context()
		logger := ExtractLoggerFromContext(ctx)
		defer func() {
			if err := response.Body.Close(); err != nil {
				logger.Error(ctx, err, "could not close the body")
			}
		}()
		// Limit drain to prevent memory issues with huge responses
		_, err := io.Copy(io.Discard, io.LimitReader(response.Body, 1<<20)) // 1MB limit
		if err != nil {
			logger.Error(ctx, err, "could not drain the body")

			return err
		}
	}

	return nil
}
