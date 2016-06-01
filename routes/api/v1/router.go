package v1

import (
	"net/http"
	"github.com/gorilla/mux"
)

func NewRouter(pathPrefixes ...string) http.Handler {
	router := mux.NewRouter()
	addUserRoutes(router, pathPrefixes...)
	addWizardRoutes(router, pathPrefixes...)
	return router
}