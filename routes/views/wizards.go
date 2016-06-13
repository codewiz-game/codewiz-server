package views

import (
	"github.com/crob1140/codewiz/log"
	"github.com/crob1140/codewiz/models/wizards"
	"github.com/gorilla/sessions"
	"net/http"
)

func listWizardsPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

}

func createWizardPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

	errList := session.Flashes("errs")
	var errs *validationErrors
	if len(errList) != 0 {
		validationErrs := errList[0].(validationErrors)
		errs = &validationErrs
	}

	// Save the session to ensure the flash messages are removed.
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		SubmitPath string
		Errors *validationErrors
	}{
		router.WizardCreation().String(),
		errs,
	}

	render(w, "createwizard.html", data)
}

func createWizardActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {
	
	user, err := getUserForSession(session, router.userDao)
	if err != nil {
		custom500Handler(w,r)
		return;
	}

	wizard, errs := validateWizardCreationRequest(r)
	if errs != nil {
		session.AddFlash(*errs, "errs")
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the user back to the creation page
		wizardCreationUrl := router.WizardCreation()
		http.Redirect(w, r, wizardCreationUrl.String(), http.StatusSeeOther)
	} else {
		wizard.OwnerID = user.ID
		if err := router.wizardDao.Insert(wizard); err != nil {
			log.Error("Error occurred while inserting wizard", log.Fields{"error" : err})
			custom500Handler(w,r)
			return
		}

		// Send the user back to the dasboard
		dashboardUrl := router.Dashboard()
		http.Redirect(w, r, dashboardUrl.String(), http.StatusSeeOther)
	}
}

func validateWizardCreationRequest(r *http.Request) (*wizards.Wizard, *validationErrors) {

	fieldErrs := make(map[string][]string)
	
	name := r.FormValue("name")
	if name == "" {
		fieldErrs["name"] = append(fieldErrs["name"], "Name must be provided.")
	}

	sex := r.FormValue("sex")
	if !(sex == "M" || sex == "F") {
		fieldErrs["sex"] = append(fieldErrs["sex"], "Sex must be either male or female.")
	}

	if len(fieldErrs) != 0 {
		return nil, newValidationErrors(nil, fieldErrs)
	}

	return wizards.NewWizard(name, sex, 0), nil

}

func viewWizardPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

}

func modifyWizardActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

}