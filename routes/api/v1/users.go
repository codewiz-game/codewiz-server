package v1

import (
	"log"
	"path"
	"net/http"
	"github.com/gorilla/mux"
)

const (
	usersPath = "/users"
)

type User struct {
	Username string `json:"username"`
	Name string `json:"name"`
	TimeZone string `json:"timeZone"`
	Email string	`json:"emailAddress"`
}

func addUserRoutes(router *mux.Router, pathPrefixes ...string) {
	pathPrefixes = append(pathPrefixes, usersPath)
	pathPrefix := path.Join(pathPrefixes...)

	router.Path(path.Join(pathPrefix)).HandlerFunc(getAllUsersHandler).Methods("GET")
	router.Path(path.Join(pathPrefix)).HandlerFunc(addUserHandler).Methods("POST")

	router.Path(path.Join(pathPrefix, "/{id}")).HandlerFunc(getUserHandler).Methods("GET")
	router.Path(path.Join(pathPrefix, "/{id}")).HandlerFunc(modifyUserHandler).Methods("POST", "PUT")
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
