package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	path string
	mux  *mux.Router
	middleware *middlewareList
}

func NewRouter(path string) *Router {
	router := &Router{
		mux:  mux.NewRouter(),
		path: path,
	}
	router.middleware = newMiddlewareList(router)
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) Use(middleware Middleware) *Router {
	router.middleware.Add(middleware)
	return router
}

func (router *Router) StrictSlash(strictSlash bool) *Router {
	router.mux.StrictSlash(strictSlash)
	return router
}

func (router *Router) Path(path string) *Route {
	return newRoute(path, router)
}