# Easy to use HTTP client in Go
`Inpu` is a Go HTTP client that allows developers to create request with builder pattern. It also 
provides some utility methods and common constants.

To download:`go get github.com/denizgursoy/inpu`

## Build the request and send
```go
	err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
        QueryBool("completed", true).
        QueryInt("userId", 2).
        OnReply(inpu.StatusIsSuccess, inpu.UnmarshalJson(&filteredTodos)).
        OnReply(inpu.StatusIs(http.StatusNotFound), inpu.ReturnError(errors.New("could not find any item"))).
        OnReply(inpu.StatusIs(http.StatusInternalServerError), inpu.ReturnError(errors.New("server could not handle the request"))).
        OnReply(inpu.StatusAny, inpu.ReturnError(errors.New("could not fetch the todo items"))).
        Send()
```
Does the following call
```
https://jsonplaceholder.typicode.com/todos?completed=1&userId=bar1 
```
## Check the status code and unmarshall the body
`OnReply` method allows developers to perform certain operation matched by some status conditions
```go
OnReply(inpu.StatusIsSuccess, inpu.UnmarshalJson(&filteredTodos)). // it marshals the body to the array 
OnReply(inpu.StatusAny, inpu.ReturnError(errors.New("could not fetch the todo items"))). // it returns the error if status does not match any condition
```
Other status matchers are:
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
```
Other response handler
```go
UnmarshalJson(t any) // it marshals the response body into the 
ReturnError(err error) // it returns the error provided
ReturnDefaultError() // it returns default error that prints status code, url and method
```
You can also add custom handler.
```go
err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
    QueryBool("completed", true).
    QueryInt("userId", 2).
    OnReply(inpu.StatusAny, func(r *http.Response) error {
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
		AuthToken("bar-password")

        err :=client.Get("/todos/1").Send()
```
It creates the same get call
```
https://jsonplaceholder.typicode.com/todos/1?completed=1&userId=bar1&foo=1&foo1=bar1 
Authorization: Bearer bar-password
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
Request body can be `io.Reader` or any value. If no marshaler found, JSON marshaler is used by default. You can also use
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