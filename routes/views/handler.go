package views

import (
	"github.com/crob1140/codewiz/models/users"
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	sessionName  = "codewiz-session"
)

type context struct {
	Router *Router
	User *users.User
	Session *sessions.Session
}

type handlerFunc func(http.ResponseWriter, *http.Request, *context)

type handler struct {
	Router *Router
	RequiresLogin bool
	HandlerFunc handlerFunc
}

func newContext(router *Router, user *users.User, session *sessions.Session) *context {
	return &context{User : user, Session : session, Router : router}
}

func newHandler(handlerFunc handlerFunc, router *Router, requiresLogin bool) *handler {
	return &handler{Router: router, HandlerFunc: handlerFunc, RequiresLogin : requiresLogin}
}

// This method performs all of the common code and passes down the
// frequently used components to the handlers.
func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	router := handler.Router

	session, err := router.sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUserForSession(session, router.userDao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil && handler.RequiresLogin {
		loginUrl := router.Login()
		http.Redirect(w, r, loginUrl.String(), http.StatusUnauthorized)
		return
	}

	context := newContext(router, user, session)
	handler.HandlerFunc(w, r, context)
}

func getUserForSession(session *sessions.Session, userDao *users.Dao) (*users.User, error) {
	userID := session.Values["userID"]
	if userID == nil {
		return nil, nil
	}

	return userDao.GetByID(userID.(uint64))
}