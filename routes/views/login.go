package views

import (
	"github.com/crob1140/codewiz/log"
	"github.com/gorilla/sessions"
	"net/http"
)


func loginPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {
	// If the user is reaching this page after being
	// redirected due to a validation error, the errors
	// be error messages stored in a flash message which
	// needs to be read and loaded into the context.
	errList := session.Flashes("errs")

	// Save the session to ensure the flash messages are removed.
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var errs map[string][]string
	if len(errList) != 0 {
		errs = errList[0].(map[string][]string)
	} else {
		// Make an empty map to save having to check for nil values
		// in the template
		errs = make(map[string][]string)
	}

	loginUrl := router.Login()
	data := struct {
		SubmitPath  string
		FieldErrors map[string][]string
	}{loginUrl.String(), errs}

	render(w, "login.html", data)
}

func loginActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

	username, password, errs := validateLoginRequest(r)

	if len(errs) == 0 {
		// Get the user from the database
		user, err := router.userDao.GetByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user != nil && user.VerifyPassword(password) {
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
			return
		}

		// TODO: this is more of a general error, not field specific...need a way to handle these
		errs["username"] = append(errs["username"], "The username or password you have entered is invalid.")		
	} 

	// Add the errors to a flash message so that we can access them
	// after redirection
	session.AddFlash(errs, "errs")
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the user back to the login page
	loginUrl := router.Login()
	http.Redirect(w, r, loginUrl.String(), http.StatusSeeOther)
}

func validateLoginRequest(r *http.Request) (string, string, map[string][]string) {
	errs := make(map[string][]string)

	username := r.FormValue("username")
	if username == "" {
		errs["username"] = append(errs["username"], "Username must be provided.")
	}

	password := r.FormValue("password")
	if password == "" {
		errs["password"] = append(errs["password"], "Password must be provided.")
	}

	return username, password, errs
}
