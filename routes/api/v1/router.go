package v1

import (
	"net/http"
	"github.com/crob1140/codewiz/routes"
)

func NewRouter() http.Handler {
	router := routes.NewRouter()
	addUserRoutes(router)
	addWizardRoutes(router)
	return router
}