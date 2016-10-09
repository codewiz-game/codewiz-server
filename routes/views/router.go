package views

import (
	"github.com/crob1140/codewiz-server/models/users"
	"github.com/crob1140/codewiz-server/models/wizards"
	"github.com/crob1140/codewiz-server/config"
	"github.com/crob1140/codewiz-server/config/keys"
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
	dashboardRoute := router.addHandler("GET", dashboardPath, dashboardPageHandler, false)
	router.dashboardURL, _ = dashboardRoute.URL()

	// Add registration page
	registrationPath := path.Join(router.path, "/register")
	registrationRoute := router.addHandler("GET", registrationPath, registerPageHandler, false)
	router.registrationURL, _ = registrationRoute.URL()
	router.addHandler("POST", registrationPath, registerActionHandler, false)

	// Add login page
	loginPath := path.Join(router.path, "/login")
	loginRoute := router.addHandler("GET", loginPath, loginPageHandler, false)
	router.loginURL, _ = loginRoute.URL() 
	router.addHandler("POST", loginPath, loginActionHandler, false)

	// Add wizard list page
	wizardListPath := path.Join(router.path, "/wizards")
	wizardListRoute := router.addHandler("GET", wizardListPath, listWizardsPageHandler, true)
	router.wizardListURL, _ = wizardListRoute.URL()

	// Add wizard creation page
	wizardCreationPath := path.Join(router.path, "/wizards/create")
	wizardCreationRoute := router.addHandler("GET", wizardCreationPath, createWizardPageHandler, true)
	router.wizardCreationURL, _ = wizardCreationRoute.URL()
	router.addHandler("POST", wizardCreationPath, createWizardActionHandler, true)

	// Add wizard view/update page
	wizardViewPath := path.Join(router.path, "/wizards/{id}")
	router.wizardViewRoute = router.addHandler("GET", wizardViewPath, viewWizardPageHandler, true)
	router.addHandler("POST", wizardViewPath, modifyWizardActionHandler, true)
}

func (router *Router) addHandler(method string, path string, handlerFunc handlerFunc, requiresLogin bool) *mux.Route {
	return router.Path(path).Handler(newHandler(handlerFunc, router, requiresLogin)).Methods(method)
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