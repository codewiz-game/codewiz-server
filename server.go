package main

import (
	"github.com/crob1140/codewiz/datastore"
	"github.com/crob1140/codewiz/models/users"
	"github.com/crob1140/codewiz/routes/api"
	"github.com/crob1140/codewiz/routes/views"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	apiPath   = "/api"
	viewsPath = "/"
)

type CodewizServer struct {
	Router http.Handler
}

func NewServer(db *datastore.DB) *CodewizServer {
	userDao := users.NewDao(db)
	router := mux.NewRouter()

	// Add API endpoints
	apiRouter := api.NewRouter(apiPath)
	router.PathPrefix(apiPath).Handler(apiRouter)

	// Add view endpoints
	viewsRouter := views.NewRouter(viewsPath, userDao)
	router.PathPrefix(viewsPath).Handler(viewsRouter)

	return &CodewizServer{Router: router}
}

func (server *CodewizServer) ListenAndServe(address string) {
	http.ListenAndServe(address, server.Router)
}
