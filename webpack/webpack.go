package webpack

import (
	"log"
	"net/http"
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
		gs     []*RouterGroup // store all gs
	}
)

// Engine Factory
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.gs = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		engine: engine,
	}
	engine.gs = append(engine.gs, newGroup)
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

func (e *Engine) ServeHTTP(respone http.ResponseWriter, request *http.Request) {
	context := NewContext(respone, request)
	e.router.handle(context)
}
