package webpack

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implements http.Handler interface
type Engine struct {
	// routing table
	router map[string]HandlerFunc
}

// Engine Factory
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (e *Engine) AddRoute(method, url string, handler HandlerFunc) {
	key := method + "-" + url
	log.Printf("Route: %04s - %s", method, url)
	e.router[key] = handler
}

func (e *Engine) GET(URL string, handler HandlerFunc) {
	e.AddRoute("GET", URL, handler)
}

func (e *Engine) POST(URL string, handler HandlerFunc) {
	e.AddRoute("POST", URL, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(respone http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(respone, request)
	} else {
		fmt.Fprint(respone, "404 Not Found\n")
	}
}
