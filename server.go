package main

import (
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/crob1140/codewiz/routes/api"
	"github.com/crob1140/codewiz/routes/views"
	"github.com/crob1140/codewiz/datastore"
	"github.com/crob1140/codewiz/models"


	_ "github.com/mattn/go-sqlite3"
)

type CodewizServer struct {
	*mux.Router
}

func NewServer() *CodewizServer {

	// TODO: read from config file to build a Configuration (which should really just be environment variables)

	// TODO: this should be the constructor for UserDao
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	ds := datastore.NewDataStore(db, datastore.SqliteDialect)
	userDao := models.NewUserDao(ds)

	router := mux.NewRouter()	
	viewrouter := views.NewRouter(userDao)
	router.PathPrefix("/").Handler(viewrouter)

/**
	apirouter := api.Router{UserDao : userDao, WizardDao : wizardDao} etc.
	router.Path("/api/").Handler(apirouter)
*/
	api.AddRoutes(router)

	return &CodewizServer{Router : router}
}

func (server *CodewizServer) ListenAndServe(address string) {
	http.ListenAndServe(address, server.Router)
}