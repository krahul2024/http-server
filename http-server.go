package main

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type HttpServer struct {
	Port string
	Routers []*Router
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
    print ("[%s] %s %s\n", time.Now().Format(time.RFC3339), req.Method, req.Path)
	return true
}

func NewHttpServer(port string) *HttpServer {
    return &HttpServer{
        Port: port,
        Routers: []*Router{},
        StartTime: time.Now(),
        ReqLogger: DefaultLogger{},
    }
}

func (h *HttpServer) Listen() error {
	if len(h.Port) < 5 || !strings.HasPrefix(h.Port, ":") {
		return errors.New("invalid port configuration")
	}

	// call connect here with port
	return nil
}

func (h *HttpServer) AddRouter (r *Router) error {
	if r == nil {
		return errors.New("can't add nil router")
	}

	h.Routers = append(h.Routers, r)
	return nil
}
