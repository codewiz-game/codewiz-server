package views

import (
	"github.com/crob1140/codewiz-server/models"
	"github.com/crob1140/codewiz-server/models/users"
	"github.com/crob1140/codewiz-server/log"
	"net/http"
)


func loginPageHandler(w http.ResponseWriter, r *http.Request, context *context) {

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

	loginUrl := router.Login()
	data := struct {
		SubmitPath  string
		ValidationErrors models.ValidationErrors
	}{
		loginUrl.String(), 
		validationErrs,
	}

	render(w, "login.html", data)
}

func loginActionHandler(w http.ResponseWriter, r *http.Request, context *context) {

	router := context.Router
	session := context.Session

	user, validationErrs, err := validateLoginRequest(r, router.userDao)
	if err != nil {
		custom500Handler(w,r)
		return
	}

	if len(validationErrs) == 0 {
		// Log the user in by saving their username as a session attribute
		session.Values["userID"] = user.ID
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Debug("User has logged in", log.Fields{"username" : user.Username})

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

		// Send the user back to the login page
		loginUrl := router.Login()
		http.Redirect(w, r, loginUrl.String(), http.StatusSeeOther)
	}
}

func validateLoginRequest(r *http.Request, userDao *users.Dao) (*users.User, models.ValidationErrors, error) {
	errs := make(models.ValidationErrors)

	password := r.FormValue("password")
	if password == "" {
		errs.Add("Password", "This field is required.")
	}

	username := r.FormValue("username")
	var user *users.User
	if username == "" {
		errs.Add("Username", "This field is required.")
	} 

	if username != "" && password != "" {
		var err error
		user, err = userDao.GetByUsername(username)
		if err != nil {
			return nil, nil, err
		}

		if user == nil || !user.VerifyPassword(password) {
			errs.Add("Username", "The username or password you have entered is invalid.")
		}
	}

	return user, errs, nil
}