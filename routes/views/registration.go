package views

import (
	"github.com/gorilla/sessions"
	"github.com/crob1140/codewiz/models/users"
	"net/http"
)


func registerPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *router) {
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

	registerUrl := router.Registration()
	data := struct {
		SubmitPath  string
		FieldErrors map[string][]string
	}{registerUrl.String(), errs}

	render(w, "register.html", data)
}

func registerActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *router) {

	user, errs := validateRegistrationRequest(r)

	if len(errs) == 0 {
		// Create a new User and store it in the database
		if err := router.userDao.Insert(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Log the user in by saving their username as a session attribute
		session.Values["username"] = user.Username
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
		session.AddFlash(errs, "errs")
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the user back to the registration page
		registerUrl := router.Registration()
		http.Redirect(w, r, registerUrl.String(), http.StatusSeeOther)
	}
}

func validateRegistrationRequest(r *http.Request) (*users.User, map[string][]string) {
	errs := make(map[string][]string)

	username := r.FormValue("username")
	if username == "" {
		errs["username"] = append(errs["username"], "Username must be provided.")
	}

	password := r.FormValue("password")
	if password == "" {
		errs["password"] = append(errs["password"], "Password must be provided.")
	}

	email := r.FormValue("email")
	if email == "" {
		errs["email"] = append(errs["email"], "Email must be provided.")
	}

	if len(errs) != 0 {
		return nil, errs
	}

	user := users.NewUser(username, password, email)
	return user, errs
}

