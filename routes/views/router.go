package views

import (
	"github.com/crob1140/codewiz/models/users"
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

type router struct {
	*mux.Router

	path         string
	sessionStore sessions.Store
	userDao      *users.Dao

	// Static URL's
	resourceURL *url.URL
	dashboardURL   *url.URL
	registrationURL *url.URL
	loginURL        *url.URL
}

func NewRouter(viewsPath string, userDao *users.Dao) http.Handler {

	// Initialise the session store with the necessary keys
	sessionStore := sessions.NewCookieStore([]byte(config.GetString(keys.SessionKey))) // TODO: read this directly from config? make it another arg?
	sessionStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: secondsPerHour,
		Secure: config.GetBool(keys.SessionSecure, false),	
	}

	// Create a new router instance with the obtained data
	router := &router{Router: mux.NewRouter(), 
		path: viewsPath, 
		userDao: userDao, 
		sessionStore: sessionStore,
	}

	// Add all of the routes to the router
	router.NotFoundHandler = &custom404Handler{}
	initRoutes(router)
	return router
}

func initRoutes(router *router) {
	// Set static resource directory
	resourcePath := path.Join(router.path, "/resources")
	resourceHandler := http.StripPrefix(resourcePath, http.FileServer(http.Dir(resourceDirectory)))
	resourceRoute := router.PathPrefix(resourcePath).Handler(resourceHandler)
	router.resourceURL, _ = resourceRoute.URL()

	// Add dashboard page
	dashboardPath := router.path
	dashboardRoute := router.Path(dashboardPath).Handler(createHandler(dashboardPageHandler, router)).Methods("GET")
	router.dashboardURL, _ = dashboardRoute.URL()

	// Add registration page
	registrationPath := path.Join(router.path, "/register")
	registrationRoute := router.Path(registrationPath).Handler(createHandler(registerPageHandler, router)).Methods("GET")
	router.registrationURL, _ = registrationRoute.URL()
	router.Path(registrationPath).Handler(createHandler(registerActionHandler, router)).Methods("POST")

	// Add login page
	loginPath := path.Join(router.path, "/login")
	loginRoute := router.Path(loginPath).Handler(createHandler(loginPageHandler, router)).Methods("GET")
	router.loginURL, _ = loginRoute.URL() 
	router.Path(loginPath).Handler(createHandler(loginActionHandler, router)).Methods("POST")
}

func (router *router) Dashboard() *url.URL {
	return router.dashboardURL
}

func (router *router) Registration() *url.URL {
	return router.registrationURL
}

func (router *router) Login() *url.URL {
	return router.loginURL
}


