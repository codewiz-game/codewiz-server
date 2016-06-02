package v1

import (
	"net/http"
	"github.com/gorilla/mux"
)

func NewRouter(v1Path string) http.Handler {
	router := mux.NewRouter()
	addUserRoutes(router, v1Path)
	addWizardRoutes(router, v1Path)
	return router
}