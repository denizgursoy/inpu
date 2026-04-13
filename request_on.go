package inpu

// OnOneOf is a shorthand for On(StatusIsOneOf(statusCodes...), responseHandler).
// It matches any of the provided status codes.
// Usage:
//
//	OnOneOf(ThenUnmarshalJsonTo(&items), http.StatusOK, http.StatusCreated).
//	Send()
func (r *Req) OnOneOf(responseHandler ResponseHandler, statusCodes ...int) *Req {
	return r.On(StatusIsOneOf(statusCodes...), responseHandler)
}

// OnAny is a shorthand for On(StatusAny, responseHandler).
// It matches any status code. Useful as a fallback.
// Usage:
//
//	OnOk(ThenUnmarshalJsonTo(&items)).
//	OnAny(ThenReturnDefaultError).
//	Send()
func (r *Req) OnAny(responseHandler ResponseHandler) *Req {
	return r.On(StatusAny, responseHandler)
}

// OnAnyExcept is a shorthand for On(StatusAnyExcept(statusCode), responseHandler).
// It matches any status code except the one provided.
// Usage:
//
//	OnAnyExcept(http.StatusOK, ThenReturnDefaultError).
//	Send()
func (r *Req) OnAnyExcept(statusCode int, responseHandler ResponseHandler) *Req {
	return r.On(StatusAnyExcept(statusCode), responseHandler)
}

// OnAnyExceptOneOf is a shorthand for On(StatusAnyExceptOneOf(statusCodes...), responseHandler).
// It matches any status code except those provided.
// Usage:
//
//	OnAnyExceptOneOf(ThenReturnDefaultError, http.StatusOK, http.StatusCreated).
//	Send()
func (r *Req) OnAnyExceptOneOf(responseHandler ResponseHandler, statusCodes ...int) *Req {
	return r.On(StatusAnyExceptOneOf(statusCodes...), responseHandler)
}

// --- Category matchers ---

// OnSuccess is a shorthand for On(StatusIsSuccess, responseHandler).
// It matches any status code in the range [200, 300).
// Usage:
//
//	OnSuccess(ThenUnmarshalJsonTo(&items)).
//	Send()
func (r *Req) OnSuccess(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsSuccess, responseHandler)
}

// OnInformational is a shorthand for On(StatusIsInformational, responseHandler).
// It matches any status code less than 200.
func (r *Req) OnInformational(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsInformational, responseHandler)
}

// OnRedirection is a shorthand for On(StatusIsRedirection, responseHandler).
// It matches any status code in the range [300, 400).
func (r *Req) OnRedirection(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsRedirection, responseHandler)
}

// OnClientError is a shorthand for On(StatusIsClientError, responseHandler).
// It matches any status code in the range [400, 500).
func (r *Req) OnClientError(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsClientError, responseHandler)
}

// OnServerError is a shorthand for On(StatusIsServerError, responseHandler).
// It matches any status code greater than or equal to 500.
func (r *Req) OnServerError(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsServerError, responseHandler)
}

// --- 1xx Informational ---

// OnContinue is a shorthand for On(StatusIsContinue, responseHandler).
// It matches status code 100.
func (r *Req) OnContinue(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsContinue, responseHandler)
}

// OnSwitchingProtocols is a shorthand for On(StatusIsSwitchingProtocols, responseHandler).
// It matches status code 101.
func (r *Req) OnSwitchingProtocols(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsSwitchingProtocols, responseHandler)
}

// OnProcessing is a shorthand for On(StatusIsProcessing, responseHandler).
// It matches status code 102.
func (r *Req) OnProcessing(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsProcessing, responseHandler)
}

// OnEarlyHints is a shorthand for On(StatusIsEarlyHints, responseHandler).
// It matches status code 103.
func (r *Req) OnEarlyHints(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsEarlyHints, responseHandler)
}

// --- 2xx Success ---

// OnOk is a shorthand for On(StatusIsOk, responseHandler).
// It matches status code 200.
// Usage:
//
//	OnOk(ThenUnmarshalJsonTo(&items)).
//	OnAny(ThenReturnDefaultError).
//	Send()
func (r *Req) OnOk(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsOk, responseHandler)
}

// OnCreated is a shorthand for On(StatusIsCreated, responseHandler).
// It matches status code 201.
func (r *Req) OnCreated(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsCreated, responseHandler)
}

// OnAccepted is a shorthand for On(StatusIsAccepted, responseHandler).
// It matches status code 202.
func (r *Req) OnAccepted(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsAccepted, responseHandler)
}

// OnNonAuthoritativeInfo is a shorthand for On(StatusIsNonAuthoritativeInfo, responseHandler).
// It matches status code 203.
func (r *Req) OnNonAuthoritativeInfo(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNonAuthoritativeInfo, responseHandler)
}

// OnNoContent is a shorthand for On(StatusIsNoContent, responseHandler).
// It matches status code 204.
func (r *Req) OnNoContent(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNoContent, responseHandler)
}

// OnResetContent is a shorthand for On(StatusIsResetContent, responseHandler).
// It matches status code 205.
func (r *Req) OnResetContent(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsResetContent, responseHandler)
}

// OnPartialContent is a shorthand for On(StatusIsPartialContent, responseHandler).
// It matches status code 206.
func (r *Req) OnPartialContent(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsPartialContent, responseHandler)
}

// OnMultiStatus is a shorthand for On(StatusIsMultiStatus, responseHandler).
// It matches status code 207.
func (r *Req) OnMultiStatus(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsMultiStatus, responseHandler)
}

// OnAlreadyReported is a shorthand for On(StatusIsAlreadyReported, responseHandler).
// It matches status code 208.
func (r *Req) OnAlreadyReported(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsAlreadyReported, responseHandler)
}

// OnIMUsed is a shorthand for On(StatusIsIMUsed, responseHandler).
// It matches status code 226.
func (r *Req) OnIMUsed(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsIMUsed, responseHandler)
}

// --- 3xx Redirection ---

// OnMultipleChoices is a shorthand for On(StatusIsMultipleChoices, responseHandler).
// It matches status code 300.
func (r *Req) OnMultipleChoices(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsMultipleChoices, responseHandler)
}

// OnMovedPermanently is a shorthand for On(StatusIsMovedPermanently, responseHandler).
// It matches status code 301.
func (r *Req) OnMovedPermanently(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsMovedPermanently, responseHandler)
}

// OnFound is a shorthand for On(StatusIsFound, responseHandler).
// It matches status code 302.
func (r *Req) OnFound(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsFound, responseHandler)
}

// OnSeeOther is a shorthand for On(StatusIsSeeOther, responseHandler).
// It matches status code 303.
func (r *Req) OnSeeOther(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsSeeOther, responseHandler)
}

// OnNotModified is a shorthand for On(StatusIsNotModified, responseHandler).
// It matches status code 304.
func (r *Req) OnNotModified(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNotModified, responseHandler)
}

// OnUseProxy is a shorthand for On(StatusIsUseProxy, responseHandler).
// It matches status code 305.
func (r *Req) OnUseProxy(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsUseProxy, responseHandler)
}

// OnTemporaryRedirect is a shorthand for On(StatusIsTemporaryRedirect, responseHandler).
// It matches status code 307.
func (r *Req) OnTemporaryRedirect(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsTemporaryRedirect, responseHandler)
}

// OnPermanentRedirect is a shorthand for On(StatusIsPermanentRedirect, responseHandler).
// It matches status code 308.
func (r *Req) OnPermanentRedirect(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsPermanentRedirect, responseHandler)
}

// --- 4xx Client Errors ---

// OnBadRequest is a shorthand for On(StatusIsBadRequest, responseHandler).
// It matches status code 400.
func (r *Req) OnBadRequest(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsBadRequest, responseHandler)
}

// OnUnauthorized is a shorthand for On(StatusIsUnauthorized, responseHandler).
// It matches status code 401.
func (r *Req) OnUnauthorized(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsUnauthorized, responseHandler)
}

// OnPaymentRequired is a shorthand for On(StatusIsPaymentRequired, responseHandler).
// It matches status code 402.
func (r *Req) OnPaymentRequired(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsPaymentRequired, responseHandler)
}

// OnForbidden is a shorthand for On(StatusIsForbidden, responseHandler).
// It matches status code 403.
func (r *Req) OnForbidden(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsForbidden, responseHandler)
}

// OnNotFound is a shorthand for On(StatusIsNotFound, responseHandler).
// It matches status code 404.
func (r *Req) OnNotFound(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNotFound, responseHandler)
}

// OnMethodNotAllowed is a shorthand for On(StatusIsMethodNotAllowed, responseHandler).
// It matches status code 405.
func (r *Req) OnMethodNotAllowed(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsMethodNotAllowed, responseHandler)
}

// OnNotAcceptable is a shorthand for On(StatusIsNotAcceptable, responseHandler).
// It matches status code 406.
func (r *Req) OnNotAcceptable(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNotAcceptable, responseHandler)
}

// OnProxyAuthRequired is a shorthand for On(StatusIsProxyAuthRequired, responseHandler).
// It matches status code 407.
func (r *Req) OnProxyAuthRequired(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsProxyAuthRequired, responseHandler)
}

// OnRequestTimeout is a shorthand for On(StatusIsRequestTimeout, responseHandler).
// It matches status code 408.
func (r *Req) OnRequestTimeout(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsRequestTimeout, responseHandler)
}

// OnConflict is a shorthand for On(StatusIsConflict, responseHandler).
// It matches status code 409.
func (r *Req) OnConflict(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsConflict, responseHandler)
}

// OnGone is a shorthand for On(StatusIsGone, responseHandler).
// It matches status code 410.
func (r *Req) OnGone(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsGone, responseHandler)
}

// OnLengthRequired is a shorthand for On(StatusIsLengthRequired, responseHandler).
// It matches status code 411.
func (r *Req) OnLengthRequired(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsLengthRequired, responseHandler)
}

// OnPreconditionFailed is a shorthand for On(StatusIsPreconditionFailed, responseHandler).
// It matches status code 412.
func (r *Req) OnPreconditionFailed(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsPreconditionFailed, responseHandler)
}

// OnRequestEntityTooLarge is a shorthand for On(StatusIsRequestEntityTooLarge, responseHandler).
// It matches status code 413.
func (r *Req) OnRequestEntityTooLarge(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsRequestEntityTooLarge, responseHandler)
}

// OnRequestURITooLong is a shorthand for On(StatusIsRequestURITooLong, responseHandler).
// It matches status code 414.
func (r *Req) OnRequestURITooLong(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsRequestURITooLong, responseHandler)
}

// OnUnsupportedMediaType is a shorthand for On(StatusIsUnsupportedMediaType, responseHandler).
// It matches status code 415.
func (r *Req) OnUnsupportedMediaType(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsUnsupportedMediaType, responseHandler)
}

// OnRequestedRangeNotSatisfiable is a shorthand for On(StatusIsRequestedRangeNotSatisfiable, responseHandler).
// It matches status code 416.
func (r *Req) OnRequestedRangeNotSatisfiable(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsRequestedRangeNotSatisfiable, responseHandler)
}

// OnExpectationFailed is a shorthand for On(StatusIsExpectationFailed, responseHandler).
// It matches status code 417.
func (r *Req) OnExpectationFailed(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsExpectationFailed, responseHandler)
}

// OnTeapot is a shorthand for On(StatusIsTeapot, responseHandler).
// It matches status code 418.
func (r *Req) OnTeapot(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsTeapot, responseHandler)
}

// OnMisdirectedRequest is a shorthand for On(StatusIsMisdirectedRequest, responseHandler).
// It matches status code 421.
func (r *Req) OnMisdirectedRequest(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsMisdirectedRequest, responseHandler)
}

// OnUnprocessableEntity is a shorthand for On(StatusIsUnprocessableEntity, responseHandler).
// It matches status code 422.
func (r *Req) OnUnprocessableEntity(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsUnprocessableEntity, responseHandler)
}

// OnLocked is a shorthand for On(StatusIsLocked, responseHandler).
// It matches status code 423.
func (r *Req) OnLocked(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsLocked, responseHandler)
}

// OnFailedDependency is a shorthand for On(StatusIsFailedDependency, responseHandler).
// It matches status code 424.
func (r *Req) OnFailedDependency(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsFailedDependency, responseHandler)
}

// OnTooEarly is a shorthand for On(StatusIsTooEarly, responseHandler).
// It matches status code 425.
func (r *Req) OnTooEarly(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsTooEarly, responseHandler)
}

// OnUpgradeRequired is a shorthand for On(StatusIsUpgradeRequired, responseHandler).
// It matches status code 426.
func (r *Req) OnUpgradeRequired(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsUpgradeRequired, responseHandler)
}

// OnPreconditionRequired is a shorthand for On(StatusIsPreconditionRequired, responseHandler).
// It matches status code 428.
func (r *Req) OnPreconditionRequired(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsPreconditionRequired, responseHandler)
}

// OnTooManyRequests is a shorthand for On(StatusIsTooManyRequests, responseHandler).
// It matches status code 429.
func (r *Req) OnTooManyRequests(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsTooManyRequests, responseHandler)
}

// OnRequestHeaderFieldsTooLarge is a shorthand for On(StatusIsRequestHeaderFieldsTooLarge, responseHandler).
// It matches status code 431.
func (r *Req) OnRequestHeaderFieldsTooLarge(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsRequestHeaderFieldsTooLarge, responseHandler)
}

// OnUnavailableForLegalReasons is a shorthand for On(StatusIsUnavailableForLegalReasons, responseHandler).
// It matches status code 451.
func (r *Req) OnUnavailableForLegalReasons(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsUnavailableForLegalReasons, responseHandler)
}

// --- 5xx Server Errors ---

// OnInternalServerError is a shorthand for On(StatusIsInternalServerError, responseHandler).
// It matches status code 500.
func (r *Req) OnInternalServerError(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsInternalServerError, responseHandler)
}

// OnNotImplemented is a shorthand for On(StatusIsNotImplemented, responseHandler).
// It matches status code 501.
func (r *Req) OnNotImplemented(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNotImplemented, responseHandler)
}

// OnBadGateway is a shorthand for On(StatusIsBadGateway, responseHandler).
// It matches status code 502.
func (r *Req) OnBadGateway(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsBadGateway, responseHandler)
}

// OnServiceUnavailable is a shorthand for On(StatusIsServiceUnavailable, responseHandler).
// It matches status code 503.
func (r *Req) OnServiceUnavailable(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsServiceUnavailable, responseHandler)
}

// OnGatewayTimeout is a shorthand for On(StatusIsGatewayTimeout, responseHandler).
// It matches status code 504.
func (r *Req) OnGatewayTimeout(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsGatewayTimeout, responseHandler)
}

// OnHTTPVersionNotSupported is a shorthand for On(StatusIsHTTPVersionNotSupported, responseHandler).
// It matches status code 505.
func (r *Req) OnHTTPVersionNotSupported(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsHTTPVersionNotSupported, responseHandler)
}

// OnVariantAlsoNegotiates is a shorthand for On(StatusIsVariantAlsoNegotiates, responseHandler).
// It matches status code 506.
func (r *Req) OnVariantAlsoNegotiates(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsVariantAlsoNegotiates, responseHandler)
}

// OnInsufficientStorage is a shorthand for On(StatusIsInsufficientStorage, responseHandler).
// It matches status code 507.
func (r *Req) OnInsufficientStorage(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsInsufficientStorage, responseHandler)
}

// OnLoopDetected is a shorthand for On(StatusIsLoopDetected, responseHandler).
// It matches status code 508.
func (r *Req) OnLoopDetected(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsLoopDetected, responseHandler)
}

// OnNotExtended is a shorthand for On(StatusIsNotExtended, responseHandler).
// It matches status code 510.
func (r *Req) OnNotExtended(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNotExtended, responseHandler)
}

// OnNetworkAuthenticationRequired is a shorthand for On(StatusIsNetworkAuthenticationRequired, responseHandler).
// It matches status code 511.
func (r *Req) OnNetworkAuthenticationRequired(responseHandler ResponseHandler) *Req {
	return r.On(StatusIsNetworkAuthenticationRequired, responseHandler)
}
