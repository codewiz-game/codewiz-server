package routes

import (
	"github.com/crob1140/codewiz/models/users"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	paths "path"
)

// ----------
// - Context
// ----------

type Context struct {
	User    *users.User
	Session *sessions.Session
	Router  *Router
}

func newContext(router *Router) *Context {
	return &Context{Router: router}
}

// --------------------
// - Handler interface
// --------------------

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *Context)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, *Context)

func (handler HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, context *Context) {
	handler(w, r, context)
}

var voidHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request, context *Context) {})

// --------
// - Route
// --------

type Route struct {
	path        string
	mux         *mux.Route
	router      *Router
	middleware  *middlewareList
	handlerNode *middlewareNode
}

func newRoute(path string, router *Router) *Route {
	routePath := paths.Join(router.path, path)
	mux := router.mux.Path(routePath)
	route := &Route{
		path:        path,
		mux:         mux,
		router:      router,
		middleware:  router.middleware.Clone(),
		handlerNode: newMiddlewareNode(wrapHandler(voidHandler)),
	}

	route.mux.Handler(route.middleware)
	route.middleware.last.next = route.handlerNode
	return route
}

func (route *Route) Methods(methods ...string) *Route {
	route.mux.Methods(methods...)
	return route
}

func (route *Route) HandlerFunc(handlerFunc HandlerFunc) *Route {
	// TODO: this could be simplified by making 'middleware' a wrapper that just doesn't call next,
	// so that we don't need to always create the void node that follows
	route.handlerNode.middleware = wrapHandler(handlerFunc)
	return route
}

func (route *Route) Subrouter() *Router {
	path := paths.Join(route.router.path, route.path)
	mux := route.mux.Subrouter()
	middleware := route.middleware.Clone()
	return &Router{
		path : path,
		mux : mux,
		middleware : middleware,
	}
}