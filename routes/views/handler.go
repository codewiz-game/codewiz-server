package views

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"net/http"
)


const (
	sessionName  = "codewiz-session"
)

type handler struct {
	router      *router
	handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *router)
}

func init() {
	// Register all of the types that will be stored as session values
	// to allow them to be encoded to the session cookie.
	gob.Register(map[string][]string{})
}

func createHandler(handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *router), router *router) *handler {
	return &handler{router: router, handlerFunc: handlerFunc}
}

// This method performs all of the common code and passes down the
// frequently used components to the handlers.
func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get or create a session cookie
	session, err := handler.router.sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass this down to the handler function along with a reference to the router
	handler.handlerFunc(w, r, session, handler.router)
}
