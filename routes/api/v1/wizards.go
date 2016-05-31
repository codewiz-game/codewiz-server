package v1

import (
	"log"
	"net/http"
	"github.com/crob1140/codewiz/routes"
)

const (
	wizardsPath = "/wizards"
)

type Spell struct {
	URI string `json:"uri"`
}

type Wizard struct {
	Name string `json:"name"`
	Spells []Spell `json:"spells"`
}

func addWizardRoutes(router *routes.Router) {
	router.HandleFunc(wizardsPath, getAllWizardsHandler).Methods("GET")
	router.HandleFunc(wizardsPath, addWizardHandler).Methods("POST")

	router.HandleFunc(wizardsPath + "/{id}", getWizardHandler).Methods("GET")
	router.HandleFunc(wizardsPath + "/{id}", modifyWizardHandler).Methods("POST")
}

func getAllWizardsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get all wizards v1")
}


func addWizardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Add wizard v1")
}

func getWizardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get wizard v1")
}

func modifyWizardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Modify wizard v1")
}
