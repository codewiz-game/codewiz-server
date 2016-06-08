package views

import (
	"encoding/gob"
	"github.com/crob1140/codewiz/models/users"
	"github.com/crob1140/codewiz/config"
	"github.com/crob1140/codewiz/config/keys"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
)

const (
	resourceDirectory = "routes/views/resources"
	TemplateDirectory = resourceDirectory + "/templates"
	sessionName       = "codewiz-session"
	secondsPerHour    = 3600
)

type handler struct {
	router      *viewsRouter
	handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *viewsRouter)
}

type viewsRouter struct {
	*mux.Router

	path         string
	sessionStore sessions.Store
	userDao      *users.Dao

	// TODO: replace all with just URL's? only ever used for redirection
	dashboardRoute    *mux.Route
	registrationRoute *mux.Route
	loginRoute        *mux.Route
}

func init() {

	// Register all of the types that will be stored as session values
	// to allow them to be encoded to the session cookie
	gob.Register(map[string][]string{})
}

// ----------------------------------------------------------------------
// 404 Customisation
// ----------------------------------------------------------------------

type custom404Handler struct{}

func (*custom404Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	render(w, "404.html", nil)
}

// ----------------------------------------------------------------------

// This method performs all of the common code and passes down the
// frequently used components to the handlers.
func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get or create a session cookie
	session, err := handler.router.sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.handlerFunc(w, r, session, handler.router)
}


func NewRouter(viewsPath string, userDao *users.Dao) http.Handler {

	// Initialise the session store with the necessary keys
	sessionStore := sessions.NewCookieStore([]byte(config.GetString(keys.SessionKey))) // TODO: read this directly from config? make it another arg?
	sessionStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: secondsPerHour,
		Secure: config.GetBool(keys.SessionSecure, false),	
	}

	viewsRouter := &viewsRouter{Router: mux.NewRouter(), path: viewsPath, userDao: userDao, sessionStore: sessionStore}
	viewsRouter.NotFoundHandler = &custom404Handler{}
	initRoutes(viewsRouter)
	return viewsRouter
}

func initRoutes(router *viewsRouter) {

	// Set static resource directory
	resourceHandler := http.StripPrefix("/resources", http.FileServer(http.Dir(resourceDirectory)))
	router.PathPrefix(path.Join(router.path, "/resources")).Handler(resourceHandler)

	// Add dashboard page
	router.dashboardRoute = router.Path(router.path).Handler(createHandler(dashboardPageHandler, router)).Methods("GET")

	// Add registration page
	router.registrationRoute = router.Path(path.Join(router.path, "/register")).Handler(createHandler(registerPageHandler, router)).Methods("GET")
	router.Path(path.Join(router.path, "/register")).Handler(createHandler(registerActionHandler, router)).Methods("POST")

	// Add login page
	router.loginRoute = router.Path(path.Join(router.path, "/login")).Handler(createHandler(loginPageHandler, router)).Methods("GET")
	router.Path(path.Join(router.path, "/login")).Handler(createHandler(loginActionHandler, router)).Methods("POST")

}

func createHandler(handlerFunc func(http.ResponseWriter, *http.Request, *sessions.Session, *viewsRouter), router *viewsRouter) *handler {
	return &handler{router: router, handlerFunc: handlerFunc}
}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	path, _ := filepath.Abs(TemplateDirectory + "/" + templateName)
	tmpl, _ := template.ParseFiles(path)
	tmpl.Execute(w, data)
}

func dashboardPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *viewsRouter) {
	data := struct {
		Username string
	}{session.Values["username"].(string)}

	render(w, "dashboard.html", data)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *viewsRouter) {
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

func registerActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *viewsRouter) {

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

func loginPageHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *viewsRouter) {
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

func loginActionHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, router *viewsRouter) {
	
	username, password, errs := validateLoginRequest(r)

	if len(errs) == 0 {
		// Get the user from the database
		user, err := router.userDao.Get(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user != nil && user.VerifyPassword(password) {
			// Log the user in by saving their username as a session attribute
			session.Values["username"] = user.Username
			if err := session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

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

func (router *viewsRouter) Dashboard() *url.URL {
	url, _ := router.dashboardRoute.URL()
	return url
}

func (router *viewsRouter) Registration() *url.URL {
	url, _ := router.registrationRoute.URL()
	return url
}

func (router *viewsRouter) Login() *url.URL {
	url, _ := router.loginRoute.URL()
	return url
}
