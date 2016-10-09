package v1

import (
	"path"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"github.com/crob1140/codewiz-server/routes"
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


func addUserRoutes(router *routes.Router) {
	router.Path(usersPath).HandlerFunc(getAllUsersHandler).Methods("GET")
	router.Path(usersPath).HandlerFunc(addUserHandler).Methods("POST")
	router.Path(path.Join(usersPath, "/{id:[0-9]+}")).HandlerFunc(getUserHandler).Methods("GET")
	router.Path(path.Join(usersPath, "/{id:[0-9]+}")).HandlerFunc(modifyUserHandler).Methods("POST", "PUT")
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	user := context.User

	isAuthorised := true // user.Role == users.AdminRole
	if isAuthorised{
		w.WriteHeader(http.StatusOK)
		w.Write(toJson(User{
			Username : user.Username,
			// Name : "TODO",
			Email : user.Email,
			// TimeZone : "TODO",	
		}))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(toJson(Error{
			Message : "Resource is only available to admin users.",
			Code : CodeAdminOnly,
		}))
	}
}


func addUserHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	// TODO: auth with user??
	w.Write([]byte(`{"user" : true}`))
}

func getUserHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	user := context.User
	pathUserID, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32) 

	isAuthorised := user != nil && (user.ID == pathUserID) // || user.Role == users.AdminRole)
	if isAuthorised{
		w.WriteHeader(http.StatusOK)
		w.Write(toJson(User{
			Username : user.Username,
			// Name : "TODO",
			Email : user.Email,
			// TimeZone : "TODO",	
		}))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(toJson(Error{
			Message : "User does not have permission to access this resource.",
			Code : CodeOwnerOnly,
		}))
	}
}

func modifyUserHandler(w http.ResponseWriter, r *http.Request, context *routes.Context) {
	w.Write([]byte(`{"user" : true}`))
}
