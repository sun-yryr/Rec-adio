package router

import "github.com/gin-gonic/gin"

type Router interface {
	Routes() []Route
}

type Route interface {
	Handler() gin.HandlerFunc
	Method() string
	Path() string
}

type localRouter struct {
	handler gin.HandlerFunc
	method  string
	path    string
}

func (r *localRouter) Handler() gin.HandlerFunc {
	return r.handler
}

func (r *localRouter) Method() string {
	return r.method
}

func (r *localRouter) Path() string {
	return r.path
}

func NewGetRouter(path string, handler gin.HandlerFunc) Route {
	return &localRouter{
		handler: handler,
		method:  "GET",
		path:    path,
	}
}

func NewPostRouter(path string, handler gin.HandlerFunc) Route {
	return &localRouter{
		handler: handler,
		method:  "POST",
		path:    path,
	}
}

func NewPutRouter(path string, handler gin.HandlerFunc) Route {
	return &localRouter{
		handler: handler,
		method:  "PUT",
		path:    path,
	}
}

func NewPatchRouter(path string, handler gin.HandlerFunc) Route {
	return &localRouter{
		handler: handler,
		method:  "PATCH",
		path:    path,
	}
}

func NewDeleteRouter(path string, handler gin.HandlerFunc) Route {
	return &localRouter{
		handler: handler,
		method:  "DELETE",
		path:    path,
	}
}
