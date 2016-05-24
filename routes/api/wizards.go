package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

const (
	GET_ALL_WIZARDS_ROUTE = "api.getallwizards"
	CREATE_WIZARD_ROUTE string = "api.createwizard"
)

type Spell struct {
	URI string `json:"uri"`
}

type Wizard struct {
	Name string `json:"name"`
	Spells []Spell `json:"spells"`
}

func addWizardRoutes(router *mux.Router) {
/**
	subrouter := router.PathPrefix("/wizards/").Subrouter()
	subrouter.Methods("GET").HandleFunc("/", GetAllWizardsHandler).Name(GET_ALL_WIZARDS_API_ROUTE)
	subrouter.Methods("POST").HandleFunc("/", CreateWizardHandler).Name(CREATE_WIZARD_API_ROUTE)
*/
}

func GetAllWizardsHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func CreateWizardHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}