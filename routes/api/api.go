package api

import (
	"path"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/crob1140/codewiz-server/models/users"
	"github.com/crob1140/codewiz-server/routes/api/v1"
)

func NewRouter(apiPath string, userDao *users.Dao) http.Handler {

	router := mux.NewRouter()

	// Add version one
	v1Path := path.Join(apiPath, "/v1")
	v1Router := v1.NewRouter(v1Path, userDao)
	router.PathPrefix(v1Path).Handler(v1Router)
	
	// ----------------------------------------------------------------
	// NOTE: new versions can be added here
	// Each version should be a complete implementation, rather than
	// just the differences from the previous version. This can be 
	// done by adding all the previous versions routes to the new router,
	// and overriding the changes.
	// ----------------------------------------------------------------

	latestVersionPath := path.Join(apiPath, "/latest")
	latestVersionRouter := v1.NewRouter(latestVersionPath, userDao)
	router.PathPrefix(latestVersionPath).Handler(latestVersionRouter)

	return router

}