package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var AllRoutes = map[string]*Router{}

type Router struct {
	pathPrefix string
	handlers []RouteHandler
}

type RouteHandler struct {
	route string
	handler HandlerFunc
	pathParts []string
}

func NewRouter (pathPrefix string) *Router {
	router := Router {
		pathPrefix: pathPrefix,
		handlers: []RouteHandler{},
	}

	AllRoutes[pathPrefix] = &router
	return &router
}

func (r *Router) Add(routeStr string, handler HandlerFunc) {
	r.handlers = append(r.handlers, RouteHandler{
		route: routeStr,
		handler: handler,
		pathParts: strings.Split(routeStr, "/"),
	})
}

func main() {
	registerUserRoute()
	// fmt.Printf("%+v\n", AllRoutes)
	connect()
}

func addUserHandler(req *Request, res *Response) {
	fmt.Println("This is from addUserHandler")
	res.Body = []byte("Hello from add user")
}

func allUserHandler(req *Request, res *Response) {
	fmt.Println("This is from allUserHandler")
	res.Body = []byte("Hello from all user")
}

func registerUserRoute() {
	router := NewRouter("/user")
	router.Add("/add", addUserHandler)
	router.Add("/all", allUserHandler)
}

const CRLF    = "\r\n"
const NewLine = "\n"

type HeaderType int
const (
	TypeString HeaderType = iota
	TypeInt
	TypeBool
)

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

func connect() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("Error = %v\n", err.Error())
	}

	log.Printf ("Listening for connections on localhost%v", ":4000")

	connCounter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error = %v\n", err.Error())
			break
		}

		connCounter += 1
		go handleConn(conn, connCounter)
	}
}

type HandlerFunc func(req *Request, res *Response)
type HandleRoute func(path string, handler HandlerFunc)

func handleConn(conn net.Conn, connCounter int) {
	defer conn.Close()

	for {
		var (
			req Request
			res Response
		)

		if err := readMsg(conn, &req); err != nil {
			log.Printf ("Conn[%v]: Read Error = %v\n", connCounter, err.Error())
			return;
		}

		handleReq(&req, &res)

		res.Version    = req.Version
		res.Headers    = req.Headers
		res.Body       = req.Body

		if err := writeMsg(conn, &res); err != nil {
			log.Printf("Error = %v\n", err.Error())
			return
		}
	}
}

type UrlContent struct {
	url string
	queryParams map[string]interface{}
	pathParams map[string]interface{}
}

func matchUrlStr(reqUrlParts *[]string, srcUrlParts *[]string, urlContent *UrlContent) (bool, error) {
	// dynamic path values are denoted by :key
	didMatch := false

	for i :=0; i < len(*reqUrlParts); i++ {
		// matches perfectly, then okay
		if *urlStrParts[i] == utmParts[i] || strings.HasPrefix(utmParts[i], ":") {
			// here get the dynamic path key and value
		} else {
			return nil, nil
		}
	}

	// reset the urlContent if it didn't match
	if !didMatch {
		urlContent = &UrlContent{}
	}

	return didMatch, nil
}

func parseUrl(url string, handlers *[]RouteHandler, res *Response) (UrlContent, error) {
	// get the string upto ?
	reqUrlParts := strings.Split("?", url)
	// urlStr := strings.TrimSpace(urlStrParts[0])

	urlContent := UrlContent{}
	for _, h := range *handlers {
		if len(reqUrlParts) == len(h.pathParts) {
			didMatch, err := matchUrlStr(&reqUrlParts, &h.pathParts, &urlContent)
			if err != nil {
				return urlContent, err
			}

			if didMatch {
				// we are done here, nothing to do
			}
		}

		// otherwise we don't care
	}

	// match the obtained string with routerHandler urls
	// if matched then okay, move to query params
	// else 404

	return UrlContent{}, nil
}

func handleReq(req *Request, res *Response) {
	path := req.Path
	parts := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })

	// if the path was "/", handle it as home
	if len(parts) == 0 {
		parts = append(parts, "")
	}

	// prepending "/"
	for i, p := range parts {
		parts[i] = "/" + p
	}

	fmt.Printf("Path Parts(%v):%+v\n", len(parts), parts)

	if router, ok := AllRoutes[parts[0]]; ok {
		fmt.Println("Found the router =", router)

		// now looking for the handler function
		restUrl := path[len(parts[0]):]
		fmt.Println("Rest URL =", restUrl)

		urlCont, _ := parseUrl(restUrl, &router.handlers)
		req.Headers["QueryParams"] = urlCont.queryParams
		req.Headers["PathParams"] = urlCont.pathParams

		res.StatusCode = 200
		res.StatusText = "OK"
	} else {
		res.StatusCode = StatusNotFound.Code
		res.StatusText = StatusNotFound.Msg
	}
}

func readMsg (conn net.Conn, req *Request) (error) {
	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return  err
	}
	defer conn.SetReadDeadline(time.Time{})

	reader := bufio.NewReader(conn)
	reqStartLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	req.Headers = make(map[string]any)
	if err = parseReqLine(reqStartLine, req); err != nil {
		return err
	}

	if err = parseHeaders(reader, req); err != nil {
		return err
	}

	cl := req.Headers["content-length"]
	contentLength, ok := cl.(int)
	if !ok {
		contentLength = 0;
	}

	req.Body = make([]byte, contentLength)
	if err = readBody(reader, req); err != nil {
		return err
	}

	fmt.Printf("%+v\n", req)
	fmt.Println("Body = ", string(req.Body))

	return nil
}

func parseReqLine (s string, r *Request) error {
	s = strings.TrimSuffix(s, CRLF)
	strs := strings.Split(s, " ")
	if len(strs) != 3 {
		return errors.New("error parsing request")
	}

	// method validation
	switch strs[0] {
	case "GET": r.Method = GET; break;
	case "POST": r.Method = POST; break;
	case "PUT" : r.Method = PUT; break;
	case "DELETE": r.Method = DELETE; break;
	default: return errors.New("invalid request method")
	}

	r.Path = strs[1]

	// version validation
	if strs[2] != "HTTP/1.1" {
		return errors.New("invalid http version")
	}
	r.Version = strs[2]

	return nil
}

func parseHeaders(reader *bufio.Reader, req *Request) (error) {

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("error parsing headers")
		}

		s := strings.TrimSuffix(line, CRLF)
		if line == "" {
			break
		}

		ss := strings.SplitN(s, ":", 2)
		if len(ss) != 2 {
			break
		}

		key := strings.ToLower(strings.TrimSpace(ss[0]))
		value := strings.TrimSpace(ss[1])

		valueType, ok := HeaderValTypes[key]
		if !ok {
			valueType = TypeString
		}

		if v, err := parseHeaderValue(value, valueType); err == nil {
			req.Headers[key] = v
		} else {
			return err
		}
	}

	return nil
}

func parseHeaderValue(value string, typ HeaderType) (any, error) {
    switch typ {
    case TypeInt:
        return strToInt(value)
    case TypeBool:
        return strToBool(value)
    default:
        return value, nil
    }
}

func readBody(reader *bufio.Reader, req *Request) error {
	_, err := reader.Read(req.Body)
	return err
}

func strToInt(value string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(value))
}

func strToBool(value string) (bool, error) {
	value = strings.TrimSpace(strings.ToLower(value))

	switch value {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid value for bool parsing : %q", value)
	}
}

func writeMsg(conn net.Conn, res *Response) error {
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}
	defer conn.SetWriteDeadline(time.Time{})

	var buffer bytes.Buffer
	start, err := resStartStr(res)
	if err != nil {
		return err
	}
	buffer.WriteString(start)

	hStr, err := headersStr(res)
	if err != nil {
		return err
	}

	buffer.WriteString(hStr)
	buffer.WriteString(CRLF)
	buffer.Write(res.Body)

	writer := bufio.NewWriter(conn)
	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return writer.Flush()
}

func resStartStr (res *Response) (string, error) {
	//TODO: check validity of version, and rest of the stuff
	resStr := fmt.Sprintf("%v %v %v%v", res.Version, res.StatusCode, res.StatusText, CRLF)
	return resStr, nil
}

func headersStr (res *Response) (string, error) {
	headerStr := ""

	for key, val := range res.Headers {
		if _, ok := HopByHopHeaders[strings.ToLower(key)]; !ok {
			headerStr += fmt.Sprintf("%v: %v%v", key, val, CRLF)
		}
	}

	return headerStr, nil
}

//TODO: add a queue to handle connection limit
