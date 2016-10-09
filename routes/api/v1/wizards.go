package v1

import (
	"log"
	"path"
	"net/http"
	"github.com/crob1140/codewiz/routes"
)


type Spell struct {
	URI string `json:"uri"`
}

type Wizard struct {
	Name string `json:"name"`
	Spells []Spell `json:"spells"`
}

func addWizardRoutes(router *routes.Router) {
	wizardsPath := "/wizards"

	router.Path(wizardsPath).HandlerFunc(getAllWizardsHandler).Methods("GET")
	router.Path(wizardsPath).HandlerFunc(addWizardHandler).Methods("POST")

	router.Path(path.Join(wizardsPath, "/{id}")).HandlerFunc(getWizardHandler).Methods("GET")
	router.Path(path.Join(wizardsPath, "/{id}")).HandlerFunc(modifyWizardHandler).Methods("POST")
}

func getAllWizardsHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	log.Printf("Get all wizards v1")
}


func addWizardHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	log.Printf("Add wizard v1")
}

func getWizardHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	log.Printf("Get wizard v1")
}

func modifyWizardHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	log.Printf("Modify wizard v1")
}
