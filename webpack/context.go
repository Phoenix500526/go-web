package webpack

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	// for Requestuest
	Path   string
	Method string
	// for Responseonse
	StatusCode int
	// for router
	Params map[string]string

	// for middleware
	handlers []HandlerFunc
	index    int

	// Engine pointer
	engine *Engine
}

// Context Factory
func newContext(response http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Response: response,
		Request:  request,
		Path:     request.URL.Path,
		Method:   request.Method,
		index:    -1,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Response.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Response.Header().Set(key, value)
}

func (c *Context) FormatResponse(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Response.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Response)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Response, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	c.Response.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Response, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
