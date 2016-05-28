package views

import (
	"encoding/gob"
	"net/http"
	"path/filepath"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/crob1140/codewiz/models"
	"github.com/crob1140/codewiz/routes"
	"html/template"
)

const (
	resourceDirectory = "routes/views/resources"
	sessionName = "codewiz-session"
	secondsPerHour = 3600
)

type handler struct {
	router *Router
	handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *Router)
}

type Router struct {
	muxRouter *mux.Router
	sessionStore sessions.Store
	userDao *models.UserDao

	HomePage routes.Route
	Dashboard routes.Route
	Registration routes.Route
	Login routes.Route

	// The following routes should not be exported, 
	// since they should only be accessed through
	// actions made on the views
	registrationSubmit routes.Route
}

func init() {

	// Register all of the types that will be stored as session values
	// to allow them to be encoded to the session cookie
	gob.Register(map[string][]string{})
}

func createHandler(handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *Router), router *Router) *handler {
	return &handler{router : router, handlerFunc : handlerFunc}
}

func NewRouter(userDao *models.UserDao) http.Handler {
	// Initialise the session store with the necessary keys
	sessionStore := sessions.NewCookieStore([]byte("top-secret-keks"))
	router := &Router{muxRouter : mux.NewRouter(), userDao : userDao, sessionStore : sessionStore}
	initRoutes(router)
	return router
}

func initRoutes(router *Router) {
	// Set static resource directory
	resourceHandler := http.StripPrefix("/resources", http.FileServer(http.Dir(resourceDirectory)))
	router.muxRouter.PathPrefix("/resources").Handler(resourceHandler)

	// Add dashboard page
	router.HomePage = router.muxRouter.Path("/").Handler(createHandler(dashboardPageHandler, router)).Methods("GET")
	router.Dashboard = router.muxRouter.Path("/dashboard").Handler(createHandler(dashboardPageHandler, router)).Methods("GET")
	
	// Add registration page
	router.Registration = router.muxRouter.Path("/register").Handler(createHandler(registerPageHandler, router)).Methods("GET")
	router.registrationSubmit = router.muxRouter.Path("/register").Handler(createHandler(registerActionHandler, router)).Methods("POST")
	
	// Add login page
	router.Login = router.muxRouter.Path("/login").Handler(createHandler(loginPageHandler, router)).Methods("GET")
}

// This method performs all of the common code and passes down the
// frequently used components to the handlers.
func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	// Get or create a session cookie
	session, err := handler.router.sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialise the session if the cookie didn't already exist
	if session.IsNew {
		session.Options.Path = "/"
		session.Options.MaxAge = secondsPerHour
		session.Options.Secure = true
	}

	handler.handlerFunc(w, r, session, handler.router)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.muxRouter.ServeHTTP(w,r)
}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	path, _ := filepath.Abs(resourceDirectory + "/templates/" + templateName)
	tmpl, _ := template.ParseFiles(path)
	tmpl.Execute(w, data)
}

func validateRegistrationRequest(r *http.Request) (*models.User, map[string][]string) {	
	errs := make(map[string][]string)

	username := r.FormValue("username")
	if username == "" {
		errs["username"] = append(errs["username"], "Username must be provided.");
	}

	password := r.FormValue("password")
	if password == "" {
		errs["password"] = append(errs["password"], "Password must be provided.");
	}

	email := r.FormValue("email")
	if email == "" {
		errs["email"] = append(errs["email"], "Email must be provided.");
	}

	if len(errs) != 0 {
		return nil, errs
	}

	user := models.NewUser(username, password, email)
	return user, errs
}


func dashboardPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {
	data := struct {
		Username string
	}{ session.Values["username"].(string) }

	render(w, "dashboard.html", data)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {
	// If the user is reaching this page after being 
	// redirected due to a validation error, the errors
	// be error messages stored in a flash message which
	// needs to be read and loaded into the context.
	errList := session.Flashes("errs")
	var errs map[string][]string
	if len(errList) != 0 {
		errs = errList[0].(map[string][]string)
	} else {
		// Make an empty map to save having to check for nil values
		// in the template
		errs = make(map[string][]string) 
	}

	registerUrl, _ := router.registrationSubmit.URL()
	data := struct {
		SubmitPath string
		FieldErrors map[string][]string
	}{ registerUrl.String(), errs }

	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, "register.html", data)	
}

func registerActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {
	
	user, errs := validateRegistrationRequest(r)

	if len(errs) == 0 {
		// Create a new User and store it in the database
		if err := router.userDao.AddUser(user); err != nil {
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
		dashboardUrl, _ := router.Dashboard.URL()
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
		registerUrl, _ := router.Registration.URL()
		http.Redirect(w, r, registerUrl.String(), http.StatusSeeOther)	
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *Router) {

}
