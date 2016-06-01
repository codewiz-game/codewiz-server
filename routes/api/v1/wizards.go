package v1

import (
	"log"
	"path"
	"net/http"
	"github.com/gorilla/mux"
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

func addWizardRoutes(router *mux.Router, pathPrefixes ...string) {
	pathPrefixes = append(pathPrefixes, wizardsPath)
	pathPrefix := path.Join(pathPrefixes...)

	router.HandleFunc(path.Join(pathPrefix), getAllWizardsHandler).Methods("GET")
	router.HandleFunc(path.Join(pathPrefix), addWizardHandler).Methods("POST")

	router.HandleFunc(path.Join(pathPrefix, "/{id}"), getWizardHandler).Methods("GET")
	router.HandleFunc(path.Join(pathPrefix, "/{id}"), modifyWizardHandler).Methods("POST")
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
