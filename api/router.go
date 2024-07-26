package api

import (
	"fmt"
	"net/http"
	"regexp"
)

//
// Path Parameters: '/v1/user/{id}' accessed with r.PathValue("id")
//
// rg := NewRouterGroup("/v1")
// rg.Group("/user", AuthMiddleware)
// 	rg.GET("/{id}", GetUser)

var duplicateSlashes = regexp.MustCompile("/{2,}")

type RouterGroup struct {
	mux        *http.ServeMux
	root       *RouterGroup
	basePath   string
	middleware []Middleware
	groups     []*RouterGroup
}

func cleanPath(path string) string {
	if path == "" {
		path = "/"
	}

	if path[0] != '/' {
		path = "/" + path
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	return duplicateSlashes.ReplaceAllString(path, "/")
}

func NewRouterGroup(mux *http.ServeMux, path string, middleware ...Middleware) (*RouterGroup, error) {
	if mux == nil {
		return nil, fmt.Errorf("no mux provided")
	}

	path = cleanPath(path)

	if middleware == nil {
		middleware = []Middleware{}
	}

	return &RouterGroup{
		mux:        mux,
		basePath:   path,
		middleware: middleware,
		groups:     []*RouterGroup{},
	}, nil
}

func (group *RouterGroup) Group(path string, middleware ...Middleware) *RouterGroup {
	r := group.root
	if r == nil {
		r = group
	}

	path = cleanPath(path)
	return &RouterGroup{
		root:       r,
		basePath:   cleanPath(group.basePath + "/" + path),
		middleware: append(group.middleware, middleware...),
	}
}

func (group *RouterGroup) GET(path string, handler http.Handler, middleware ...Middleware) {
	h := group.genHandler(handler, middleware)
	path = cleanPath(path)

	mux := group.mux
	if mux == nil {
		mux = group.root.mux
	}
	mux.Handle("GET "+cleanPath(group.basePath+"/"+path), h)
}

func (group RouterGroup) genHandler(h http.Handler, middleware []Middleware) http.Handler {
	if len(middleware) != 0 {
		h = compile(h, middleware)
	}

	if len(group.middleware) != 0 {
		h = compile(h, group.middleware)
	}

	return h
}

func compile(h http.Handler, middleware []Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}
