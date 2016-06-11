package views

import (
	"github.com/gorilla/sessions"
	"net/http"
)


func dashboardPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *router) {
	data := struct {
		Username string
	}{session.Values["username"].(string)}

	render(w, "dashboard.html", data)
}