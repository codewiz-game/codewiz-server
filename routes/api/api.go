package api

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/crob1140/codewiz/routes/api/v1"
)

const (
	// Versioning
	LatestVersion = 1
	versionPathFormat = "/v%d/"
	latestVersionPath = "/latest/"
)

func Handler() http.Handler {
	router := mux.NewRouter()

	// Add version one
	v1Path := fmt.Sprintf(versionPathFormat, 1)
	v1Router := router.Handle(v1Path, v1.Router()).Subrouter()

	// ----------------------------------------------------------------
	// NOTE: new versions can be added here
	// Each version should be a complete implementation, rather than
	// just the differences from the previous version. This can be 
	// done by importing any unchanged routes from the previous version
	// and attaching them to the new version's handler instance.
	// ----------------------------------------------------------------

	// Use the most recent versions handlers for the lastVersionPath
	latest := v1Router
	router.Handle(latestVersionPath, latest)
	
	return router
}