package views

import (
	"github.com/crob1140/codewiz/models/wizards"
	"github.com/gorilla/sessions"
	"net/http"
)


func dashboardPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

	user, err := getUserForSession(session, router.userDao)
	if err != nil {
		custom500Handler(w, r)
		return
	}

	if user == nil {
		custom401Handler(w, r)
		return
	}

	userWizards, err := router.wizardDao.GetByOwnerID(user.ID)
	if err != nil {
		custom500Handler(w, r)
	}

	data := struct {
		Username string
		Wizards []*wizards.Wizard
		CreateWizardPath string
	}{
		user.Username,
		userWizards,
		router.WizardCreation().String(),
	}

	render(w, "dashboard.html", data)
}