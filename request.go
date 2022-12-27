package gan

import (
	"net/http"
	"net/url"
)

type Request struct {
	r          *http.Request
	query      url.Values
	formParsed bool

	Method string
	Path   string
	Params map[string]string
}

func NewRequest(r *http.Request) *Request {
	req := &Request{
		r:      r,
		Method: r.Method,
		Path:   r.URL.Path,
	}
	return req
}

func (req *Request) GetFormValue(name string) string {
	if !req.formParsed {
		req.r.ParseForm()
		req.formParsed = true
	}
	return req.r.FormValue(name)
}

func (req *Request) GetQueryValue(name string) string {
	if req.query == nil {
		req.query = req.r.URL.Query()
	}
	return req.query.Get(name)
}

func (req *Request) GetCookieValue(name string) (string, error) {
	var value string
	cookie, err := req.r.Cookie(name)
	if err == nil {
		value = cookie.Value
	}
	return value, err
}

func (req *Request) QueryString() string {
	return req.r.URL.RawQuery
}
