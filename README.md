# Easy to use HTTP client in Go
Inpu is a Go HTTP client.

To download:`go get github.com/denizgursoy/inpu`

## Build the request and send
```go
response, err := Get("https://swapi.dev/api/people/1").
    QueryInt("foo", 1).
    QueryString("foo1", "bar1").
    Header("foo", "bar").
    Header("foo1", "bar1").
    AuthToken("bar-password").
    Send()
```
Does the following call
```
https://swapi.dev/api/people/1?foo=1&foo1=bar1 
Authorization: Bearer bar-password
Foo: bar
Foo1: bar1

```
## Check the status code and unmarshall the body
```go
lukeSkywalker := StarWarsCharacter{}
if response.Status() == http.StatusOK {
    if err := response.UnmarshalJson(&lukeSkywalker); err != nil {
        log.Fatal(err)
    }
}
```

## Create clients
```go
client := New().
		UseMiddlewares(RetryMiddleware(2)).// add middlewares
		BasePath("https://swapi.dev/api"). // prepends base path to every call uri
		TimeOutIn(time.Second *5). // causes every request created from the client to expire in the duration
		// following are added to every request created form the client
		QueryInt("foo", 1).
		QueryString("foo1", "bar1").
		Header("foo", "bar").
		Header("foo1", "bar1").
		AuthToken("bar-password")

	response, err :=client.Get("/people/1").Send()
```
it creates the same get call
```
https://swapi.dev/api/people/1?foo=1&foo1=bar1 
Authorization: Bearer bar-password
Foo: bar
Foo1: bar1

```
client is reusable
```go
	response, err :=client.Get("/people/1").Send()
	response, err =client.Patch("/people/1", payload).Send()
	response, err =client.Post("/people", payload).Send()
	response, err =client.Put("/people/1", payload).Send()
```