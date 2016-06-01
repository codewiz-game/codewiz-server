package api

import (
	"fmt"
	"path"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/crob1140/codewiz/routes/api/v1"
)

const (
	// Versioning
	LatestVersion = 1
	versionPathFormat = "/v%d"
	latestVersionPath = "/latest"
)

func NewRouter(pathPrefixes ...string) http.Handler {

	router := mux.NewRouter()

	apiPath := path.Join(pathPrefixes...)

	// Add version one
	v1Path := path.Join(apiPath, fmt.Sprintf(versionPathFormat, 1))
	v1Router := v1.NewRouter(v1Path)
	router.PathPrefix(v1Path).Handler(v1Router)
	
	// ----------------------------------------------------------------
	// NOTE: new versions can be added here
	// Each version should be a complete implementation, rather than
	// just the differences from the previous version. This can be 
	// done by adding all the previous versions routes to the new router,
	// and overriding the changes.
	// ----------------------------------------------------------------

	latestPath := path.Join(apiPath, latestVersionPath)
	latestRouter := v1.NewRouter(latestPath)
	router.PathPrefix(latestPath).Handler(latestRouter)

	return router

}
