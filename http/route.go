package http

type Route struct {
	path           string
	name           string
	beforeHandlers []HandlerFunc
	afterHandlers  []HandlerFunc
	handler        HandlerFunc
}

func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}

func (r *Route) Before(handlers ...HandlerFunc) *Route {
	r.beforeHandlers = append(r.beforeHandlers, handlers...)
	return r
}
func (r *Route) After(handlers ...HandlerFunc) *Route {
	r.afterHandlers = append(r.afterHandlers, handlers...)
	return r
}

func (r *Route) Handlers() []HandlerFunc {
	return append(append(r.beforeHandlers, r.handler), r.afterHandlers...)
}

func newRoute(path string, handler HandlerFunc) *Route {
	return &Route{
		path:           path,
		handler:        handler,
		beforeHandlers: make([]HandlerFunc, 0),
		afterHandlers:  make([]HandlerFunc, 0),
	}
}
