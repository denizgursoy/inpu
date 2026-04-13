package inpu

// On is a shorthand for OnWhen(StatusIs(statusCode), responseHandler).
// It matches a single exact status code.
// Usage:
//
//	On(http.StatusOK, ThenUnmarshalJsonTo(&items)).
//	On(http.StatusNotFound, ThenReturnError(ErrNotFound)).
//	Send()
func (r *Req) On(statusCode int, responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIs(statusCode), responseHandler)
}

// OnOneOf is a shorthand for OnWhen(StatusIsOneOf(statusCodes...), responseHandler).
// It matches any of the provided status codes.
// Usage:
//
//	OnOneOf(ThenUnmarshalJsonTo(&items), http.StatusOK, http.StatusCreated).
//	Send()
func (r *Req) OnOneOf(responseHandler ResponseHandler, statusCodes ...int) *Req {
	return r.OnWhen(StatusIsOneOf(statusCodes...), responseHandler)
}

// OnAny is a shorthand for OnWhen(StatusAny, responseHandler).
// It matches any status code. Useful as a fallback.
// Usage:
//
//	OnOk(ThenUnmarshalJsonTo(&items)).
//	OnAny(ThenReturnDefaultError).
//	Send()
func (r *Req) OnAny(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusAny, responseHandler)
}

// OnAnyExcept is a shorthand for OnWhen(StatusAnyExcept(statusCode), responseHandler).
// It matches any status code except the one provided.
// Usage:
//
//	OnAnyExcept(http.StatusOK, ThenReturnDefaultError).
//	Send()
func (r *Req) OnAnyExcept(statusCode int, responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusAnyExcept(statusCode), responseHandler)
}

// OnAnyExceptOneOf is a shorthand for OnWhen(StatusAnyExceptOneOf(statusCodes...), responseHandler).
// It matches any status code except those provided.
// Usage:
//
//	OnAnyExceptOneOf(ThenReturnDefaultError, http.StatusOK, http.StatusCreated).
//	Send()
func (r *Req) OnAnyExceptOneOf(responseHandler ResponseHandler, statusCodes ...int) *Req {
	return r.OnWhen(StatusAnyExceptOneOf(statusCodes...), responseHandler)
}

// --- Category matchers ---

// OnSuccess is a shorthand for OnWhen(StatusIsSuccess, responseHandler).
// It matches any status code in the range [200, 300).
// Usage:
//
//	OnSuccess(ThenUnmarshalJsonTo(&items)).
//	Send()
func (r *Req) OnSuccess(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsSuccess, responseHandler)
}

// OnInformational is a shorthand for OnWhen(StatusIsInformational, responseHandler).
// It matches any status code less than 200.
func (r *Req) OnInformational(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsInformational, responseHandler)
}

// OnRedirection is a shorthand for OnWhen(StatusIsRedirection, responseHandler).
// It matches any status code in the range [300, 400).
func (r *Req) OnRedirection(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsRedirection, responseHandler)
}

// OnClientError is a shorthand for OnWhen(StatusIsClientError, responseHandler).
// It matches any status code in the range [400, 500).
func (r *Req) OnClientError(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsClientError, responseHandler)
}

// OnServerError is a shorthand for OnWhen(StatusIsServerError, responseHandler).
// It matches any status code greater than or equal to 500.
func (r *Req) OnServerError(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsServerError, responseHandler)
}

// --- 1xx Informational ---

// OnContinue is a shorthand for OnWhen(StatusIsContinue, responseHandler).
// It matches status code 100.
func (r *Req) OnContinue(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsContinue, responseHandler)
}

// OnSwitchingProtocols is a shorthand for OnWhen(StatusIsSwitchingProtocols, responseHandler).
// It matches status code 101.
func (r *Req) OnSwitchingProtocols(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsSwitchingProtocols, responseHandler)
}

// OnProcessing is a shorthand for OnWhen(StatusIsProcessing, responseHandler).
// It matches status code 102.
func (r *Req) OnProcessing(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsProcessing, responseHandler)
}

// OnEarlyHints is a shorthand for OnWhen(StatusIsEarlyHints, responseHandler).
// It matches status code 103.
func (r *Req) OnEarlyHints(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsEarlyHints, responseHandler)
}

// --- 2xx Success ---

// OnOk is a shorthand for OnWhen(StatusIsOk, responseHandler).
// It matches status code 200.
// Usage:
//
//	OnOk(ThenUnmarshalJsonTo(&items)).
//	OnAny(ThenReturnDefaultError).
//	Send()
func (r *Req) OnOk(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsOk, responseHandler)
}

// OnCreated is a shorthand for OnWhen(StatusIsCreated, responseHandler).
// It matches status code 201.
func (r *Req) OnCreated(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsCreated, responseHandler)
}

// OnAccepted is a shorthand for OnWhen(StatusIsAccepted, responseHandler).
// It matches status code 202.
func (r *Req) OnAccepted(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsAccepted, responseHandler)
}

// OnNonAuthoritativeInfo is a shorthand for OnWhen(StatusIsNonAuthoritativeInfo, responseHandler).
// It matches status code 203.
func (r *Req) OnNonAuthoritativeInfo(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNonAuthoritativeInfo, responseHandler)
}

// OnNoContent is a shorthand for OnWhen(StatusIsNoContent, responseHandler).
// It matches status code 204.
func (r *Req) OnNoContent(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNoContent, responseHandler)
}

// OnResetContent is a shorthand for OnWhen(StatusIsResetContent, responseHandler).
// It matches status code 205.
func (r *Req) OnResetContent(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsResetContent, responseHandler)
}

// OnPartialContent is a shorthand for OnWhen(StatusIsPartialContent, responseHandler).
// It matches status code 206.
func (r *Req) OnPartialContent(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsPartialContent, responseHandler)
}

// OnMultiStatus is a shorthand for OnWhen(StatusIsMultiStatus, responseHandler).
// It matches status code 207.
func (r *Req) OnMultiStatus(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsMultiStatus, responseHandler)
}

// OnAlreadyReported is a shorthand for OnWhen(StatusIsAlreadyReported, responseHandler).
// It matches status code 208.
func (r *Req) OnAlreadyReported(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsAlreadyReported, responseHandler)
}

// OnIMUsed is a shorthand for OnWhen(StatusIsIMUsed, responseHandler).
// It matches status code 226.
func (r *Req) OnIMUsed(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsIMUsed, responseHandler)
}

// --- 3xx Redirection ---

// OnMultipleChoices is a shorthand for OnWhen(StatusIsMultipleChoices, responseHandler).
// It matches status code 300.
func (r *Req) OnMultipleChoices(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsMultipleChoices, responseHandler)
}

// OnMovedPermanently is a shorthand for OnWhen(StatusIsMovedPermanently, responseHandler).
// It matches status code 301.
func (r *Req) OnMovedPermanently(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsMovedPermanently, responseHandler)
}

// OnFound is a shorthand for OnWhen(StatusIsFound, responseHandler).
// It matches status code 302.
func (r *Req) OnFound(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsFound, responseHandler)
}

// OnSeeOther is a shorthand for OnWhen(StatusIsSeeOther, responseHandler).
// It matches status code 303.
func (r *Req) OnSeeOther(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsSeeOther, responseHandler)
}

// OnNotModified is a shorthand for OnWhen(StatusIsNotModified, responseHandler).
// It matches status code 304.
func (r *Req) OnNotModified(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNotModified, responseHandler)
}

// OnUseProxy is a shorthand for OnWhen(StatusIsUseProxy, responseHandler).
// It matches status code 305.
func (r *Req) OnUseProxy(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsUseProxy, responseHandler)
}

// OnTemporaryRedirect is a shorthand for OnWhen(StatusIsTemporaryRedirect, responseHandler).
// It matches status code 307.
func (r *Req) OnTemporaryRedirect(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsTemporaryRedirect, responseHandler)
}

// OnPermanentRedirect is a shorthand for OnWhen(StatusIsPermanentRedirect, responseHandler).
// It matches status code 308.
func (r *Req) OnPermanentRedirect(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsPermanentRedirect, responseHandler)
}

// --- 4xx Client Errors ---

// OnBadRequest is a shorthand for OnWhen(StatusIsBadRequest, responseHandler).
// It matches status code 400.
func (r *Req) OnBadRequest(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsBadRequest, responseHandler)
}

// OnUnauthorized is a shorthand for OnWhen(StatusIsUnauthorized, responseHandler).
// It matches status code 401.
func (r *Req) OnUnauthorized(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsUnauthorized, responseHandler)
}

// OnPaymentRequired is a shorthand for OnWhen(StatusIsPaymentRequired, responseHandler).
// It matches status code 402.
func (r *Req) OnPaymentRequired(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsPaymentRequired, responseHandler)
}

// OnForbidden is a shorthand for OnWhen(StatusIsForbidden, responseHandler).
// It matches status code 403.
func (r *Req) OnForbidden(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsForbidden, responseHandler)
}

// OnNotFound is a shorthand for OnWhen(StatusIsNotFound, responseHandler).
// It matches status code 404.
func (r *Req) OnNotFound(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNotFound, responseHandler)
}

// OnMethodNotAllowed is a shorthand for OnWhen(StatusIsMethodNotAllowed, responseHandler).
// It matches status code 405.
func (r *Req) OnMethodNotAllowed(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsMethodNotAllowed, responseHandler)
}

// OnNotAcceptable is a shorthand for OnWhen(StatusIsNotAcceptable, responseHandler).
// It matches status code 406.
func (r *Req) OnNotAcceptable(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNotAcceptable, responseHandler)
}

// OnProxyAuthRequired is a shorthand for OnWhen(StatusIsProxyAuthRequired, responseHandler).
// It matches status code 407.
func (r *Req) OnProxyAuthRequired(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsProxyAuthRequired, responseHandler)
}

// OnRequestTimeout is a shorthand for OnWhen(StatusIsRequestTimeout, responseHandler).
// It matches status code 408.
func (r *Req) OnRequestTimeout(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsRequestTimeout, responseHandler)
}

// OnConflict is a shorthand for OnWhen(StatusIsConflict, responseHandler).
// It matches status code 409.
func (r *Req) OnConflict(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsConflict, responseHandler)
}

// OnGone is a shorthand for OnWhen(StatusIsGone, responseHandler).
// It matches status code 410.
func (r *Req) OnGone(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsGone, responseHandler)
}

// OnLengthRequired is a shorthand for OnWhen(StatusIsLengthRequired, responseHandler).
// It matches status code 411.
func (r *Req) OnLengthRequired(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsLengthRequired, responseHandler)
}

// OnPreconditionFailed is a shorthand for OnWhen(StatusIsPreconditionFailed, responseHandler).
// It matches status code 412.
func (r *Req) OnPreconditionFailed(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsPreconditionFailed, responseHandler)
}

// OnRequestEntityTooLarge is a shorthand for OnWhen(StatusIsRequestEntityTooLarge, responseHandler).
// It matches status code 413.
func (r *Req) OnRequestEntityTooLarge(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsRequestEntityTooLarge, responseHandler)
}

// OnRequestURITooLong is a shorthand for OnWhen(StatusIsRequestURITooLong, responseHandler).
// It matches status code 414.
func (r *Req) OnRequestURITooLong(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsRequestURITooLong, responseHandler)
}

// OnUnsupportedMediaType is a shorthand for OnWhen(StatusIsUnsupportedMediaType, responseHandler).
// It matches status code 415.
func (r *Req) OnUnsupportedMediaType(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsUnsupportedMediaType, responseHandler)
}

// OnRequestedRangeNotSatisfiable is a shorthand for OnWhen(StatusIsRequestedRangeNotSatisfiable, responseHandler).
// It matches status code 416.
func (r *Req) OnRequestedRangeNotSatisfiable(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsRequestedRangeNotSatisfiable, responseHandler)
}

// OnExpectationFailed is a shorthand for OnWhen(StatusIsExpectationFailed, responseHandler).
// It matches status code 417.
func (r *Req) OnExpectationFailed(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsExpectationFailed, responseHandler)
}

// OnTeapot is a shorthand for OnWhen(StatusIsTeapot, responseHandler).
// It matches status code 418.
func (r *Req) OnTeapot(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsTeapot, responseHandler)
}

// OnMisdirectedRequest is a shorthand for OnWhen(StatusIsMisdirectedRequest, responseHandler).
// It matches status code 421.
func (r *Req) OnMisdirectedRequest(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsMisdirectedRequest, responseHandler)
}

// OnUnprocessableEntity is a shorthand for OnWhen(StatusIsUnprocessableEntity, responseHandler).
// It matches status code 422.
func (r *Req) OnUnprocessableEntity(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsUnprocessableEntity, responseHandler)
}

// OnLocked is a shorthand for OnWhen(StatusIsLocked, responseHandler).
// It matches status code 423.
func (r *Req) OnLocked(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsLocked, responseHandler)
}

// OnFailedDependency is a shorthand for OnWhen(StatusIsFailedDependency, responseHandler).
// It matches status code 424.
func (r *Req) OnFailedDependency(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsFailedDependency, responseHandler)
}

// OnTooEarly is a shorthand for OnWhen(StatusIsTooEarly, responseHandler).
// It matches status code 425.
func (r *Req) OnTooEarly(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsTooEarly, responseHandler)
}

// OnUpgradeRequired is a shorthand for OnWhen(StatusIsUpgradeRequired, responseHandler).
// It matches status code 426.
func (r *Req) OnUpgradeRequired(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsUpgradeRequired, responseHandler)
}

// OnPreconditionRequired is a shorthand for OnWhen(StatusIsPreconditionRequired, responseHandler).
// It matches status code 428.
func (r *Req) OnPreconditionRequired(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsPreconditionRequired, responseHandler)
}

// OnTooManyRequests is a shorthand for OnWhen(StatusIsTooManyRequests, responseHandler).
// It matches status code 429.
func (r *Req) OnTooManyRequests(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsTooManyRequests, responseHandler)
}

// OnRequestHeaderFieldsTooLarge is a shorthand for OnWhen(StatusIsRequestHeaderFieldsTooLarge, responseHandler).
// It matches status code 431.
func (r *Req) OnRequestHeaderFieldsTooLarge(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsRequestHeaderFieldsTooLarge, responseHandler)
}

// OnUnavailableForLegalReasons is a shorthand for OnWhen(StatusIsUnavailableForLegalReasons, responseHandler).
// It matches status code 451.
func (r *Req) OnUnavailableForLegalReasons(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsUnavailableForLegalReasons, responseHandler)
}

// --- 5xx Server Errors ---

// OnInternalServerError is a shorthand for OnWhen(StatusIsInternalServerError, responseHandler).
// It matches status code 500.
func (r *Req) OnInternalServerError(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsInternalServerError, responseHandler)
}

// OnNotImplemented is a shorthand for OnWhen(StatusIsNotImplemented, responseHandler).
// It matches status code 501.
func (r *Req) OnNotImplemented(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNotImplemented, responseHandler)
}

// OnBadGateway is a shorthand for OnWhen(StatusIsBadGateway, responseHandler).
// It matches status code 502.
func (r *Req) OnBadGateway(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsBadGateway, responseHandler)
}

// OnServiceUnavailable is a shorthand for OnWhen(StatusIsServiceUnavailable, responseHandler).
// It matches status code 503.
func (r *Req) OnServiceUnavailable(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsServiceUnavailable, responseHandler)
}

// OnGatewayTimeout is a shorthand for OnWhen(StatusIsGatewayTimeout, responseHandler).
// It matches status code 504.
func (r *Req) OnGatewayTimeout(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsGatewayTimeout, responseHandler)
}

// OnHTTPVersionNotSupported is a shorthand for OnWhen(StatusIsHTTPVersionNotSupported, responseHandler).
// It matches status code 505.
func (r *Req) OnHTTPVersionNotSupported(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsHTTPVersionNotSupported, responseHandler)
}

// OnVariantAlsoNegotiates is a shorthand for OnWhen(StatusIsVariantAlsoNegotiates, responseHandler).
// It matches status code 506.
func (r *Req) OnVariantAlsoNegotiates(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsVariantAlsoNegotiates, responseHandler)
}

// OnInsufficientStorage is a shorthand for OnWhen(StatusIsInsufficientStorage, responseHandler).
// It matches status code 507.
func (r *Req) OnInsufficientStorage(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsInsufficientStorage, responseHandler)
}

// OnLoopDetected is a shorthand for OnWhen(StatusIsLoopDetected, responseHandler).
// It matches status code 508.
func (r *Req) OnLoopDetected(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsLoopDetected, responseHandler)
}

// OnNotExtended is a shorthand for OnWhen(StatusIsNotExtended, responseHandler).
// It matches status code 510.
func (r *Req) OnNotExtended(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNotExtended, responseHandler)
}

// OnNetworkAuthenticationRequired is a shorthand for OnWhen(StatusIsNetworkAuthenticationRequired, responseHandler).
// It matches status code 511.
func (r *Req) OnNetworkAuthenticationRequired(responseHandler ResponseHandler) *Req {
	return r.OnWhen(StatusIsNetworkAuthenticationRequired, responseHandler)
}
