package v1

import (
	"net/http"
	"github.com/crob1140/codewiz/routes"
)

var (
	GetUser routes.Route
	AddUser routes.Route
)

type User struct {
	Username string `json:"username"`
	Name string `json:"name"`
	TimeZone string `json:"timeZone"`
	Email string	`json:"emailAddress"`
}

func addUserRoutes() {
	userRouter := router.PathPrefix("/users/").Subrouter()
	GetUser = userRouter.HandleFunc("/", getUserHandler).Methods("GET")
	AddUser = userRouter.HandleFunc("/", addUserHandler).Methods("POST")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {

}


func addUserHandler(w http.ResponseWriter, r *http.Request) {

}