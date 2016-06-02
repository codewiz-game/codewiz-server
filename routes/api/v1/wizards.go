package v1

import (
	"log"
	"path"
	"net/http"
	"github.com/gorilla/mux"
)

type Spell struct {
	URI string `json:"uri"`
}

type Wizard struct {
	Name string `json:"name"`
	Spells []Spell `json:"spells"`
}

func addWizardRoutes(router *mux.Router, v1Path string) {
	wizardsPath := path.Join(v1Path, "/wizards")

	router.HandleFunc(wizardsPath, getAllWizardsHandler).Methods("GET")
	router.HandleFunc(wizardsPath, addWizardHandler).Methods("POST")

	router.HandleFunc(path.Join(wizardsPath, "/{id}"), getWizardHandler).Methods("GET")
	router.HandleFunc(path.Join(wizardsPath, "/{id}"), modifyWizardHandler).Methods("POST")
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
