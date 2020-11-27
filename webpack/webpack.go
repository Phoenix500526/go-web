package webpack

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
		router        *router
		groups        []*RouterGroup     // store all groups
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   //for html render
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
	log.Printf("Route: %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(URL string, handler HandlerFunc) {
	g.addRoute("GET", URL, handler)
}

func (g *RouterGroup) POST(URL string, handler HandlerFunc) {
	g.addRoute("POST", URL, handler)
}

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Response, c.Request)
	}
}

func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	g.GET(urlPattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
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
	context.engine = e
	e.router.handle(context)
}
