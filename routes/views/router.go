package views

import (
	"github.com/crob1140/codewiz/models/users"
	"github.com/crob1140/codewiz/models/wizards"
	"github.com/crob1140/codewiz/config"
	"github.com/crob1140/codewiz/config/keys"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
	"path"
)


const (
	secondsPerHour    = 3600
)

type Router struct {
	*mux.Router

	path         string
	sessionStore sessions.Store

	userDao      *users.Dao
	wizardDao	 *wizards.Dao

	// Static URLs
	resourceURL *url.URL
	dashboardURL   *url.URL
	registrationURL *url.URL
	loginURL        *url.URL
	wizardListURL   *url.URL
	wizardCreationURL *url.URL

	// Dynamic URLs
	wizardViewRoute *mux.Route
}

func NewRouter(viewsPath string, userDao *users.Dao, wizardDao *wizards.Dao) http.Handler {

	// Initialise the session store with the necessary keys
	sessionStore := sessions.NewCookieStore([]byte(config.GetString(keys.SessionKey))) // TODO: read this directly from config? make it another arg?
	sessionStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: secondsPerHour,
		Secure: config.GetBool(keys.SessionSecure, false),	
	}

	// Create a new router instance with the obtained data
	router := &Router{Router: mux.NewRouter(), 
		path: viewsPath, 
		userDao: userDao,
		wizardDao : wizardDao,
		sessionStore: sessionStore,
	}

	// Add all of the routes to the router
	router.NotFoundHandler = &custom404Handler{}
	initRoutes(router)
	return router
}

func initRoutes(router *Router) {
	// Set static resource directory
	resourcePath := path.Join(router.path, "/resources")
	resourceHandler := http.StripPrefix(resourcePath, http.FileServer(http.Dir(resourceDirectory)))
	resourceRoute := router.PathPrefix(resourcePath).Handler(resourceHandler)
	router.resourceURL, _ = resourceRoute.URL()

	// Add dashboard page
	dashboardPath := router.path
	dashboardRoute := router.Path(dashboardPath).Handler(newHandler(dashboardPageHandler, router)).Methods("GET")
	router.dashboardURL, _ = dashboardRoute.URL()

	// Add registration page
	registrationPath := path.Join(router.path, "/register")
	registrationRoute := router.Path(registrationPath).Handler(newHandler(registerPageHandler, router)).Methods("GET")
	router.registrationURL, _ = registrationRoute.URL()
	router.Path(registrationPath).Handler(newHandler(registerActionHandler, router)).Methods("POST")

	// Add login page
	loginPath := path.Join(router.path, "/login")
	loginRoute := router.Path(loginPath).Handler(newHandler(loginPageHandler, router)).Methods("GET")
	router.loginURL, _ = loginRoute.URL() 
	router.Path(loginPath).Handler(newHandler(loginActionHandler, router)).Methods("POST")

	// Add wizard list page
	wizardListPath := path.Join(router.path, "/wizards")
	wizardListRoute := router.Path(wizardListPath).Handler(newHandler(listWizardsPageHandler, router)).Methods("GET")
	router.wizardListURL, _ = wizardListRoute.URL()

	// Add wizard creation page
	wizardCreationPath := path.Join(router.path, "/wizards/create")
	wizardCreationRoute := router.Path(wizardCreationPath).Handler(newHandler(createWizardPageHandler, router)).Methods("GET")
	router.wizardCreationURL, _ = wizardCreationRoute.URL()
	router.Path(wizardCreationPath).Handler(newHandler(createWizardActionHandler, router)).Methods("POST")

	// Add wizard view/update page
	wizardViewPath := path.Join(router.path, "/wizards/{id}")
	router.wizardViewRoute = router.Path(wizardViewPath).Handler(newHandler(viewWizardPageHandler, router)).Methods("GET")
	router.Path(wizardViewPath).Handler(newHandler(modifyWizardActionHandler, router)).Methods("POST")
}

func (router *Router) Dashboard() *url.URL {
	return router.dashboardURL
}

func (router *Router) Registration() *url.URL {
	return router.registrationURL
}

func (router *Router) Login() *url.URL {
	return router.loginURL
}

func (router *Router) WizardList() *url.URL {
	return router.wizardListURL
}

func (router *Router) WizardCreation() *url.URL {
	return router.wizardCreationURL
}

func (router *Router) WizardDetails(wizardID int) *url.URL {
	url, _ := router.wizardViewRoute.URL(string(wizardID))
	return url
}