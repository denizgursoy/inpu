package inpu

// On is a shorthand for OnReplyIf(StatusIs(statusCode), responseHandler).
// It matches a single exact status code.
// Usage:
//
//	On(http.StatusOK, ThenUnmarshalJsonTo(&items)).
//	On(http.StatusNotFound, ThenReturnError(ErrNotFound)).
//	Send()
func (r *Req) On(statusCode int, responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIs(statusCode), responseHandler)
}

// OnOneOf is a shorthand for OnReplyIf(StatusIsOneOf(statusCodes...), responseHandler).
// It matches any of the provided status codes.
// Usage:
//
//	OnOneOf(ThenUnmarshalJsonTo(&items), http.StatusOK, http.StatusCreated).
//	Send()
func (r *Req) OnOneOf(responseHandler ResponseHandler, statusCodes ...int) *Req {
	return r.OnReplyIf(StatusIsOneOf(statusCodes...), responseHandler)
}

// OnAny is a shorthand for OnReplyIf(StatusAny, responseHandler).
// It matches any status code. Useful as a fallback.
// Usage:
//
//	OnOk(ThenUnmarshalJsonTo(&items)).
//	OnAny(ThenReturnDefaultError).
//	Send()
func (r *Req) OnAny(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusAny, responseHandler)
}

// OnAnyExcept is a shorthand for OnReplyIf(StatusAnyExcept(statusCode), responseHandler).
// It matches any status code except the one provided.
// Usage:
//
//	OnAnyExcept(http.StatusOK, ThenReturnDefaultError).
//	Send()
func (r *Req) OnAnyExcept(statusCode int, responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusAnyExcept(statusCode), responseHandler)
}

// OnAnyExceptOneOf is a shorthand for OnReplyIf(StatusAnyExceptOneOf(statusCodes...), responseHandler).
// It matches any status code except those provided.
// Usage:
//
//	OnAnyExceptOneOf(ThenReturnDefaultError, http.StatusOK, http.StatusCreated).
//	Send()
func (r *Req) OnAnyExceptOneOf(responseHandler ResponseHandler, statusCodes ...int) *Req {
	return r.OnReplyIf(StatusAnyExceptOneOf(statusCodes...), responseHandler)
}

// --- Category matchers ---

// OnSuccess is a shorthand for OnReplyIf(StatusIsSuccess, responseHandler).
// It matches any status code in the range [200, 300).
// Usage:
//
//	OnSuccess(ThenUnmarshalJsonTo(&items)).
//	Send()
func (r *Req) OnSuccess(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsSuccess, responseHandler)
}

// OnInformational is a shorthand for OnReplyIf(StatusIsInformational, responseHandler).
// It matches any status code less than 200.
func (r *Req) OnInformational(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsInformational, responseHandler)
}

// OnRedirection is a shorthand for OnReplyIf(StatusIsRedirection, responseHandler).
// It matches any status code in the range [300, 400).
func (r *Req) OnRedirection(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsRedirection, responseHandler)
}

// OnClientError is a shorthand for OnReplyIf(StatusIsClientError, responseHandler).
// It matches any status code in the range [400, 500).
func (r *Req) OnClientError(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsClientError, responseHandler)
}

// OnServerError is a shorthand for OnReplyIf(StatusIsServerError, responseHandler).
// It matches any status code greater than or equal to 500.
func (r *Req) OnServerError(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsServerError, responseHandler)
}

// --- 1xx Informational ---

// OnContinue is a shorthand for OnReplyIf(StatusIsContinue, responseHandler).
// It matches status code 100.
func (r *Req) OnContinue(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsContinue, responseHandler)
}

// OnSwitchingProtocols is a shorthand for OnReplyIf(StatusIsSwitchingProtocols, responseHandler).
// It matches status code 101.
func (r *Req) OnSwitchingProtocols(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsSwitchingProtocols, responseHandler)
}

// OnProcessing is a shorthand for OnReplyIf(StatusIsProcessing, responseHandler).
// It matches status code 102.
func (r *Req) OnProcessing(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsProcessing, responseHandler)
}

// OnEarlyHints is a shorthand for OnReplyIf(StatusIsEarlyHints, responseHandler).
// It matches status code 103.
func (r *Req) OnEarlyHints(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsEarlyHints, responseHandler)
}

// --- 2xx Success ---

// OnOk is a shorthand for OnReplyIf(StatusIsOk, responseHandler).
// It matches status code 200.
// Usage:
//
//	OnOk(ThenUnmarshalJsonTo(&items)).
//	OnAny(ThenReturnDefaultError).
//	Send()
func (r *Req) OnOk(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsOk, responseHandler)
}

// OnCreated is a shorthand for OnReplyIf(StatusIsCreated, responseHandler).
// It matches status code 201.
func (r *Req) OnCreated(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsCreated, responseHandler)
}

// OnAccepted is a shorthand for OnReplyIf(StatusIsAccepted, responseHandler).
// It matches status code 202.
func (r *Req) OnAccepted(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsAccepted, responseHandler)
}

// OnNonAuthoritativeInfo is a shorthand for OnReplyIf(StatusIsNonAuthoritativeInfo, responseHandler).
// It matches status code 203.
func (r *Req) OnNonAuthoritativeInfo(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNonAuthoritativeInfo, responseHandler)
}

// OnNoContent is a shorthand for OnReplyIf(StatusIsNoContent, responseHandler).
// It matches status code 204.
func (r *Req) OnNoContent(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNoContent, responseHandler)
}

// OnResetContent is a shorthand for OnReplyIf(StatusIsResetContent, responseHandler).
// It matches status code 205.
func (r *Req) OnResetContent(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsResetContent, responseHandler)
}

// OnPartialContent is a shorthand for OnReplyIf(StatusIsPartialContent, responseHandler).
// It matches status code 206.
func (r *Req) OnPartialContent(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsPartialContent, responseHandler)
}

// OnMultiStatus is a shorthand for OnReplyIf(StatusIsMultiStatus, responseHandler).
// It matches status code 207.
func (r *Req) OnMultiStatus(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsMultiStatus, responseHandler)
}

// OnAlreadyReported is a shorthand for OnReplyIf(StatusIsAlreadyReported, responseHandler).
// It matches status code 208.
func (r *Req) OnAlreadyReported(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsAlreadyReported, responseHandler)
}

// OnIMUsed is a shorthand for OnReplyIf(StatusIsIMUsed, responseHandler).
// It matches status code 226.
func (r *Req) OnIMUsed(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsIMUsed, responseHandler)
}

// --- 3xx Redirection ---

// OnMultipleChoices is a shorthand for OnReplyIf(StatusIsMultipleChoices, responseHandler).
// It matches status code 300.
func (r *Req) OnMultipleChoices(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsMultipleChoices, responseHandler)
}

// OnMovedPermanently is a shorthand for OnReplyIf(StatusIsMovedPermanently, responseHandler).
// It matches status code 301.
func (r *Req) OnMovedPermanently(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsMovedPermanently, responseHandler)
}

// OnFound is a shorthand for OnReplyIf(StatusIsFound, responseHandler).
// It matches status code 302.
func (r *Req) OnFound(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsFound, responseHandler)
}

// OnSeeOther is a shorthand for OnReplyIf(StatusIsSeeOther, responseHandler).
// It matches status code 303.
func (r *Req) OnSeeOther(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsSeeOther, responseHandler)
}

// OnNotModified is a shorthand for OnReplyIf(StatusIsNotModified, responseHandler).
// It matches status code 304.
func (r *Req) OnNotModified(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNotModified, responseHandler)
}

// OnUseProxy is a shorthand for OnReplyIf(StatusIsUseProxy, responseHandler).
// It matches status code 305.
func (r *Req) OnUseProxy(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsUseProxy, responseHandler)
}

// OnTemporaryRedirect is a shorthand for OnReplyIf(StatusIsTemporaryRedirect, responseHandler).
// It matches status code 307.
func (r *Req) OnTemporaryRedirect(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsTemporaryRedirect, responseHandler)
}

// OnPermanentRedirect is a shorthand for OnReplyIf(StatusIsPermanentRedirect, responseHandler).
// It matches status code 308.
func (r *Req) OnPermanentRedirect(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsPermanentRedirect, responseHandler)
}

// --- 4xx Client Errors ---

// OnBadRequest is a shorthand for OnReplyIf(StatusIsBadRequest, responseHandler).
// It matches status code 400.
func (r *Req) OnBadRequest(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsBadRequest, responseHandler)
}

// OnUnauthorized is a shorthand for OnReplyIf(StatusIsUnauthorized, responseHandler).
// It matches status code 401.
func (r *Req) OnUnauthorized(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsUnauthorized, responseHandler)
}

// OnPaymentRequired is a shorthand for OnReplyIf(StatusIsPaymentRequired, responseHandler).
// It matches status code 402.
func (r *Req) OnPaymentRequired(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsPaymentRequired, responseHandler)
}

// OnForbidden is a shorthand for OnReplyIf(StatusIsForbidden, responseHandler).
// It matches status code 403.
func (r *Req) OnForbidden(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsForbidden, responseHandler)
}

// OnNotFound is a shorthand for OnReplyIf(StatusIsNotFound, responseHandler).
// It matches status code 404.
func (r *Req) OnNotFound(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNotFound, responseHandler)
}

// OnMethodNotAllowed is a shorthand for OnReplyIf(StatusIsMethodNotAllowed, responseHandler).
// It matches status code 405.
func (r *Req) OnMethodNotAllowed(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsMethodNotAllowed, responseHandler)
}

// OnNotAcceptable is a shorthand for OnReplyIf(StatusIsNotAcceptable, responseHandler).
// It matches status code 406.
func (r *Req) OnNotAcceptable(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNotAcceptable, responseHandler)
}

// OnProxyAuthRequired is a shorthand for OnReplyIf(StatusIsProxyAuthRequired, responseHandler).
// It matches status code 407.
func (r *Req) OnProxyAuthRequired(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsProxyAuthRequired, responseHandler)
}

// OnRequestTimeout is a shorthand for OnReplyIf(StatusIsRequestTimeout, responseHandler).
// It matches status code 408.
func (r *Req) OnRequestTimeout(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsRequestTimeout, responseHandler)
}

// OnConflict is a shorthand for OnReplyIf(StatusIsConflict, responseHandler).
// It matches status code 409.
func (r *Req) OnConflict(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsConflict, responseHandler)
}

// OnGone is a shorthand for OnReplyIf(StatusIsGone, responseHandler).
// It matches status code 410.
func (r *Req) OnGone(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsGone, responseHandler)
}

// OnLengthRequired is a shorthand for OnReplyIf(StatusIsLengthRequired, responseHandler).
// It matches status code 411.
func (r *Req) OnLengthRequired(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsLengthRequired, responseHandler)
}

// OnPreconditionFailed is a shorthand for OnReplyIf(StatusIsPreconditionFailed, responseHandler).
// It matches status code 412.
func (r *Req) OnPreconditionFailed(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsPreconditionFailed, responseHandler)
}

// OnRequestEntityTooLarge is a shorthand for OnReplyIf(StatusIsRequestEntityTooLarge, responseHandler).
// It matches status code 413.
func (r *Req) OnRequestEntityTooLarge(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsRequestEntityTooLarge, responseHandler)
}

// OnRequestURITooLong is a shorthand for OnReplyIf(StatusIsRequestURITooLong, responseHandler).
// It matches status code 414.
func (r *Req) OnRequestURITooLong(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsRequestURITooLong, responseHandler)
}

// OnUnsupportedMediaType is a shorthand for OnReplyIf(StatusIsUnsupportedMediaType, responseHandler).
// It matches status code 415.
func (r *Req) OnUnsupportedMediaType(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsUnsupportedMediaType, responseHandler)
}

// OnRequestedRangeNotSatisfiable is a shorthand for OnReplyIf(StatusIsRequestedRangeNotSatisfiable, responseHandler).
// It matches status code 416.
func (r *Req) OnRequestedRangeNotSatisfiable(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsRequestedRangeNotSatisfiable, responseHandler)
}

// OnExpectationFailed is a shorthand for OnReplyIf(StatusIsExpectationFailed, responseHandler).
// It matches status code 417.
func (r *Req) OnExpectationFailed(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsExpectationFailed, responseHandler)
}

// OnTeapot is a shorthand for OnReplyIf(StatusIsTeapot, responseHandler).
// It matches status code 418.
func (r *Req) OnTeapot(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsTeapot, responseHandler)
}

// OnMisdirectedRequest is a shorthand for OnReplyIf(StatusIsMisdirectedRequest, responseHandler).
// It matches status code 421.
func (r *Req) OnMisdirectedRequest(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsMisdirectedRequest, responseHandler)
}

// OnUnprocessableEntity is a shorthand for OnReplyIf(StatusIsUnprocessableEntity, responseHandler).
// It matches status code 422.
func (r *Req) OnUnprocessableEntity(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsUnprocessableEntity, responseHandler)
}

// OnLocked is a shorthand for OnReplyIf(StatusIsLocked, responseHandler).
// It matches status code 423.
func (r *Req) OnLocked(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsLocked, responseHandler)
}

// OnFailedDependency is a shorthand for OnReplyIf(StatusIsFailedDependency, responseHandler).
// It matches status code 424.
func (r *Req) OnFailedDependency(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsFailedDependency, responseHandler)
}

// OnTooEarly is a shorthand for OnReplyIf(StatusIsTooEarly, responseHandler).
// It matches status code 425.
func (r *Req) OnTooEarly(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsTooEarly, responseHandler)
}

// OnUpgradeRequired is a shorthand for OnReplyIf(StatusIsUpgradeRequired, responseHandler).
// It matches status code 426.
func (r *Req) OnUpgradeRequired(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsUpgradeRequired, responseHandler)
}

// OnPreconditionRequired is a shorthand for OnReplyIf(StatusIsPreconditionRequired, responseHandler).
// It matches status code 428.
func (r *Req) OnPreconditionRequired(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsPreconditionRequired, responseHandler)
}

// OnTooManyRequests is a shorthand for OnReplyIf(StatusIsTooManyRequests, responseHandler).
// It matches status code 429.
func (r *Req) OnTooManyRequests(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsTooManyRequests, responseHandler)
}

// OnRequestHeaderFieldsTooLarge is a shorthand for OnReplyIf(StatusIsRequestHeaderFieldsTooLarge, responseHandler).
// It matches status code 431.
func (r *Req) OnRequestHeaderFieldsTooLarge(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsRequestHeaderFieldsTooLarge, responseHandler)
}

// OnUnavailableForLegalReasons is a shorthand for OnReplyIf(StatusIsUnavailableForLegalReasons, responseHandler).
// It matches status code 451.
func (r *Req) OnUnavailableForLegalReasons(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsUnavailableForLegalReasons, responseHandler)
}

// --- 5xx Server Errors ---

// OnInternalServerError is a shorthand for OnReplyIf(StatusIsInternalServerError, responseHandler).
// It matches status code 500.
func (r *Req) OnInternalServerError(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsInternalServerError, responseHandler)
}

// OnNotImplemented is a shorthand for OnReplyIf(StatusIsNotImplemented, responseHandler).
// It matches status code 501.
func (r *Req) OnNotImplemented(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNotImplemented, responseHandler)
}

// OnBadGateway is a shorthand for OnReplyIf(StatusIsBadGateway, responseHandler).
// It matches status code 502.
func (r *Req) OnBadGateway(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsBadGateway, responseHandler)
}

// OnServiceUnavailable is a shorthand for OnReplyIf(StatusIsServiceUnavailable, responseHandler).
// It matches status code 503.
func (r *Req) OnServiceUnavailable(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsServiceUnavailable, responseHandler)
}

// OnGatewayTimeout is a shorthand for OnReplyIf(StatusIsGatewayTimeout, responseHandler).
// It matches status code 504.
func (r *Req) OnGatewayTimeout(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsGatewayTimeout, responseHandler)
}

// OnHTTPVersionNotSupported is a shorthand for OnReplyIf(StatusIsHTTPVersionNotSupported, responseHandler).
// It matches status code 505.
func (r *Req) OnHTTPVersionNotSupported(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsHTTPVersionNotSupported, responseHandler)
}

// OnVariantAlsoNegotiates is a shorthand for OnReplyIf(StatusIsVariantAlsoNegotiates, responseHandler).
// It matches status code 506.
func (r *Req) OnVariantAlsoNegotiates(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsVariantAlsoNegotiates, responseHandler)
}

// OnInsufficientStorage is a shorthand for OnReplyIf(StatusIsInsufficientStorage, responseHandler).
// It matches status code 507.
func (r *Req) OnInsufficientStorage(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsInsufficientStorage, responseHandler)
}

// OnLoopDetected is a shorthand for OnReplyIf(StatusIsLoopDetected, responseHandler).
// It matches status code 508.
func (r *Req) OnLoopDetected(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsLoopDetected, responseHandler)
}

// OnNotExtended is a shorthand for OnReplyIf(StatusIsNotExtended, responseHandler).
// It matches status code 510.
func (r *Req) OnNotExtended(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNotExtended, responseHandler)
}

// OnNetworkAuthenticationRequired is a shorthand for OnReplyIf(StatusIsNetworkAuthenticationRequired, responseHandler).
// It matches status code 511.
func (r *Req) OnNetworkAuthenticationRequired(responseHandler ResponseHandler) *Req {
	return r.OnReplyIf(StatusIsNetworkAuthenticationRequired, responseHandler)
}
