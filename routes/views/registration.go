package views

import (
	"github.com/crob1140/codewiz/models"
	"github.com/crob1140/codewiz/models/users"
	"net/http"
)


func registerPageHandler(w http.ResponseWriter, r *http.Request, context *context) {
	
	router := context.Router
	session := context.Session

	// If the user is reaching this page after being
	// redirected due to a validation error, the errors
	// be error messages stored in a flash message which
	// needs to be read and loaded into the context.
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

	registerUrl := router.Registration()
	data := struct {
		SubmitPath  string
		ValidationErrors models.ValidationErrors
	}{
		registerUrl.String(), 
		validationErrs,
	}

	render(w, "register.html", data)
}

func registerActionHandler(w http.ResponseWriter, r *http.Request, context *context) {

	router := context.Router
	session := context.Session

	user := extractUserFromRegistrationRequest(r)
	validator := users.NewValidator(router.userDao)
	validationErrs, err := validator.Validate(user)
	if err != nil {
		custom500Handler(w,r)
		return
	}

	if len(validationErrs) == 0 {
		// Create a new User and store it in the database
		if err := router.userDao.Insert(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Log the user in by saving their username as a session attribute
		session.Values["userID"] = user.ID
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect the user to the dashboard
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

		// Send the user back to the registration page
		registerUrl := router.Registration()
		http.Redirect(w, r, registerUrl.String(), http.StatusSeeOther)
	}
}

func extractUserFromRegistrationRequest(r *http.Request) *users.User {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	return users.NewUser(username, password, email)
}
