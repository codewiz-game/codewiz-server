package api

import (
	"fmt"
	"net/http"
	"github.com/crob1140/codewiz/routes"
	"github.com/crob1140/codewiz/routes/api/v1"
)

const (
	// Versioning
	LatestVersion = 1
	versionPathFormat = "/v%d"
	latestVersionPath = "/latest"
)

func NewRouter() http.Handler {

	router := routes.NewRouter()

	// Add version one
	v1Path := fmt.Sprintf(versionPathFormat, 1)
	v1Router := v1.NewRouter()
	router.AddSubrouter(v1Path, v1Router)
	
	// ----------------------------------------------------------------
	// NOTE: new versions can be added here
	// Each version should be a complete implementation, rather than
	// just the differences from the previous version. This can be 
	// done by adding all the previous versions routes to the new router,
	// and overriding the changes.
	// ----------------------------------------------------------------

	router.AddSubrouter(latestVersionPath, v1Router)
	return router

}
