package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pradist/apidemo/todo"
)

type MyContext struct {
	*gin.Context
}

func NewMyContext(c *gin.Context) *MyContext {
	return &MyContext{Context: c}
}

func (c MyContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}

func (c MyContext) JSON(statusCode int, v interface{}) {
	c.Context.JSON(statusCode, v)
}
func (c MyContext) TransactionID() string {
	return c.Context.Request.Header.Get("TransactionID")
}
func (c MyContext) Audience() string {
	if aud, ok := c.Context.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}
	return ""
}

func NewGinHandler(handler func(todo.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}

type MyRouter struct {
	*gin.Engine
}

func NewMyRouter() *MyRouter {
	r := gin.Default()
	return &MyRouter{r}
}

func (r *MyRouter) POST(path string, handler func(todo.Context)) {
	r.Engine.POST(path, NewGinHandler(handler))
}
