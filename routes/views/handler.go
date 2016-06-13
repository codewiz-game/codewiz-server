package views

import (
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	sessionName  = "codewiz-session"
)

type handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *Router)

type handler struct {
	Router      *Router
	HandlerFunc handlerFunc
}

func newHandler(handlerFunc handlerFunc, router *Router) *handler {
	return &handler{Router: router, HandlerFunc: handlerFunc}
}

// This method performs all of the common code and passes down the
// frequently used components to the handlers.
func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := handler.Router.sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.HandlerFunc(w, r, session, handler.Router)
}