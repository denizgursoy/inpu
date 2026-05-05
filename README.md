# Inpu — Easy to use HTTP client in Go

`Inpu` is a Go HTTP client with a builder pattern that reads like natural sentences. It includes typed query parameters,
status-based response handling, pluggable middleware, structured logging, and OpenTelemetry observability.

```
go get github.com/denizgursoy/inpu
```

## Table of Contents

- [Quick Start](#quick-start)
- [Client](#client)
- [Requests](#requests)
- [Request Bodies](#request-bodies)
- [Response Handling](#response-handling)
- [Middlewares](#middlewares)
- [Logging](#logging)
- [Errors](#errors)
- [Utilities](#utilities)

## Quick Start

```go
var filteredTodos []Todo

err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
    QueryBool("completed", true).
    QueryInt("userId", 2).
    OnOk(inpu.ThenUnmarshalJsonTo(&filteredTodos)).
    OnAny(inpu.ThenReturnDefaultError).
    Send()
```

This sends:

```
GET https://jsonplaceholder.typicode.com/todos?completed=true&userId=2
```

If the response is `200`, the body is unmarshalled into `filteredTodos`. Any other status returns an error like
`called [GET] -> https://jsonplaceholder.typicode.com/todos?completed=true&userId=2 and got 500`.

## Client

### Creating a Client

```go
client := inpu.New().
    BasePath("https://api.example.com").
    TimeOutIn(5 * time.Second).
    AuthToken("eyJhbGciOiJSUzI1NiJ9").
    ContentTypeJson().
    AcceptJson()
```

`New()` returns a `*Client` with a pooled HTTP transport. All configuration methods return `*Client` for chaining.

`NewWithContext(ctx)` works the same but uses the provided context as the parent. When that context
is cancelled, all in-flight requests from the client are automatically cancelled:

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()

client := inpu.NewWithContext(ctx)
```

### Client Configuration

| Method | Description |
|---|---|
| `BasePath(url)` | Prepends base URL to every request path |
| `TimeOutIn(d)` | Sets timeout for all requests from this client |
| `Header(key, val)` | Adds a header to all requests |
| `AuthBasic(user, pass)` | Sets Basic auth header |
| `AuthToken(token)` | Sets Bearer token header |
| `UserAgent(ua)` | Sets User-Agent header |
| `ContentType(ct)` | Sets Content-Type header |
| `ContentTypeJson()` | Sets Content-Type to `application/json` |
| `ContentTypeXml()` | Sets Content-Type to `application/xml` |
| `ContentTypeText()` | Sets Content-Type to `text/plain` |
| `ContentTypeHtml()` | Sets Content-Type to `text/html` |
| `ContentTypeFormUrlEncoded()` | Sets Content-Type to `application/x-www-form-urlencoded` |
| `AcceptJson()` | Sets Accept to `application/json` |
| `TlsConfig(cfg)` | Sets custom `*tls.Config` |
| `DisableTLSVerification()` | Skips TLS certificate verification |
| `DisableHTTP2()` | Forces HTTP/1.1 only |
| `DisableRedirects()` | Disables automatic redirect following |
| `FollowRedirects(n)` | Follows up to `n` redirects |
| `EnableCookies()` | Enables a cookie jar for the client |
| `Close()` | Cancels pending requests and closes idle connections |
| `ToStandardClient()` | Returns the underlying `*http.Client` |

### Reusable Client

```go
client := inpu.New().
    BasePath("https://jsonplaceholder.typicode.com").
    Use(inpu.RetryMiddleware(2)).
    QueryInt("foo", 1).
    QueryString("foo1", "bar1").
    Header("foo", "bar").
    AuthToken("eyJhbGciOiJSUzI1NiJ9")

err := client.Get("/todos/1").Send()
err = client.Patch("/todos/1", inpu.BodyJson(payload)).Send()
err = client.Post("/todos", inpu.BodyJson(payload)).Send()
err = client.Put("/todos/1", inpu.BodyJson(payload)).Send()
```

The client above prepends the base path, adds query parameters `foo=1&foo1=bar1`, sets header `Foo: bar`,
and includes the Bearer token on every request.

## Requests

### HTTP Methods

Package-level functions (use a default shared client):

```go
inpu.Get(url)
inpu.Post(url, body)
inpu.Put(url, body)
inpu.Patch(url, body)
inpu.Delete(url, body)
inpu.Head(url)
```

Client methods (use the configured client):

```go
client.Get(url)
client.Post(url, body)
client.Put(url, body)
client.Patch(url, body)
client.Delete(url, body)
client.Head(url)
```

### Context-Aware Methods

Every HTTP method has a `*Ctx` variant that accepts a `context.Context` as the first argument:

```go
inpu.GetCtx(ctx, url)
client.GetCtx(ctx, url)
client.PostCtx(ctx, url, body)
// ... and so on for Put, Patch, Delete, Head
```

### Request Headers and Auth

These are available on both `*Client` and `*Req`:

```go
req := inpu.Get("https://api.example.com/items").
    Header("X-Custom", "value").
    AuthToken("eyJhbGciOiJSUzI1NiJ9").
    AuthBasic("user", "pass").
    ContentTypeJson().
    AcceptJson().
    UserAgent("my-app/1.0")
```

### Request Timeout

```go
req := inpu.Get("https://api.example.com/items").
    TimeOutIn(3 * time.Second)
```

This sets a per-request timeout independent of the client timeout.

### Query Parameters

Typed query parameter methods are available on both `*Client` and `*Req`:

```go
inpu.Get("https://example.com/search").
    QueryString("q", "golang").
    QueryInt("page", 1).
    QueryBool("active", true).
    QueryFloat64("score", 3.14).
    Send()
```

Full list of types: `QueryString`, `QueryInt`, `QueryInt8`, `QueryInt16`, `QueryInt32`, `QueryInt64`,
`QueryUint`, `QueryUint8`, `QueryUint16`, `QueryUint32`, `QueryUint64`, `QueryFloat32`, `QueryFloat64`, `QueryBool`.

Every type also has a `*Ptr` variant (e.g. `QueryIntPtr`, `QueryStringPtr`) that accepts a pointer
and is a no-op when the pointer is `nil`. This is useful for optional filter parameters.

## Request Bodies

Request body must implement the `Requester` interface. Use the built-in constructors:

```go
inpu.BodyJson(body any)                       // marshals to JSON
inpu.BodyXml(body any)                        // marshals to XML
inpu.BodyString(body string)                  // plain string
inpu.BodyReader(body io.Reader)               // raw reader
inpu.BodyFormData(body map[string][]string)   // URL-encoded form data
inpu.BodyFormDataFromMap(body map[string]string) // simplified form data
```

## Response Handling

### Shorthand Methods

Every status matcher has a corresponding shorthand method on the request. These read like natural sentences:

```go
OnOk(inpu.ThenUnmarshalJsonTo(&result)).          // on 200, unmarshal JSON
OnCreated(inpu.ThenDoNothing).                     // on 201, do nothing
OnNotFound(inpu.ThenReturnError(ErrItemNotFound)). // on 404, return custom error
OnUnauthorized(inpu.ThenReturnError(ErrNoAuth)).   // on 401, return custom error
OnSuccess(inpu.ThenUnmarshalJsonTo(&result)).      // on any 2xx, unmarshal JSON
OnClientError(inpu.ThenReturnDefaultError).        // on any 4xx, return default error
OnServerError(inpu.ThenReturnDefaultError).        // on any 5xx, return default error
OnAny(inpu.ThenReturnDefaultError).                // fallback for any status
OnAnyExcept(http.StatusOK, inpu.ThenReturnDefaultError). // any status except 200
```

Parameterized shorthands:

```go
On(http.StatusOK, inpu.ThenUnmarshalJsonTo(&result)).                                  // match a single status code
OnOneOf(inpu.ThenDoNothing, http.StatusOK, http.StatusCreated, http.StatusAccepted).   // match any of several codes
OnAnyExceptOneOf(inpu.ThenReturnDefaultError, http.StatusOK, http.StatusCreated).      // match any except several codes
```

Available shorthand methods for individual status codes:

| 1xx Informational | 2xx Success | 3xx Redirection | 4xx Client Error | 5xx Server Error |
|---|---|---|---|---|
| OnContinue (100) | OnOk (200) | OnMultipleChoices (300) | OnBadRequest (400) | OnInternalServerError (500) |
| OnSwitchingProtocols (101) | OnCreated (201) | OnMovedPermanently (301) | OnUnauthorized (401) | OnNotImplemented (501) |
| OnProcessing (102) | OnAccepted (202) | OnFound (302) | OnPaymentRequired (402) | OnBadGateway (502) |
| OnEarlyHints (103) | OnNonAuthoritativeInfo (203) | OnSeeOther (303) | OnForbidden (403) | OnServiceUnavailable (503) |
| | OnNoContent (204) | OnNotModified (304) | OnNotFound (404) | OnGatewayTimeout (504) |
| | OnResetContent (205) | OnUseProxy (305) | OnMethodNotAllowed (405) | OnHTTPVersionNotSupported (505) |
| | OnPartialContent (206) | OnTemporaryRedirect (307) | OnNotAcceptable (406) | OnVariantAlsoNegotiates (506) |
| | OnMultiStatus (207) | OnPermanentRedirect (308) | OnProxyAuthRequired (407) | OnInsufficientStorage (507) |
| | OnAlreadyReported (208) | | OnRequestTimeout (408) | OnLoopDetected (508) |
| | OnIMUsed (226) | | OnConflict (409) | OnNotExtended (510) |
| | | | OnGone (410) | OnNetworkAuthenticationRequired (511) |
| | | | OnLengthRequired (411) | |
| | | | OnPreconditionFailed (412) | |
| | | | OnRequestEntityTooLarge (413) | |
| | | | OnRequestURITooLong (414) | |
| | | | OnUnsupportedMediaType (415) | |
| | | | OnRequestedRangeNotSatisfiable (416) | |
| | | | OnExpectationFailed (417) | |
| | | | OnTeapot (418) | |
| | | | OnMisdirectedRequest (421) | |
| | | | OnUnprocessableEntity (422) | |
| | | | OnLocked (423) | |
| | | | OnFailedDependency (424) | |
| | | | OnTooEarly (425) | |
| | | | OnUpgradeRequired (426) | |
| | | | OnPreconditionRequired (428) | |
| | | | OnTooManyRequests (429) | |
| | | | OnRequestHeaderFieldsTooLarge (431) | |
| | | | OnUnavailableForLegalReasons (451) | |

Category-level shorthands: `OnSuccess`, `OnInformational`, `OnRedirection`, `OnClientError`, `OnServerError`

Wildcard shorthands: `OnAny`, `OnAnyExcept`, `OnAnyExceptOneOf`

### On (Advanced Usage)

`On` provides full control over status matching. Use it when combining matchers with `Not()` or for other advanced patterns:

```go
err := inpu.Get("https://api.example.com/items").
    On(inpu.StatusIsSuccess, inpu.ThenUnmarshalJsonTo(&items)).
    On(inpu.Not(inpu.StatusIsSuccess), inpu.ThenReturnDefaultError).
    Send()
```

Status matchers have priorities. When multiple matchers match, the one with the lowest priority value wins.
Priorities: `StatusIs` (1), `StatusIsOneOf` (2), category matchers like `StatusIsSuccess` (3),
`StatusAnyExcept` (8), `StatusAnyExceptOneOf` (9), `StatusAny` (10).

Available status matchers:

```go
StatusAny                              // matches any status code
StatusAnyExcept(statusCode int)        // matches any except the one provided
StatusAnyExceptOneOf(statusCodes ...int) // matches any except those provided
StatusIsSuccess                        // matches 2xx
StatusIsInformational                  // matches 1xx
StatusIsRedirection                    // matches 3xx
StatusIsClientError                    // matches 4xx
StatusIsServerError                    // matches 5xx
StatusIsOneOf(statusCodes ...int)      // matches any in the provided list
StatusIs(expectedStatus int)           // matches a specific status code
// Plus individual status matchers: StatusIsOk, StatusIsCreated, StatusIsNotFound, etc.
```

### Response Handlers

```go
ThenUnmarshalJsonTo(target any)                        // unmarshals the response body JSON into the pointer provided
ThenUnmarshalJsonAndReturnError(target any, err error) // unmarshals JSON and returns the provided error
ThenReturnError(err error)                             // returns the provided error
ThenReturnDefaultError                                 // returns an error with method, URL, and status code
ThenDoNothing                                          // returns nil (placeholder)
```

Note: `ThenReturnDefaultError` and `ThenDoNothing` are `ResponseHandler` values, not factories.
Pass them without parentheses. `ThenUnmarshalJsonTo`, `ThenUnmarshalJsonAndReturnError`, and `ThenReturnError`
are factories that return a `ResponseHandler`.

### Custom Handlers

You can pass any `func(r *http.Response) error` as a handler:

```go
err := inpu.Get("https://api.example.com/items").
    OnOk(func(r *http.Response) error {
        // custom processing
        return nil
    }).
    Send()
```

## Middlewares

Add middlewares to a client with `Use()`. Middlewares are sorted by priority (lower = closer to transport).

```go
client := inpu.New().
    Use(inpu.RetryMiddleware(2)).
    Use(inpu.NewLoggingMiddleware(inpu.WithVerbose())).
    Use(inpu.RequestIDMiddleware())
```

### Built-in Middlewares

| Middleware | Priority | Description |
|---|---|---|
| `NewLoggingMiddleware(opts...)` | 1 | Logs requests/responses. Masks sensitive headers. |
| `RequestIDMiddleware()` | 100 | Adds `X-Request-ID` header and stores ID in context |
| `ErrorHandlerMiddleware(handler)` | 50 | Calls handler on connection errors |
| `RetryMiddleware(maxRetries)` | 25 | Retries on server errors and 429 with exponential backoff |

### Logging Middleware Options

```go
inpu.NewLoggingMiddleware(
    inpu.WithVerbose(),            // log headers and bodies
    inpu.WithMaxBodyLogSize(8192), // truncate bodies larger than 8KB (default: 4KB)
)
```

| Option | Description |
|---|---|
| `WithVerbose()` | Log request/response headers and bodies |
| `WithDisabled()` | Create middleware in disabled state (no-op passthrough) |
| `WithMaxBodyLogSize(n)` | Max bytes to log for bodies; larger bodies are truncated (default: 4096) |

In verbose mode, bodies exceeding the max size are logged as:
`{"key":"value"...} ... (truncated, 52431 bytes total)`

### Retry Configuration

```go
client := inpu.New().Use(inpu.RetryMiddlewareWithConfig(inpu.RetryConfig{
    MaxRetries:        3,
    InitialBackoff:    500 * time.Millisecond,
    MaxBackoff:        30 * time.Second,
    BackoffMultiplier: 2.0,
    CustomRetryChecker: func(resp *http.Response, err error) bool {
        // custom retry logic
        return false
    },
}))
```

The retry middleware respects the `Retry-After` header on 429 and 503 responses. It retries on server errors
(5xx, except 501/505/508/506/511) and 429. TLS certificate errors are never retried.

### Custom Middleware

Implement the `Middleware` interface:

```go
type Middleware interface {
    ID() string
    Priority() int
    Apply(next http.RoundTripper) http.RoundTripper
}
```

Or use the helper constructors:

```go
// Modify outgoing requests
mw := inpu.RequestModifierMiddleware(
    func(req *http.Request) (*http.Request, error) {
        req.Header.Set("X-Custom", "value")
        return req, nil
    },
    "my-request-modifier", // unique ID
    50,                    // priority
)

// Modify incoming responses
mw := inpu.ResponseModifierMiddleware(
    func(resp *http.Response, err error) (*http.Response, error) {
        // inspect or modify the response
        return resp, err
    },
    "my-response-modifier",
    50,
)
```

### OpenTelemetry

```
go get github.com/denizgursoy/inpu/middlewares/otel
```

```go
import inpuotel "github.com/denizgursoy/inpu/middlewares/otel"

client := inpu.New().Use(inpuotel.NewMiddleware())
```

The OTel middleware (priority 2) sits inside the retry loop, so each attempt gets its own metrics and span.

**Options:**

| Option | Description |
|---|---|
| `WithMeterProvider(mp)` | Use a custom MeterProvider (default: global) |
| `WithTracerProvider(tp)` | Use a custom TracerProvider (default: global) |
| `WithPropagator(p)` | Use a custom TextMapPropagator (default: global) |
| `WithoutMetrics()` | Disable metric collection |
| `WithoutTracing()` | Disable tracing and context propagation |

**Collected metrics:**

| Metric | Type | Unit | Description |
|---|---|---|---|
| `http.client.request.duration` | Float64Histogram | s | Duration of each request attempt |
| `http.client.request.body.size` | Int64Histogram | By | Request body size |
| `http.client.response.body.size` | Int64Histogram | By | Response body size |
| `http.client.active_requests` | Int64UpDownCounter | {request} | Number of in-flight requests |
| `http.client.request.total` | Int64Counter | {request} | Total requests |
| `http.client.request.retry.count` | Int64Counter | {retry} | Total retries (attempt > 0) |

**Attributes:** `http.request.method`, `server.address`, `url.scheme`, `server.port`, `http.response.status_code`,
`http.resend_count`, `inpu.request.id` (when `RequestIDMiddleware` is used), `error.type` (on errors).

**Tracing:** Each request attempt creates a client span named `METHOD hostname`. Trace context is automatically
injected into outgoing request headers via the configured propagator. Spans are marked as error for 4xx/5xx responses.

### OAuth2 Client Credentials

```
go get github.com/denizgursoy/inpu/middlewares/oauth2
```

```go
import inpuoauth2 "github.com/denizgursoy/inpu/middlewares/oauth2"

client := inpu.New().Use(inpuoauth2.NewClientCredentialsMiddleware(
    clientcredentials.Config{
        ClientID:     "my-client-id",
        ClientSecret: "my-client-secret",
        TokenURL:     "https://auth.example.com/oauth/token",
    },
))
```

Automatically obtains and refreshes OAuth2 tokens using the client credentials flow (priority 75).

## Logging

By default, inpu is **silent** — no log output is produced. Logging is opt-in.

### Logger Interface

```go
type Logger interface {
    Error(ctx context.Context, err error, msg string, fields ...any)
    Warn(ctx context.Context, msg string, fields ...any)
    Info(ctx context.Context, msg string, fields ...any)
    Debug(ctx context.Context, msg string, fields ...any)
}
```

### Enabling Logging

There are three ways to enable logging:

**1. Set the global default logger:**

```go
// Use the built-in slog JSON logger
inpu.DefaultLogger = inpu.SlogLogger

// Or use a custom logger
inpu.DefaultLogger = myCustomLogger
```

This affects all requests across all clients. Middlewares like `RetryMiddleware` that log
internally will use this logger as a fallback when no context logger is set.

**2. Inject a logger into the context (per-request):**

```go
ctx := inpu.ContextWithLogger(ctx, logger)
err := inpu.GetCtx(ctx, "https://api.example.com/items").Send()
```

This takes precedence over `DefaultLogger` for that specific request.

**3. Add the LoggingMiddleware (logs HTTP traffic):**

```go
client := inpu.New().
    Use(inpu.LoggingMiddleware(false, false))  // (verbose, disabled)
```

The `LoggingMiddleware` logs request/response details (method, URL, status, duration).
With `verbose=true`, it also logs headers and bodies. Sensitive headers
(`Authorization`, `Cookie`, etc.) are automatically masked.

### Built-in Logger (slog)

A ready-to-use JSON logger is available as `inpu.SlogLogger`. You can also create one
with a specific level:

```go
logger := inpu.NewInpuLoggerFromSlog(inpu.LogLevelInfo) // LogLevelDebug, LogLevelWarn, LogLevelError
```

### Zap Adapter

```
go get github.com/denizgursoy/inpu/loggers/zap
```

```go
import inpuzap "github.com/denizgursoy/inpu/loggers/zap"

logger := inpuzap.NewInpuLoggerFromZapLogger(zapLogger)
```

### Zerolog Adapter

```
go get github.com/denizgursoy/inpu/loggers/zero
```

```go
import inpuzero "github.com/denizgursoy/inpu/loggers/zero"

logger := inpuzero.NewInpuLoggerFromZeroLog()
```

### Context Logger Injection

Retrieve the active logger with `inpu.ExtractLoggerFromContext(ctx)`. The logging middleware
and built-in logger automatically include `request_id` when `RequestIDMiddleware` is active.

## Errors

Sentinel errors returned by inpu:

| Error | Description |
|---|---|
| `ErrRequestCreationFailed` | Could not create the HTTP request |
| `ErrInvalidBody` | Could not create the request body |
| `ErrConnectionFailed` | Connection to the server failed |
| `ErrCouldNotParseBaseUrl` | Invalid base path URL |
| `ErrCouldNotParsePath` | Invalid request path |
| `ErrMarshalToNil` | Tried to unmarshal into nil |
| `ErrNotPointerParameter` | Tried to unmarshal into non-pointer type |

`DefaultError` is returned by `ThenReturnDefaultError` and formats as
`called [METHOD] -> URL and got STATUS_CODE`.

## Utilities

```go
// Extract request ID from context (set by RequestIDMiddleware)
inpu.ExtractRequestIDFromContext(ctx) // returns *string

// Extract retry attempt number from context (set by RetryMiddleware)
inpu.ExtractRetryAttemptFromContext(ctx) // returns int (0 for first attempt)

// Pre-configured HTTP clients and transports
inpu.DefaultPooledClient()    // reusable client with connection pooling
inpu.DefaultClient()          // client with no keepalives
inpu.DefaultPooledTransport() // transport with connection pooling
inpu.DefaultTransport()       // transport with no keepalives

// Auth header helpers
inpu.GetTokenHeaderValue(token)              // returns "Bearer <token>"
inpu.GetBasicAuthHeaderValue(user, password) // returns "Basic <base64>"

// Response body helper
inpu.DrainBodyAndClose(resp) // drains and closes the response body
```

The package also exports `Header*` constants (e.g. `HeaderContentType`, `HeaderAuthorization`, `HeaderXRequestID`)
and `MimeType*` constants (e.g. `MimeTypeJson`, `MimeTypeXml`, `MimeTypeFormUrlEncoded`).
