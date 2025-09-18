package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	connect()
}

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

type RequestMethod string
const (
	GET RequestMethod = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"
)

type RequestMeta struct {
	Method RequestMethod
	Path string
	Version string
	Headers map[string]any
}

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

func handleConn(conn net.Conn, connCounter int) {
	defer conn.Close()

	for {
		if err := readMsg(conn); err != nil {
			log.Printf ("Conn[%v]: Read Error = %v\n", connCounter, err.Error())
			return;
		}

		writeMsg        := "Status: MSG_READ"
		bytesWrite, err := conn.Write([]byte(writeMsg))
		if err != nil {
			log.Printf("Error = %v\n", err.Error())
			return
		}

		fmt.Println (bytesWrite, "bytes Written to connection ", connCounter)
	}
}

func readMsg (conn net.Conn) (error) {
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

	var req RequestMeta
	req.Headers = make(map[string]any)
	if err = parseReqLine(reqStartLine, &req); err != nil {
		return err
	}

	if err = parseHeaders(reader, &req); err != nil {
		return err
	}

	fmt.Printf("%+v\n", req)

	// body parsing left

	return nil
}

func parseReqLine (s string, r *RequestMeta) error {
	s = strings.TrimSuffix(s, "\r\n")
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

func parseHeaders(reader *bufio.Reader, req *RequestMeta) (error) {

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("error parsing headers")
		}

		s := strings.TrimSuffix(line, "\r\n")
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

//TODO: add a queue to handle connection limit
