package httptool

import (
	"net/http"
)

type middleware interface {
	mid_handle(http.Handler) http.Handler
}

type Router struct {
	middlewareChain []middleware
	mux map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]http.Handler),
	}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		middleware := r.middlewareChain[i]
		mergedHandler = middleware.mid_handle(mergedHandler)
	}
	r.mux[route] = mergedHandler
}

func (r *Router) BindMux(mux *http.ServeMux) {
	for pattern, handler := range r.mux {
		mux.Handle(pattern, handler)
	}
}
