package webpack

import (
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method, url string, handler HandlerFunc) {
	key := method + "-" + url
	r.handlers[key] = handler
}

func (r *router) handle(context *Context) {
	key := context.Method + "-" + context.Path
	if handler, ok := r.handlers[key]; ok {
		handler(context)
	} else {
		context.FormatResponse(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
	}
}
