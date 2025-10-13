package main

const CRLF    = "\r\n"
const NewLine = "\n"

var HeaderValTypes = map[string]HeaderType{
	"content-length"                   : TypeInt,
	"connection"                       : TypeString,
	"content-type"                     : TypeString,
	"accept"                           : TypeString,
	"host"                             : TypeString,
	"user-agent"                       : TypeString,
	"authorization"                    : TypeString,
	"accept-encoding"                  : TypeString,
	"cache-control"                    : TypeString,
	"upgrade"                          : TypeString,
	"origin"                           : TypeString,
	"access-control-request-method"    : TypeString,
	"access-control-request-headers"   : TypeString,
	"access-control-allow-origin"      : TypeString,
	"access-control-allow-methods"     : TypeString,
	"access-control-allow-headers"     : TypeString,
	"access-control-allow-credentials" : TypeString,
	"access-control-max-age"           : TypeInt,
}

// empty struct holds 0 bytes(damn)
var HopByHopHeaders = map[string]struct{}{
	"connection":          {},
	"keep-alive":          {},
	"proxy-authenticate":  {},
	"proxy-authorization": {},
	"te":                  {},
	"trailer":             {},
	"transfer-encoding":   {},
	"upgrade":             {},
}

var (
	StatusOK                  = ReqStatus{Code: 200, Msg: "OK"}
	StatusCreated             = ReqStatus{Code: 201, Msg: "Created"}
	StatusAccepted            = ReqStatus{Code: 202, Msg: "Accepted"}
	StatusNoContent           = ReqStatus{Code: 204, Msg: "No Content"}

	StatusBadRequest          = ReqStatus{Code: 400, Msg: "Bad Request"}
	StatusUnauthorized        = ReqStatus{Code: 401, Msg: "Unauthorized"}
	StatusForbidden           = ReqStatus{Code: 403, Msg: "Forbidden"}
	StatusNotFound            = ReqStatus{Code: 404, Msg: "Not Found"}
	StatusMethodNotAllowed    = ReqStatus{Code: 405, Msg: "Method Not Allowed"}

	StatusInternalServerError = ReqStatus{Code: 500, Msg: "Internal Server Error"}
	StatusNotImplemented      = ReqStatus{Code: 501, Msg: "Not Implemented"}
	StatusBadGateway          = ReqStatus{Code: 502, Msg: "Bad Gateway"}
	StatusServiceUnavailable  = ReqStatus{Code: 503, Msg: "Service Unavailable"}
)
