# Easy to use HTTP client in Go

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
## Check the status code and unmarshall the body

```go
type StarWarsCharacter struct {
	Name      string    `json:"name"`
	Height    string    `json:"height"`
	Mass      string    `json:"mass"`
	HairColor string    `json:"hair_color"`
	SkinColor string    `json:"skin_color"`
	EyeColor  string    `json:"eye_color"`
	BirthYear string    `json:"birth_year"`
	Gender    string    `json:"gender"`
	Homeworld string    `json:"homeworld"`
	Films     []string  `json:"films"`
	Species   []string  `json:"species"`
	Vehicles  []string  `json:"vehicles"`
	Starships []string  `json:"starships"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	Url       string    `json:"url"`
}
```
```go
if response.Status() == http.StatusOK {
    lukeSkywalker := StarWarsCharacter{}
    if err := response.UnmarshalJson(&lukeSkywalker); err != nil {
        log.Fatal(err)
    }
}
```
## All methods of request 
### HTTP Client Configuration
    UseHttpClient(client *http.Client) *Req
### Headers
    Header(key, val string) *Req
    ContentTypeJson() *Req
    ContentTypeText() *Req
    ContentTypeHtml() *Req
    ContentType(contentType string) *Req
    AcceptJson() *Req
###  Authentication
    AuthBasic(username, password string) *Req
    AuthToken(token string) *Req
### Query Parameters (Primitive Types)
    QueryInt8(name string, v int8) *Req
    QueryInt16(name string, v int16) *Req
    QueryInt32(name string, v int32) *Req
    QueryInt(name string, v int) *Req
    QueryInt64(name string, v int64) *Req
    QueryUint8(name string, v uint8) *Req
    QueryUint16(name string, v uint16) *Req
    QueryUint32(name string, v uint32) *Req
    QueryUint(name string, v uint) *Req
    QueryUint64(name string, v uint64) *Req
    QueryFloat32(name string, v float32) *Req
    QueryFloat64(name string, v float64) *Req
    QueryBool(name string, v bool) *Req
    QueryString(name string, v string) *Req
###  Query Parameters (Pointer Types)
    QueryInt8Ptr(name string, v *int8) *Req
    QueryInt16Ptr(name string, v *int16) *Req
    QueryInt32Ptr(name string, v *int32) *Req
    QueryIntPtr(name string, v *int) *Req
    QueryInt64Ptr(name string, v *int64) *Req
    QueryUint8Ptr(name string, v *uint8) *Req
    QueryUint16Ptr(name string, v *uint16) *Req
    QueryUint32Ptr(name string, v *uint32) *Req
    QueryUintPtr(name string, v *uint) *Req
    QueryUint64Ptr(name string, v *uint64) *Req
    QueryFloat32Ptr(name string, v *float32) *Req
    QueryFloat64Ptr(name string, v *float64) *Req
    QueryBoolPtr(name string, v *bool) *Req
    QueryStringPtr(name string, v *string) *Req

## All methods of response 

### Status Code Checking
    IsSuccess() bool
    IsInformational() bool
    IsRedirection() bool
    IsClientError() bool
    IsServerError() bool
    Is(statusCode int) bool
    IsOneOf(statusCodes ...int) bool
### Status Code Access
    Status() int
###  HTTP Response Access
    HttpResponse() *http.Response
###  Response Body Processing
    UnmarshalJson(t any) error