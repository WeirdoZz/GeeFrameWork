package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 是字符串-空接口映射的别名
type H map[string]interface{}

type Context struct {
	// 基本部件
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	Params map[string]string
	// 返回信息
	StatusCode int

	// 中间件
	handlers []HandlerFunc
	index    int

	engine *Engine
}

// newContext 根绝输入的w和req返回一个Context的指针
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 用于控制执行多个中间件的
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

// PostForm 解析表单中的数据
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 解析URL中指定的参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 将状态码添加到context中，并且写入到返回头
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置返回头中的键值对
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回是string类型的数据
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 返回编码好的json数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 直接返回数据流
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 返回html数据，本质上也是将html转成[]byte数据流再传回去
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
