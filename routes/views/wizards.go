package views

import (
	"github.com/crob1140/codewiz-server/log"
	"github.com/crob1140/codewiz-server/models"
	"github.com/crob1140/codewiz-server/models/wizards"
	"net/http"
)

func listWizardsPageHandler(w http.ResponseWriter, r *http.Request, context *context) {

}

func createWizardPageHandler(w http.ResponseWriter, r *http.Request, context *context) {

	router := context.Router
	session := context.Session

	errList := session.Flashes("errs")
	var validationErrs models.ValidationErrors
	if len(errList) != 0 {
		validationErrs = errList[0].(models.ValidationErrors)
	}

	// Save the session to ensure the flash messages are removed.
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wizardCreationPath := router.WizardCreation().String()

	data := struct {
		SubmitPath string
		ValidationErrors models.ValidationErrors
	}{
		wizardCreationPath,
		validationErrs,
	}

	render(w, "createwizard.html", data)
}

func createWizardActionHandler(w http.ResponseWriter, r *http.Request, context *context) {
	
	user := context.User
	router := context.Router
	session := context.Session

	wizard := extractWizardFromRequest(r, user.ID)
	validator := wizards.NewValidator(router.wizardDao)
	validationErrs, err := validator.Validate(wizard)
	if err != nil {
		custom500Handler(w,r)
		return
	} 

	if len(validationErrs) == 0 {
		if err := router.wizardDao.Insert(wizard); err != nil {
			log.Error("Error occurred while inserting wizard", log.Fields{"error" : err})
			custom500Handler(w,r)
			return
		}

		// Send the user back to the dasboard
		dashboardUrl := router.Dashboard()
		http.Redirect(w, r, dashboardUrl.String(), http.StatusSeeOther)
	} else {
		// Add the errors to a flash message so that we can access them
		// after redirection
		session.AddFlash(validationErrs, "errs")
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the user back to the creation page
		wizardCreationUrl := router.WizardCreation()
		http.Redirect(w, r, wizardCreationUrl.String(), http.StatusSeeOther)
	}
}

func extractWizardFromRequest(r *http.Request, userID uint64) *wizards.Wizard {
	name := r.FormValue("name")
	sex := r.FormValue("sex")
	return wizards.NewWizard(name, sex, userID)
}

func viewWizardPageHandler(w http.ResponseWriter, r *http.Request, context *context) {

}

func modifyWizardActionHandler(w http.ResponseWriter, r *http.Request, context *context) {

}