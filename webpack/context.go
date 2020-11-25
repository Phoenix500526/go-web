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
}

// Context Factory
func NewContext(response http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Response: response,
		Request:  request,
		Path:     request.URL.Path,
		Method:   request.Method,
	}
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

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	c.Response.Write([]byte(html))
}
