package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots   map[string]*trie
	handler map[string]HandlerFunc
}

func newRouter() *router {
	return &router{roots: make(map[string]*trie), handler: make(map[string]HandlerFunc)}
}

func parsePattern(pattern string) []string {
	sp := strings.Split(pattern, "/")

	parts := []string{}
	for _, v := range sp {
		if v != "" {
			parts = append(parts, v)
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &trie{part: "/"}
	}
	// 将路由添加到trie树上
	r.roots[method].insert(pattern, parts, 0)
	r.handler[key] = handler
}

func (r *router) getRoute(method string, path string) (*trie, map[string]string) {
	parts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	tire := root.search(parts, 0)

	if tire == nil {
		return nil, nil
	}

	param := map[string]string{}
	for index, part := range parsePattern(tire.pattern) {
		if part[0] == ':' {
			param[part[1:]] = parts[index]
		}
		if part[0] == '*' {
			param[part[1:]] = strings.Join(parts[index:], "/")
			break
		}
	}

	return tire, param
}

func (r *router) handle(c *Context) {
	trie, params := r.getRoute(c.Method, c.Path)

	if trie == nil || trie.pattern == "" {
		c.handler = append(c.handler, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	key := c.Method + "-" + trie.pattern
	c.Params = params
	if handle, ok := r.handler[key]; ok {
		c.handler = append(c.handler, handle)
	} else {
		c.handler = append(c.handler, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
