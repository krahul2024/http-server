package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// TODO:
// Add queue, deadline, timeouts and connection pool
// middleware support
// proper response return to client

func main() {
	server := HttpServer {
		Port: ":4000",
		StartTime: time.Now(),
		ReqLogger: DefaultLogger{},
	}

	server.AddRouter(userRouter())

	log.Printf ("Starting the server at localhost%v\n", server.Port)
	if err := server.Listen(); err != nil {
		log.Fatalf ("Error listening starting server localhost%v\nError:%v\n", server.Port, err)
	}
}

func setQueryParams(lastSeg string, urlContent *UrlContent) error {
	queryMarkCount := strings.Count(lastSeg, "?")
	print("Query mark count = ", queryMarkCount)

	if queryMarkCount > 1 {
		*urlContent = UrlContent{}
		return errors.New("invalid query params and url formation")
	}

	parts := strings.SplitN(lastSeg, "?", 2)
	print("last-segment = ", lastSeg, ", parts = ", parts)

	if len(parts) == 2  && len(parts[1]) > 0{
		queryStr := parts[1]
		if urlContent.queryParams == nil {
			urlContent.queryParams = make(map[string]string)
		}
		for _, pair := range strings.FieldsFunc(queryStr, func(r rune) bool { return r == '&' }) {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				continue // might want to do something
			}

			key, _ := url.QueryUnescape(kv[0])
			value, _ := url.QueryUnescape(kv[1])
			urlContent.queryParams[key] = value
		}
	}

	return nil
}

func hasQueryParams(reqUrlStrSeq string, srcPrefix string) bool {
	return strings.TrimSpace(strings.SplitN(reqUrlStrSeq, "?", 2)[0]) == srcPrefix
}

func matchUrlStr(reqUrlParts []string, srcUrlParts []string, urlContent *UrlContent) (bool, error) {
	n := len(srcUrlParts)

	err := setQueryParams(reqUrlParts[n-1], urlContent)
	if err != nil {
		*urlContent = UrlContent{}
		return false, err
	}

	for i := range n {
		reqStr := reqUrlParts[i]
		srcStr := srcUrlParts[i]

		if reqStr== srcStr || (i == n-1 && hasQueryParams(reqStr, srcStr)){
			continue
		} else if strings.HasPrefix(srcStr, ":") {
			key := strings.TrimPrefix(srcStr, ":")
			if urlContent.pathParams == nil {
				urlContent.pathParams = make(map[string]string)
			}
			pathParamValue, _, _ := strings.Cut(reqUrlParts[i], "?")
			urlContent.pathParams[key] = pathParamValue
		} else {
			*urlContent = UrlContent{}
			return false, nil
		}
	}

	return true, nil
}

func parseUrl(url string, handlers []RouteHandler) (UrlContent, error) {
	reqUrlParts := strings.FieldsFunc(url, func(r rune) bool { return r == '/' })
	urlContent := UrlContent{}

	for _, h := range handlers {
		if len(reqUrlParts) == len(h.pathParts) && len(reqUrlParts) > 0 {
			didMatch, err := matchUrlStr(reqUrlParts, h.pathParts, &urlContent)
			if err != nil || didMatch {
				urlContent.handler = h.handler
				return urlContent, err
			}
		}
	}

	urlContent.reqStatus = StatusNotFound
	urlContent.handler = nil

	return urlContent, errors.New("invalid url, not found")
}

func handleReq(req *Request, res *Response, routers map[string]*Router) {
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

	printf("Path Parts(%v):%+v\n", len(parts), parts)

	if router, ok := routers[parts[0]]; ok {
		print("Found the router =", router)

		// now looking for the handler function
		restUrl := path[len(parts[0]):]
		print("Rest URL =", restUrl)

		urlCont, err := parseUrl(restUrl, router.handlers)
		printf("%+v", urlCont)

		if err == nil {
			req.Headers["QueryParams"] = urlCont.queryParams
			req.Headers["PathParams"] = urlCont.pathParams
			urlCont.handler(req, res)
		} else {
			res.StatusCode = urlCont.reqStatus.Code
			res.StatusText = urlCont.reqStatus.Msg
		}
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
