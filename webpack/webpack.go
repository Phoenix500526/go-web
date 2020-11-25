package webpack

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

// Engine implements http.Handler interface
type Engine struct {
	// routing table
	router *router
}

// Engine Factory
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) addRoute(method, url string, handler HandlerFunc) {
	log.Printf("Route: %04s - %s", method, url)
	e.router.addRoute(method, url, handler)
}

func (e *Engine) GET(URL string, handler HandlerFunc) {
	e.addRoute("GET", URL, handler)
}

func (e *Engine) POST(URL string, handler HandlerFunc) {
	e.addRoute("POST", URL, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(respone http.ResponseWriter, request *http.Request) {
	context := NewContext(respone, request)
	e.router.handle(context)
}
