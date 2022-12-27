package gan

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	w          http.ResponseWriter
	headers    map[string]string
	body       interface{}
	candicates map[string]any
	headerSent bool
	statusCode int
}

var engine = NewEngine()

func NewResponse(w http.ResponseWriter) *Response {
	res := &Response{
		w:          w,
		headers:    make(map[string]string),
		candicates: make(map[string]any),
		statusCode: http.StatusOK,
	}
	return res
}

func (res *Response) WithStatus(statusCode int) *Response {
	res.statusCode = statusCode
	return res
}

func (res *Response) WithHeader(key, value string) *Response {
	res.headers[key] = value
	return res
}

func (res *Response) WithBody(body interface{}) *Response {
	res.body = body
	return res
}

func (res *Response) WithCookie(name string, value string) *Response {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: 86400,
	}
	http.SetCookie(res.w, cookie)
	return res
}

func (res *Response) Redirect(url string, code int) {
	res.WithStatus(code)
	res.WithHeader("Location", url)
	res.Html("Redirect to " + url)
}

func (res *Response) WriteHeaders() *Response {
	for key, value := range res.headers {
		res.w.Header().Set(key, value)
	}
	res.w.WriteHeader(res.statusCode)
	res.headerSent = true
	return res
}

func (res *Response) End() {
	switch v := res.body.(type) {
	case []byte:
		res.Write(v)
	case string:
		res.Write([]byte(v))
	}
}

func (res *Response) Write(data []byte) *Response {
	if !res.headerSent {
		res.WriteHeaders()
	}
	res.w.Write(data)
	return res
}

func (res *Response) Text(content string) *Response {
	res.Write([]byte(content))
	return res
}

func (res *Response) Json(data any) *Response {
	res.WithHeader("content-type", "application/json")
	res.WriteHeaders()
	encoder := json.NewEncoder(res.w)
	encoder.Encode(data)
	return res
}

func (res *Response) Html(html string) *Response {
	res.WithHeader("content-type", "text/html; charset=utf-8")
	res.Write([]byte(html))
	return res
}

func (res *Response) SetRenderData(name string, value any) *Response {
	engine.Set(name, value)
	return res
}

func (res *Response) Render(name string, data map[string]any) *Response {
	res.WithHeader("content-type", "text/html; charset=utf-8")
	res.WriteHeaders()
	engine.Render(res.w, name, data)
	return res
}

func (res *Response) RenderError(err error) *Response {
	return res.Render("error", map[string]any{
		"error": err,
	})
}
