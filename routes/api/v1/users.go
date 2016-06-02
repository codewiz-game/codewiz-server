package v1

import (
	"log"
	"path"
	"net/http"
	"github.com/gorilla/mux"
)


type User struct {
	Username string `json:"username"`
	Name string `json:"name"`
	TimeZone string `json:"timeZone"`
	Email string	`json:"emailAddress"`
}

func addUserRoutes(router *mux.Router, v1Path string) {
	usersPath := path.Join(v1Path, "/users")

	router.Path(usersPath).HandlerFunc(getAllUsersHandler).Methods("GET")
	router.Path(usersPath).HandlerFunc(addUserHandler).Methods("POST")

	router.Path(path.Join(usersPath, "/{id}")).HandlerFunc(getUserHandler).Methods("GET")
	router.Path(path.Join(usersPath, "/{id}")).HandlerFunc(modifyUserHandler).Methods("POST", "PUT")
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get all users v1")
	w.Write([]byte(`{"user" : true}`))
}


func addUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Add user v1")
	w.Write([]byte(`{"user" : true}`))
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get user v1")
	w.Write([]byte(`{"user" : true}`))
}

func modifyUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Modify user v1")
	w.Write([]byte(`{"user" : true}`))
}
