package gan

type Route struct {
	method  string
	path    string
	handler RouteFunc
	repath  *PathToRegexp
}

type Router struct {
	routes []*Route
}

func Match(r *Request, route *Route) bool {
	if r.Method != route.method {
		return false
	}
	return route.repath.Test(r.Path)
}

func Find(routes []*Route, r *Request) (match *Route) {
	for _, route := range routes {
		if Match(r, route) {
			return route
		}
	}
	return
}

func NewRouter() (router *Router) {
	router = &Router{
		routes: make([]*Route, 0),
	}
	return
}

func (router *Router) Add(route *Route) *Router {
	route.repath = Compile(route.path)
	router.routes = append(router.routes, route)
	return router
}

func (router *Router) Match(r *Request) *Route {
	return Find(router.routes, r)
}

func (router *Router) ServeHTTP(ctx *Context) {
	req := ctx.Request()
	res := ctx.Response()
	route := router.Match(req)
	if route != nil && route.handler != nil {
		req.Params = route.repath.Parse(req.Path)
		route.handler(ctx)
	} else {
		res.WithStatus(404).Text("Not Found")
	}
}
