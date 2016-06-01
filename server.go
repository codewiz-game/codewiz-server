package main

import (
	"net/http"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	// Add API endpoints
	apiRouter := api.NewRouter(apiPath)
	router.PathPrefix(apiPath).Handler(apiRouter)

	// Add view endpoints
	viewsRouter := views.NewRouter(userDao, viewsPath)
	router.PathPrefix(viewsPath).Handler(viewsRouter)

	return &CodewizServer{Router : router}
}

func (server *CodewizServer) ListenAndServe(address string) {
	http.ListenAndServe(address, server.Router)
}
