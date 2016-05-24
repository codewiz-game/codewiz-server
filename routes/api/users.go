package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

const (
	GET_ALL_USERS_ROUTE string = "api.getallusers"
	CREATE_USER_ROUTE string = "api.createuser"
	GET_USER_ROUTE string = "api.getuser"
	UPDATE_USER_ROUTE string = "api.updateuser"
)

type User struct {
	Username string `json:"username"`
	Name string `json:"name"`
	TimeZone string `json:"timeZone"`
	Email string	`json:"emailAddress"`
}

func addUserRoutes(router *mux.Router) {

	subrouter := router.PathPrefix("/users/").Subrouter()

	subrouter.HandleFunc("/", GetAllUsersHandler).Methods("GET").Name(GET_ALL_USERS_ROUTE)
	subrouter.HandleFunc("/", CreateUserHandler).Methods("POST").Name(CREATE_USER_ROUTE)
	
	subrouter.HandleFunc("/{username}", GetUserHandler).Methods("GET").Name(GET_USER_ROUTE)
	subrouter.HandleFunc("/{username}", UpdateUserHandler).Methods("PUT").Name(UPDATE_USER_ROUTE)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}