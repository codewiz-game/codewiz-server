package api

import (
	"github.com/gorilla/mux"
)

const (
	VERSION_1_ENDPOINT = "/v1/"
	LATEST_API_ENDPOINT = VERSION_1_ENDPOINT
)

func AddRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/api/").Subrouter()
	addUserRoutes(subrouter)
	addWizardRoutes(subrouter)
}