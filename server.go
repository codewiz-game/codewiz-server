package main

import (
	"github.com/crob1140/codewiz/datastore"
	"github.com/crob1140/codewiz/models/users"
	"github.com/crob1140/codewiz/models/wizards"
	"github.com/crob1140/codewiz/routes/api"
	"github.com/crob1140/codewiz/routes/views"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	apiPath   = "/api"
	viewsPath = "/"
)

type Server struct {
	Router http.Handler
}

func NewServer(db *datastore.DB) *Server {
	userDao := users.NewDao(db)
	wizardDao := wizards.NewDao(db)

	router := mux.NewRouter()

	// Add API endpoints
	apiRouter := api.NewRouter(apiPath)
	router.PathPrefix(apiPath).Handler(apiRouter)

	// Add view endpoints
	viewsRouter := views.NewRouter(viewsPath, userDao, wizardDao)
	router.PathPrefix(viewsPath).Handler(viewsRouter)

	return &Server{Router: router}
}

func (server *Server) ListenAndServe(address string) {
	http.ListenAndServe(address, server.Router)
}
