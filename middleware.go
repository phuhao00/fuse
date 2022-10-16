package fuse

type middleware interface {
	Middleware(handler Handler) Handler
}

type MiddlewareFunc func(handler Handler) Handler

func (mw MiddlewareFunc) Middleware(handler Handler) Handler {
	return mw(handler)
}

func (r *Router) Use(mwf ...MiddlewareFunc) {
	for _, fn := range mwf {
		r.middlewares = append(r.middlewares, fn)
	}
}

func (r *Router) useInterface(mw middleware) {
	r.middlewares = append(r.middlewares, mw)
}
