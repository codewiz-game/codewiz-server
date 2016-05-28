package v1

import (
	"net/http"
	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	addUserRoutes()
	addWizardRoutes()	
}

func Handler() http.Handler {
	return router
}