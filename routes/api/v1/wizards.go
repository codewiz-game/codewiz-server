package v1

import (
	"net/http"
	"github.com/crob1140/codewiz/routes"
)

var (
	GetWizard routes.Route
	AddWizard routes.Route
)

type Spell struct {
	URI string `json:"uri"`
}

type Wizard struct {
	Name string `json:"name"`
	Spells []Spell `json:"spells"`
}

func addWizardRoutes() {
	wizardRouter := router.PathPrefix("/wizards/").Subrouter()
	GetWizard = wizardRouter.HandleFunc("/", getWizardHandler).Methods("GET")
	AddWizard = wizardRouter.HandleFunc("/", addWizardHandler).Methods("POST")
}


func getWizardHandler(w http.ResponseWriter, r *http.Request) {
	
}


func addWizardHandler(w http.ResponseWriter, r *http.Request) {

}