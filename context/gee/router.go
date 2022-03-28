package gee

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

// newRouter 创建一个新的router管理器
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// addRoute 为路径管理器中的路径映射表添加映射
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//这里每添加一个路径都会在服务器端显示
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle 调用handler函数，如果存在的话
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	}
}
