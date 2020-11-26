package webpack

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, url string, handler HandlerFunc) {
	parts := parsePattern(url)

	key := method + "-" + url
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(url, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, url string) (*node, map[string]string) {
	searchParts := parsePattern(url)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(context *Context) {
	n, params := r.getRoute(context.Method, context.Path)
	if n != nil {
		context.Params = params
		key := context.Method + "-" + n.pattern
		context.handlers = append(context.handlers, r.handlers[key])
	} else {
		context.handlers = append(context.handlers, func(c *Context) {
			c.FormatResponse(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	context.Next()
}
