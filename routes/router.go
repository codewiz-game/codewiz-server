package routes

import (
	"log"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
)

type Router struct {
 	*mux.Router
}

type subRouter struct {
	handler http.Handler
	pathPrefix string
}

func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

func (parentRouter *Router) AddSubrouter(pathPrefix string, router http.Handler) {
	
	if pathPrefix == "" || pathPrefix == "/" {
		log.Printf("case 1")
		// If the subrouter is to be shared on the same level as the parent router, then
		// there is no need to do any prefix manipulation
		parentRouter.PathPrefix(pathPrefix).Handler(router)
	} else {
		// Otherwise, create a router that strips the prefix off of the requests before they
		// are passed in, so that the sub-router can match against the remainder of the path
		subrouter := &subRouter{handler : http.StripPrefix(pathPrefix, router), pathPrefix : pathPrefix}

		// Serve the stripped subrouter at the given path prefix
		parentRouter.PathPrefix(pathPrefix).Handler(subrouter)
	}

}

func (router *subRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rr := httptest.NewRecorder()
	router.handler.ServeHTTP(rr,req)

	// If the router is attempting to redirect, the URL that it is redirecting to
	// will be missing the stripped path prefix, which must be re-inserted before
	// the redirection occurs.
	if rr.Code >= 300 && rr.Code < 400 {
		redirectUrl := rr.Header().Get("Location")
		log.Printf("Redirection changed from %s to %s", redirectUrl, router.pathPrefix + redirectUrl)
		rr.Header().Set("Location", router.pathPrefix + redirectUrl)
	}

	// Write the modified response to the original response writer
	for key, vals := range rr.Header() {
		for _, val := range vals {
			w.Header().Set(key, val)
		}
	}
	w.WriteHeader(rr.Code)
	w.Write(rr.Body.Bytes())
}