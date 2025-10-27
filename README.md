# Easy to use HTTP client in Go
`Inpu` is a Go HTTP client that allows developers to create request with builder pattern. It also 
provides some utility methods and common constants.

To download:`go get github.com/denizgursoy/inpu`

## Build the request and send
```go
err :=  inpu.Get("https://jsonplaceholder.typicode.com/todos").
            QueryBool("completed", true).
            QueryInt("userId", 2).
            OnReplyIf(inpu.StatusIsOk, inpu.ThenUnmarshalJsonTo(&filteredTodos)).
	        OnReplyIf(inpu.StatusAny, inpu.ThenReturnDefaultError).
            Send()
```
Does the following call
```
https://jsonplaceholder.typicode.com/todos?completed=true&userId=2
```
It will marshal the response body to `filteredTodos` if status code `200`. If response code anything except `200`, It will
return `called [GET] -> https://jsonplaceholder.typicode.com/todos?completed=true&userId=2 and got 500` error to provide more information.

## Check the status code and unmarshall the body
`OnReplyIf` method allows developers to execute `type ResponseHandler func(r *http.Response) error` operation matched by `StatusMatcher`
```go
OnReplyIf(inpu.StatusIsSuccess, inpu.UnmarshalJson(&filteredTodos)). // it marshals the body to the variable
OnReplyIf(inpu.StatusAny, inpu.ReturnError(errors.New("could not fetch the todo items"))). // it returns the error if status does not match any condition
```
Available status matchers are:
```go
StatusAny // it matches any status code
StatusAnyExcept(statusCode int) // it matches any status code expect the one provided
StatusAnyExceptOneOf(statusCodes ...int) // it matches any status code expect those provided
StatusIsSuccess // it matches any status between [200, 300)
StatusIsInformational // it matches any status between [100, 200)
StatusIsRedirection // it matches any status between [300, 400)
StatusIsClientError // it matches any status between [400, 500)
StatusIsServerError // it matches any status >= 500
StatusIsOneOf(statusCodes ...int) // it matches any status code in those provided
StatusIs(expectedStatus int) // it checks if it matches the status provided 
// 1xx Informational
StatusIsContinue                          // it matches status 100
StatusIsSwitchingProtocols                // it matches status 101
StatusIsProcessing                        // it matches status 102
StatusIsEarlyHints                        // it matches status 103
// 2xx Success
StatusIsOk                                // it matches status 200
StatusIsCreated                           // it matches status 201
StatusIsAccepted                          // it matches status 202
StatusIsNonAuthoritativeInfo              // it matches status 203
StatusIsNoContent                         // it matches status 204
StatusIsResetContent                      // it matches status 205
StatusIsPartialContent                    // it matches status 206
StatusIsMultiStatus                       // it matches status 207
StatusIsAlreadyReported                   // it matches status 208
StatusIsIMUsed                            // it matches status 226
// 3xx Redirection
StatusIsMultipleChoices                   // it matches status 300
StatusIsMovedPermanently                  // it matches status 301
StatusIsFound                             // it matches status 302
StatusIsSeeOther                          // it matches status 303
StatusIsNotModified                       // it matches status 304
StatusIsUseProxy                          // it matches status 305
StatusIsTemporaryRedirect                 // it matches status 307
StatusIsPermanentRedirect                 // it matches status 308
// 4xx Client Errors
StatusIsBadRequest                        // it matches status 400
StatusIsUnauthorized                      // it matches status 401
StatusIsPaymentRequired                   // it matches status 402
StatusIsForbidden                         // it matches status 403
StatusIsNotFound                          // it matches status 404
StatusIsMethodNotAllowed                  // it matches status 405
StatusIsNotAcceptable                     // it matches status 406
StatusIsProxyAuthRequired                 // it matches status 407
StatusIsRequestTimeout                    // it matches status 408
StatusIsConflict                          // it matches status 409
StatusIsGone                              // it matches status 410
StatusIsLengthRequired                    // it matches status 411
StatusIsPreconditionFailed                // it matches status 412
StatusIsRequestEntityTooLarge             // it matches status 413
StatusIsRequestURITooLong                 // it matches status 414
StatusIsUnsupportedMediaType              // it matches status 415
StatusIsRequestedRangeNotSatisfiable      // it matches status 416
StatusIsExpectationFailed                 // it matches status 417
StatusIsTeapot                            // it matches status 418
StatusIsMisdirectedRequest                // it matches status 421
StatusIsUnprocessableEntity               // it matches status 422
StatusIsLocked                            // it matches status 423
StatusIsFailedDependency                  // it matches status 424
StatusIsTooEarly                          // it matches status 425
StatusIsUpgradeRequired                   // it matches status 426
StatusIsPreconditionRequired              // it matches status 428
StatusIsTooManyRequests                   // it matches status 429
StatusIsRequestHeaderFieldsTooLarge       // it matches status 431
StatusIsUnavailableForLegalReasons        // it matches status 451
// 5xx Server Errors
StatusIsInternalServerError               // it matches status 500
StatusIsNotImplemented                    // it matches status 501
StatusIsBadGateway                        // it matches status 502
StatusIsServiceUnavailable                // it matches status 503
StatusIsGatewayTimeout                    // it matches status 504
StatusIsHTTPVersionNotSupported           // it matches status 505
StatusIsVariantAlsoNegotiates             // it matches status 506
StatusIsInsufficientStorage               // it matches status 507
StatusIsLoopDetected                      // it matches status 508
StatusIsNotExtended                       // it matches status 510
StatusIsNetworkAuthenticationRequired     // it matches status 511
```
Available response handlers are:
```go
ThenUnmarshalJsonTo(t any) // it marshals the response body into the 
ThenReturnError(err error) // it returns the error provided
ThenReturnDefaultError() // it returns default error that prints status code, url and method
ThenDoNothing() // just a place holder
```
You can also add custom handler:
```go
err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
    QueryBool("completed", true).
    QueryInt("userId", 2).
    OnReplyIf(inpu.StatusAny, func(r *http.Response) error {
        // custom processing
        return nil
    }).
    Send()
```

## Create clients
```go
client := New().
		UseMiddlewares(RetryMiddleware(2)).// add middlewares
		BasePath("https://jsonplaceholder.typicode.com"). // prepends base path to every call uri
		TimeOutIn(time.Second *5). // causes every request created from the client to expire in the duration
		// following are added to every request created form the client
		QueryInt("foo", 1).
		QueryString("foo1", "bar1").
		Header("foo", "bar").
		Header("foo1", "bar1").
		AuthToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ")

        err :=client.Get("/todos/1").Send()
```
It creates the same get call
```
https://jsonplaceholder.typicode.com/todos/1?completed=1&userId=bar1&foo=1&foo1=bar1 
Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ
Foo: bar
Foo1: bar1
```
Client is reusable
```go
err := client.Get("/todos/1").Send()
err = client.Patch("/todos/1", BodyJson(payload)).Send()
err = client.Post("/todos",  BodyJson(payload)).Send()
err = client.Put("/todos/1",  BodyJson(payload)).Send()
```

## Request Bodies
Request body must be `inpu.Requester`. You can use
following functions to create request body in the specific formats.

```go
BodyString(body string)
BodyXml(body any)
BodyJson(body any)
BodyReader(body io.Reader)
BodyFormDataFromMap(body map[string]string)
BodyFormData(body map[string][]string)
```

## Middlewares

```go
RetryMiddleware(2) // retries twice in case of certain codes
LoggingMiddleware(true,false) // logs the request and responses
RequestIDMiddleware() // add request ID  header to all request
ErrorHandlerMiddleware(handler) // calls the handler in case of connection error
```