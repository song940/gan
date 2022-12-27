package gan

import (
	"net/http"
)

type Gan struct {
	router     *Router
	middleware []MiddlewareFunc
}

type RouteFunc func(ctx *Context)
type NextFunc func()
type MiddlewareFunc func(ctx *Context, next NextFunc)

func New() *Gan {
	router := NewRouter()
	app := &Gan{
		router:     router,
		middleware: make([]MiddlewareFunc, 0),
	}
	return app
}

func (g *Gan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)
	var i = 0
	var next NextFunc
	next = func() {
		if i >= len(g.middleware) {
			return
		}
		// log.Println("mw:", i)
		mw := g.middleware[i]
		i = i + 1
		mw(ctx, next)
	}
	next()
	g.router.ServeHTTP(ctx)
}

func (g *Gan) Routes() []*Route {
	return g.router.routes
}

func (g *Gan) Use(mw MiddlewareFunc) *Gan {
	g.middleware = append(g.middleware, mw)
	return g
}

func (g *Gan) Route(method string, path string, handler RouteFunc) *Gan {
	g.router.Add(&Route{
		method:  method,
		path:    path,
		handler: handler,
	})
	return g
}

func (g *Gan) GET(path string, handler RouteFunc) *Gan {
	return g.Route(http.MethodGet, path, handler)
}

func (g *Gan) POST(path string, handler RouteFunc) *Gan {
	return g.Route(http.MethodPost, path, handler)
}

func (g *Gan) DELETE(path string, handler RouteFunc) *Gan {
	return g.Route(http.MethodDelete, path, handler)
}

func (g *Gan) Run(addr string) error {
	return http.ListenAndServe(addr, g)
}
