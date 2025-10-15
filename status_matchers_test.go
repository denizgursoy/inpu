package inpu

import (
	"net/http"
)

func (c *ClientSuite) Test_Response_IsSuccess() {
	c.Require().False(StatusIsSuccess.Match(http.StatusEarlyHints))
	c.Require().True(StatusIsSuccess.Match(http.StatusOK))
	c.Require().True(StatusIsSuccess.Match(http.StatusCreated))
	c.Require().False(StatusIsSuccess.Match(http.StatusMultipleChoices))
}

func (c *ClientSuite) Test_Response_IsServerError() {
	c.Require().False(StatusIsServerError.Match(http.StatusUnavailableForLegalReasons))
	c.Require().True(StatusIsServerError.Match(http.StatusInternalServerError))
}

func (c *ClientSuite) Test_Response_IsClientError() {
	c.Require().False(StatusIsClientError.Match(http.StatusPermanentRedirect))
	c.Require().True(StatusIsClientError.Match(http.StatusBadRequest))
	c.Require().True(StatusIsClientError.Match(http.StatusUnavailableForLegalReasons))
	c.Require().False(StatusIsClientError.Match(http.StatusInternalServerError))
}

func (c *ClientSuite) Test_Response_IsRedirection() {
	c.Require().True(StatusIsRedirection.Match(http.StatusMultipleChoices))
	c.Require().False(StatusIsRedirection.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_Response_IsInformational() {
	c.Require().True(StatusIsInformational.Match(http.StatusContinue))
	c.Require().False(StatusIsInformational.Match(http.StatusOK))
}

// 1xx Informational Tests
func (c *ClientSuite) Test_StatusIsContinue() {
	c.Require().True(StatusIsContinue.Match(http.StatusContinue))
	c.Require().False(StatusIsContinue.Match(http.StatusOK))
	c.Require().False(StatusIsContinue.Match(http.StatusSwitchingProtocols))
}

func (c *ClientSuite) Test_StatusIsSwitchingProtocols() {
	c.Require().True(StatusIsSwitchingProtocols.Match(http.StatusSwitchingProtocols))
	c.Require().False(StatusIsSwitchingProtocols.Match(http.StatusContinue))
	c.Require().False(StatusIsSwitchingProtocols.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsProcessing() {
	c.Require().True(StatusIsProcessing.Match(http.StatusProcessing))
	c.Require().False(StatusIsProcessing.Match(http.StatusContinue))
	c.Require().False(StatusIsProcessing.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsEarlyHints() {
	c.Require().True(StatusIsEarlyHints.Match(http.StatusEarlyHints))
	c.Require().False(StatusIsEarlyHints.Match(http.StatusContinue))
	c.Require().False(StatusIsEarlyHints.Match(http.StatusOK))
}

// 2xx Success Tests

func (c *ClientSuite) Test_StatusIsOk() {
	c.Require().True(StatusIsOk.Match(http.StatusOK))
	c.Require().False(StatusIsOk.Match(http.StatusCreated))
	c.Require().False(StatusIsOk.Match(http.StatusNotFound))
}

func (c *ClientSuite) Test_StatusIsCreated() {
	c.Require().True(StatusIsCreated.Match(http.StatusCreated))
	c.Require().False(StatusIsCreated.Match(http.StatusOK))
	c.Require().False(StatusIsCreated.Match(http.StatusAccepted))
}

func (c *ClientSuite) Test_StatusIsAccepted() {
	c.Require().True(StatusIsAccepted.Match(http.StatusAccepted))
	c.Require().False(StatusIsAccepted.Match(http.StatusOK))
	c.Require().False(StatusIsAccepted.Match(http.StatusCreated))
}

func (c *ClientSuite) Test_StatusIsNonAuthoritativeInfo() {
	c.Require().True(StatusIsNonAuthoritativeInfo.Match(http.StatusNonAuthoritativeInfo))
	c.Require().False(StatusIsNonAuthoritativeInfo.Match(http.StatusOK))
	c.Require().False(StatusIsNonAuthoritativeInfo.Match(http.StatusAccepted))
}

func (c *ClientSuite) Test_StatusIsNoContent() {
	c.Require().True(StatusIsNoContent.Match(http.StatusNoContent))
	c.Require().False(StatusIsNoContent.Match(http.StatusOK))
	c.Require().False(StatusIsNoContent.Match(http.StatusResetContent))
}

func (c *ClientSuite) Test_StatusIsResetContent() {
	c.Require().True(StatusIsResetContent.Match(http.StatusResetContent))
	c.Require().False(StatusIsResetContent.Match(http.StatusNoContent))
	c.Require().False(StatusIsResetContent.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsPartialContent() {
	c.Require().True(StatusIsPartialContent.Match(http.StatusPartialContent))
	c.Require().False(StatusIsPartialContent.Match(http.StatusOK))
	c.Require().False(StatusIsPartialContent.Match(http.StatusNoContent))
}

func (c *ClientSuite) Test_StatusIsMultiStatus() {
	c.Require().True(StatusIsMultiStatus.Match(http.StatusMultiStatus))
	c.Require().False(StatusIsMultiStatus.Match(http.StatusOK))
	c.Require().False(StatusIsMultiStatus.Match(http.StatusPartialContent))
}

func (c *ClientSuite) Test_StatusIsAlreadyReported() {
	c.Require().True(StatusIsAlreadyReported.Match(http.StatusAlreadyReported))
	c.Require().False(StatusIsAlreadyReported.Match(http.StatusMultiStatus))
	c.Require().False(StatusIsAlreadyReported.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsIMUsed() {
	c.Require().True(StatusIsIMUsed.Match(http.StatusIMUsed))
	c.Require().False(StatusIsIMUsed.Match(http.StatusOK))
	c.Require().False(StatusIsIMUsed.Match(http.StatusAlreadyReported))
}

// 3xx Redirection Tests

func (c *ClientSuite) Test_StatusIsMultipleChoices() {
	c.Require().True(StatusIsMultipleChoices.Match(http.StatusMultipleChoices))
	c.Require().False(StatusIsMultipleChoices.Match(http.StatusOK))
	c.Require().False(StatusIsMultipleChoices.Match(http.StatusMovedPermanently))
}

func (c *ClientSuite) Test_StatusIsMovedPermanently() {
	c.Require().True(StatusIsMovedPermanently.Match(http.StatusMovedPermanently))
	c.Require().False(StatusIsMovedPermanently.Match(http.StatusOK))
	c.Require().False(StatusIsMovedPermanently.Match(http.StatusFound))
}

func (c *ClientSuite) Test_StatusIsFound() {
	c.Require().True(StatusIsFound.Match(http.StatusFound))
	c.Require().False(StatusIsFound.Match(http.StatusMovedPermanently))
	c.Require().False(StatusIsFound.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsSeeOther() {
	c.Require().True(StatusIsSeeOther.Match(http.StatusSeeOther))
	c.Require().False(StatusIsSeeOther.Match(http.StatusFound))
	c.Require().False(StatusIsSeeOther.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsNotModified() {
	c.Require().True(StatusIsNotModified.Match(http.StatusNotModified))
	c.Require().False(StatusIsNotModified.Match(http.StatusSeeOther))
	c.Require().False(StatusIsNotModified.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsUseProxy() {
	c.Require().True(StatusIsUseProxy.Match(http.StatusUseProxy))
	c.Require().False(StatusIsUseProxy.Match(http.StatusNotModified))
	c.Require().False(StatusIsUseProxy.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsTemporaryRedirect() {
	c.Require().True(StatusIsTemporaryRedirect.Match(http.StatusTemporaryRedirect))
	c.Require().False(StatusIsTemporaryRedirect.Match(http.StatusFound))
	c.Require().False(StatusIsTemporaryRedirect.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsPermanentRedirect() {
	c.Require().True(StatusIsPermanentRedirect.Match(http.StatusPermanentRedirect))
	c.Require().False(StatusIsPermanentRedirect.Match(http.StatusTemporaryRedirect))
	c.Require().False(StatusIsPermanentRedirect.Match(http.StatusMovedPermanently))
}

// 4xx Client Error Tests

func (c *ClientSuite) Test_StatusIsBadRequest() {
	c.Require().True(StatusIsBadRequest.Match(http.StatusBadRequest))
	c.Require().False(StatusIsBadRequest.Match(http.StatusOK))
	c.Require().False(StatusIsBadRequest.Match(http.StatusUnauthorized))
}

func (c *ClientSuite) Test_StatusIsUnauthorized() {
	c.Require().True(StatusIsUnauthorized.Match(http.StatusUnauthorized))
	c.Require().False(StatusIsUnauthorized.Match(http.StatusBadRequest))
	c.Require().False(StatusIsUnauthorized.Match(http.StatusForbidden))
}

func (c *ClientSuite) Test_StatusIsPaymentRequired() {
	c.Require().True(StatusIsPaymentRequired.Match(http.StatusPaymentRequired))
	c.Require().False(StatusIsPaymentRequired.Match(http.StatusUnauthorized))
	c.Require().False(StatusIsPaymentRequired.Match(http.StatusForbidden))
}

func (c *ClientSuite) Test_StatusIsForbidden() {
	c.Require().True(StatusIsForbidden.Match(http.StatusForbidden))
	c.Require().False(StatusIsForbidden.Match(http.StatusUnauthorized))
	c.Require().False(StatusIsForbidden.Match(http.StatusNotFound))
}

func (c *ClientSuite) Test_StatusIsNotFound() {
	c.Require().True(StatusIsNotFound.Match(http.StatusNotFound))
	c.Require().False(StatusIsNotFound.Match(http.StatusForbidden))
	c.Require().False(StatusIsNotFound.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsMethodNotAllowed() {
	c.Require().True(StatusIsMethodNotAllowed.Match(http.StatusMethodNotAllowed))
	c.Require().False(StatusIsMethodNotAllowed.Match(http.StatusNotFound))
	c.Require().False(StatusIsMethodNotAllowed.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsNotAcceptable() {
	c.Require().True(StatusIsNotAcceptable.Match(http.StatusNotAcceptable))
	c.Require().False(StatusIsNotAcceptable.Match(http.StatusMethodNotAllowed))
	c.Require().False(StatusIsNotAcceptable.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsProxyAuthRequired() {
	c.Require().True(StatusIsProxyAuthRequired.Match(http.StatusProxyAuthRequired))
	c.Require().False(StatusIsProxyAuthRequired.Match(http.StatusUnauthorized))
	c.Require().False(StatusIsProxyAuthRequired.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsRequestTimeout() {
	c.Require().True(StatusIsRequestTimeout.Match(http.StatusRequestTimeout))
	c.Require().False(StatusIsRequestTimeout.Match(http.StatusBadRequest))
	c.Require().False(StatusIsRequestTimeout.Match(http.StatusGatewayTimeout))
}

func (c *ClientSuite) Test_StatusIsConflict() {
	c.Require().True(StatusIsConflict.Match(http.StatusConflict))
	c.Require().False(StatusIsConflict.Match(http.StatusBadRequest))
	c.Require().False(StatusIsConflict.Match(http.StatusRequestTimeout))
}

func (c *ClientSuite) Test_StatusIsGone() {
	c.Require().True(StatusIsGone.Match(http.StatusGone))
	c.Require().False(StatusIsGone.Match(http.StatusNotFound))
	c.Require().False(StatusIsGone.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsLengthRequired() {
	c.Require().True(StatusIsLengthRequired.Match(http.StatusLengthRequired))
	c.Require().False(StatusIsLengthRequired.Match(http.StatusBadRequest))
	c.Require().False(StatusIsLengthRequired.Match(http.StatusGone))
}

func (c *ClientSuite) Test_StatusIsPreconditionFailed() {
	c.Require().True(StatusIsPreconditionFailed.Match(http.StatusPreconditionFailed))
	c.Require().False(StatusIsPreconditionFailed.Match(http.StatusBadRequest))
	c.Require().False(StatusIsPreconditionFailed.Match(http.StatusPreconditionRequired))
}

func (c *ClientSuite) Test_StatusIsRequestEntityTooLarge() {
	c.Require().True(StatusIsRequestEntityTooLarge.Match(http.StatusRequestEntityTooLarge))
	c.Require().False(StatusIsRequestEntityTooLarge.Match(http.StatusBadRequest))
	c.Require().False(StatusIsRequestEntityTooLarge.Match(http.StatusRequestHeaderFieldsTooLarge))
}

func (c *ClientSuite) Test_StatusIsRequestURITooLong() {
	c.Require().True(StatusIsRequestURITooLong.Match(http.StatusRequestURITooLong))
	c.Require().False(StatusIsRequestURITooLong.Match(http.StatusRequestEntityTooLarge))
	c.Require().False(StatusIsRequestURITooLong.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsUnsupportedMediaType() {
	c.Require().True(StatusIsUnsupportedMediaType.Match(http.StatusUnsupportedMediaType))
	c.Require().False(StatusIsUnsupportedMediaType.Match(http.StatusNotAcceptable))
	c.Require().False(StatusIsUnsupportedMediaType.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsRequestedRangeNotSatisfiable() {
	c.Require().True(StatusIsRequestedRangeNotSatisfiable.Match(http.StatusRequestedRangeNotSatisfiable))
	c.Require().False(StatusIsRequestedRangeNotSatisfiable.Match(http.StatusBadRequest))
	c.Require().False(StatusIsRequestedRangeNotSatisfiable.Match(http.StatusPartialContent))
}

func (c *ClientSuite) Test_StatusIsExpectationFailed() {
	c.Require().True(StatusIsExpectationFailed.Match(http.StatusExpectationFailed))
	c.Require().False(StatusIsExpectationFailed.Match(http.StatusBadRequest))
	c.Require().False(StatusIsExpectationFailed.Match(http.StatusPreconditionFailed))
}

func (c *ClientSuite) Test_StatusIsTeapot() {
	c.Require().True(StatusIsTeapot.Match(http.StatusTeapot))
	c.Require().False(StatusIsTeapot.Match(http.StatusBadRequest))
	c.Require().False(StatusIsTeapot.Match(http.StatusOK))
}

func (c *ClientSuite) Test_StatusIsMisdirectedRequest() {
	c.Require().True(StatusIsMisdirectedRequest.Match(http.StatusMisdirectedRequest))
	c.Require().False(StatusIsMisdirectedRequest.Match(http.StatusBadRequest))
	c.Require().False(StatusIsMisdirectedRequest.Match(http.StatusTeapot))
}

func (c *ClientSuite) Test_StatusIsUnprocessableEntity() {
	c.Require().True(StatusIsUnprocessableEntity.Match(http.StatusUnprocessableEntity))
	c.Require().False(StatusIsUnprocessableEntity.Match(http.StatusBadRequest))
	c.Require().False(StatusIsUnprocessableEntity.Match(http.StatusMisdirectedRequest))
}

func (c *ClientSuite) Test_StatusIsLocked() {
	c.Require().True(StatusIsLocked.Match(http.StatusLocked))
	c.Require().False(StatusIsLocked.Match(http.StatusUnprocessableEntity))
	c.Require().False(StatusIsLocked.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsFailedDependency() {
	c.Require().True(StatusIsFailedDependency.Match(http.StatusFailedDependency))
	c.Require().False(StatusIsFailedDependency.Match(http.StatusLocked))
	c.Require().False(StatusIsFailedDependency.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsTooEarly() {
	c.Require().True(StatusIsTooEarly.Match(http.StatusTooEarly))
	c.Require().False(StatusIsTooEarly.Match(http.StatusBadRequest))
	c.Require().False(StatusIsTooEarly.Match(http.StatusFailedDependency))
}

func (c *ClientSuite) Test_StatusIsUpgradeRequired() {
	c.Require().True(StatusIsUpgradeRequired.Match(http.StatusUpgradeRequired))
	c.Require().False(StatusIsUpgradeRequired.Match(http.StatusBadRequest))
	c.Require().False(StatusIsUpgradeRequired.Match(http.StatusTooEarly))
}

func (c *ClientSuite) Test_StatusIsPreconditionRequired() {
	c.Require().True(StatusIsPreconditionRequired.Match(http.StatusPreconditionRequired))
	c.Require().False(StatusIsPreconditionRequired.Match(http.StatusPreconditionFailed))
	c.Require().False(StatusIsPreconditionRequired.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsTooManyRequests() {
	c.Require().True(StatusIsTooManyRequests.Match(http.StatusTooManyRequests))
	c.Require().False(StatusIsTooManyRequests.Match(http.StatusBadRequest))
	c.Require().False(StatusIsTooManyRequests.Match(http.StatusServiceUnavailable))
}

func (c *ClientSuite) Test_StatusIsRequestHeaderFieldsTooLarge() {
	c.Require().True(StatusIsRequestHeaderFieldsTooLarge.Match(http.StatusRequestHeaderFieldsTooLarge))
	c.Require().False(StatusIsRequestHeaderFieldsTooLarge.Match(http.StatusRequestEntityTooLarge))
	c.Require().False(StatusIsRequestHeaderFieldsTooLarge.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsUnavailableForLegalReasons() {
	c.Require().True(StatusIsUnavailableForLegalReasons.Match(http.StatusUnavailableForLegalReasons))
	c.Require().False(StatusIsUnavailableForLegalReasons.Match(http.StatusForbidden))
	c.Require().False(StatusIsUnavailableForLegalReasons.Match(http.StatusServiceUnavailable))
}

// 5xx Server Error Tests

func (c *ClientSuite) Test_StatusIsInternalServerError() {
	c.Require().True(StatusIsInternalServerError.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsInternalServerError.Match(http.StatusOK))
	c.Require().False(StatusIsInternalServerError.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsNotImplemented() {
	c.Require().True(StatusIsNotImplemented.Match(http.StatusNotImplemented))
	c.Require().False(StatusIsNotImplemented.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsNotImplemented.Match(http.StatusMethodNotAllowed))
}

func (c *ClientSuite) Test_StatusIsBadGateway() {
	c.Require().True(StatusIsBadGateway.Match(http.StatusBadGateway))
	c.Require().False(StatusIsBadGateway.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsBadGateway.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsServiceUnavailable() {
	c.Require().True(StatusIsServiceUnavailable.Match(http.StatusServiceUnavailable))
	c.Require().False(StatusIsServiceUnavailable.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsServiceUnavailable.Match(http.StatusBadGateway))
}

func (c *ClientSuite) Test_StatusIsGatewayTimeout() {
	c.Require().True(StatusIsGatewayTimeout.Match(http.StatusGatewayTimeout))
	c.Require().False(StatusIsGatewayTimeout.Match(http.StatusRequestTimeout))
	c.Require().False(StatusIsGatewayTimeout.Match(http.StatusInternalServerError))
}

func (c *ClientSuite) Test_StatusIsHTTPVersionNotSupported() {
	c.Require().True(StatusIsHTTPVersionNotSupported.Match(http.StatusHTTPVersionNotSupported))
	c.Require().False(StatusIsHTTPVersionNotSupported.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsHTTPVersionNotSupported.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_StatusIsVariantAlsoNegotiates() {
	c.Require().True(StatusIsVariantAlsoNegotiates.Match(http.StatusVariantAlsoNegotiates))
	c.Require().False(StatusIsVariantAlsoNegotiates.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsVariantAlsoNegotiates.Match(http.StatusNotAcceptable))
}

func (c *ClientSuite) Test_StatusIsInsufficientStorage() {
	c.Require().True(StatusIsInsufficientStorage.Match(http.StatusInsufficientStorage))
	c.Require().False(StatusIsInsufficientStorage.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsInsufficientStorage.Match(http.StatusRequestEntityTooLarge))
}

func (c *ClientSuite) Test_StatusIsLoopDetected() {
	c.Require().True(StatusIsLoopDetected.Match(http.StatusLoopDetected))
	c.Require().False(StatusIsLoopDetected.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsLoopDetected.Match(http.StatusInsufficientStorage))
}

func (c *ClientSuite) Test_StatusIsNotExtended() {
	c.Require().True(StatusIsNotExtended.Match(http.StatusNotExtended))
	c.Require().False(StatusIsNotExtended.Match(http.StatusInternalServerError))
	c.Require().False(StatusIsNotExtended.Match(http.StatusNotImplemented))
}

func (c *ClientSuite) Test_StatusIsNetworkAuthenticationRequired() {
	c.Require().True(StatusIsNetworkAuthenticationRequired.Match(http.StatusNetworkAuthenticationRequired))
	c.Require().False(StatusIsNetworkAuthenticationRequired.Match(http.StatusUnauthorized))
	c.Require().False(StatusIsNetworkAuthenticationRequired.Match(http.StatusProxyAuthRequired))
}
