package main

import (
	"net/http"
	"github.com/crob1140/codewiz/routes"
	"github.com/crob1140/codewiz/routes/api"
	"github.com/crob1140/codewiz/routes/views"
	"github.com/crob1140/codewiz/datastore"
	"github.com/crob1140/codewiz/models"
)

const (
	apiPath = "/api"
	viewsPath = "/"
)

type CodewizServer struct {
	Router http.Handler
}

func NewServer(ds *datastore.SQLDataStore) *CodewizServer {
	userDao := models.NewUserDao(ds)
	router := routes.NewRouter()

	// Add API endpoints
	apiRouter := api.NewRouter()
	router.AddSubrouter(apiPath, apiRouter)

	// Add view endpoints
	viewsRouter := views.NewRouter(userDao)
	router.AddSubrouter(viewsPath, viewsRouter)

	return &CodewizServer{Router : router}
}

func (server *CodewizServer) ListenAndServe(address string) {
	http.ListenAndServe(address, server.Router)
}
