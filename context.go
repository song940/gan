package gan

import (
	"net/http"
)

type Context struct {
	request  *Request
	response *Response
	data     map[string]any
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	req := NewRequest(r)
	res := NewResponse(w)
	ctx := &Context{
		request:  req,
		response: res,
		data:     map[string]any{},
	}
	return ctx
}

func (ctx *Context) Request() *Request {
	return ctx.request
}

func (ctx *Context) Response() *Response {
	return ctx.response
}

func (ctx *Context) Set(name string, value any) *Context {
	ctx.data[name] = value
	return ctx
}

func (ctx *Context) Get(name string) any {
	return ctx.data[name]
}

func (ctx *Context) GetParam(name string) string {
	return ctx.Request().Params[name]
}
