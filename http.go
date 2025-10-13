package main

type RequestMethod string
const (
	GET RequestMethod = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"
)

type Request struct {
	Method RequestMethod
	Path string
	Version string
	Headers map[string]any
	Body    []byte
}

type Response struct {
	Version    string
	StatusCode int
	StatusText string
	Headers    map[string]any
	Body []byte
}

type ReqStatus struct {
	Code int
	Msg  string
}

type HeaderType int
const (
	TypeString HeaderType = iota
	TypeInt
	TypeBool
)

type UrlContent struct {
	url string
	reqStatus ReqStatus
	handler HandlerFunc
	queryParams map[string]string
	pathParams map[string]string
}

