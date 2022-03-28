package gee

import (
	"fmt"
	"net/http"
)

//定义我们框架中需要使用的handler函数类型
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine 实现ServeHTTP接口
type Engine struct {
	router map[string]HandlerFunc
}

// New 用于新建一个gee.Engine
func New() *Engine {
	return &Engine{
		make(map[string]HandlerFunc),
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}

// addRoute 将路径和handler添加到engine中
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
