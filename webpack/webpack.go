package webpack

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

// Engine implements http.Handler interface
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		engine      *Engine
	}

	Engine struct {
		// routing table
		*RouterGroup
		router *router
		groups []*RouterGroup // store all groups
	}
)

// Engine Factory
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route: %04s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(URL string, handler HandlerFunc) {
	g.addRoute("GET", URL, handler)
}

func (g *RouterGroup) POST(URL string, handler HandlerFunc) {
	g.addRoute("POST", URL, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (e *Engine) ServeHTTP(respone http.ResponseWriter, request *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context := newContext(respone, request)
	context.handlers = middlewares
	e.router.handle(context)
}
