package main

import "strings"

type Router struct {
	pathPrefix string
	handlers []RouteHandler
}

type RouteHandler struct {
	route string
	handler HandlerFunc
	pathParts []string
}

type HandlerFunc func(req *Request, res *Response)
type HandleRoute func(path string, handler HandlerFunc)

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
		pathParts: strings.FieldsFunc(routeStr, func(r rune) bool { return r == '/'}),
	})
}

func registerUserRoute() {
	router := NewRouter("/user")

	router.Add("/add",       addUserHandler) // make this router.Method(GET/POST/PUT/DELETE/...)
	router.Add("/all",       allUserHandler)
	router.Add("/:id",       getUserById)
	router.Add("/:id/name", getUserByNameAndId)
	router.Add("/:id/:name", getUserByNameAndId)
}
