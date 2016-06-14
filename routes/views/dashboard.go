package views

import (
	"github.com/crob1140/codewiz/models/wizards"
	"net/http"
)


func dashboardPageHandler(w http.ResponseWriter, r *http.Request, context *context) {

	user := context.User
	router := context.Router

	if user != nil {
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
	} else {
		// TODO: show index page
	}
}