package main

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"fmt"
	"time"
)

type HttpServer struct {
	Port string
	routers  map[string]*Router
	StartTime time.Time
	ReqCount int64
	NotFoundHandler HandlerFunc
	MethodNotAllowedHandler HandlerFunc
	reqCountMutex sync.Mutex
	ReqLogger
}

type ReqLogger interface {
	Log(req *Request, res *Response) bool
}

type DefaultLogger struct {}
func (l DefaultLogger) Log(req *Request, res *Response) bool {
    printf ("[%s] %s %s\n", time.Now().Format(time.RFC3339), req.Method, req.Path)
	return true
}

func NewHttpServer(port string) *HttpServer {
    return &HttpServer{
        Port: port,
        routers: map[string]*Router{},
        StartTime: time.Now(),
        ReqLogger: DefaultLogger{},
    }
}

func (h *HttpServer) Listen() error {
	if len(h.Port) < 5 || !strings.HasPrefix(h.Port, ":") {
		return errors.New("invalid port configuration")
	}

	listener, err := net.Listen("tcp", h.Port)
	if err != nil {
		return err
	}

	connCounter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error = %v\n", err.Error())
			break
		}

		connCounter += 1
		go handleConn(conn, connCounter, h.routers)
	}

	return nil
}

func (h *HttpServer) AddRouter (r *Router) error {
	if r == nil {
		return errors.New("can't add nil router")
	}

	if h.routers == nil {
		h.routers = map[string]*Router{}
	}

	if _, ok := h.routers[r.pathPrefix]; ok {
		return fmt.Errorf("path-prefix already has a router assigned, path = %v", r.pathPrefix)
	}

	h.routers[r.pathPrefix] = r;
	return nil
}

func handleConn(conn net.Conn, connCounter int, routers map[string]*Router) {
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

		handleReq(&req, &res, routers)

		res.Version    = req.Version
		res.Headers    = req.Headers
		// res.Body       = req.Body

		if err := writeMsg(conn, &res); err != nil {
			log.Printf("Error = %v\n", err.Error())
			return
		}
	}
}
