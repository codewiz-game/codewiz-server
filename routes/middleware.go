package routes

import (
	"net/http"
)

// -----------------------
// - Middleware interface
// -----------------------
type Middleware interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *Context, HandlerFunc)
}

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc)

func (middleware MiddlewareFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
	middleware(w, r, context, next)
}

func wrapHandler(handler Handler) Middleware {
	return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
		handler.ServeHTTP(w,r,context)
	})
}

var voidMiddleware = MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
	next(w,r,context)
})

// -----------------------
// - Middleware Node
// -----------------------

type middlewareNode struct {
	middleware Middleware
	next       *middlewareNode
}

func newMiddlewareNode(middleware Middleware) *middlewareNode {
	return &middlewareNode{middleware: middleware}
}

func (node *middlewareNode) ServeHTTP(w http.ResponseWriter, r *http.Request, context *Context) {
	node.middleware.ServeHTTP(w, r, context, HandlerFunc(node.next.ServeHTTP))
}

// -----------------------
// - Middleware List
// -----------------------

type middlewareList struct {
	initialised bool
	root        *middlewareNode
	last        *middlewareNode
	router      *Router
}

func newMiddlewareList(router *Router) *middlewareList {
	tempRoot := newMiddlewareNode(voidMiddleware)
	return &middlewareList{root: tempRoot, last: tempRoot, router: router}
}

func (list *middlewareList) Add(middleware Middleware) {
	if list.initialised {
		node := newMiddlewareNode(middleware)
		list.last.next = node
		list.last = node
	} else {
		list.root.middleware = middleware
		list.initialised = true
	}
}

func (list *middlewareList) Clone() *middlewareList {
	clonedList := newMiddlewareList(list.router)
	for currentNode := list.root; currentNode != nil; {
		clonedList.Add(currentNode.middleware)
		currentNode = currentNode.next
	}
	return clonedList
}

func (list *middlewareList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := newContext(list.router)
	list.root.ServeHTTP(w, r, context)
}
