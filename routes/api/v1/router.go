package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"runtime/debug"
	"github.com/crob1140/codewiz/log"
	"github.com/crob1140/codewiz/routes"
	"github.com/crob1140/codewiz/models/users"
)

const (
	CodeInternalError = 50000
	
	// Authorization
	CodeAdminOnly = 40100
	CodeOwnerOnly = 40101
)

type Error struct {
	Message string `json:"message"`
	Code int `json:"code"`
}


func NewRouter(v1Path string, userDao *users.Dao) *routes.Router {

	router := routes.NewRouter(v1Path).StrictSlash(true)
	router.Use(createRecoveryMiddleware())
	router.Use(createUserIdentificationMiddleware(userDao))
	router.Use(createLoggerMiddleware())

	addUserRoutes(router)
	addWizardRoutes(router)

	return router
}

func createRecoveryMiddleware() routes.Middleware {
	return routes.MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *routes.Context, next routes.HandlerFunc) {
		defer func() {
			if recovery := recover(); recovery != nil {

				err, ok := recovery.(error)
            	if !ok {
              	  err = fmt.Errorf("pkg: %v", recovery)
            	}

            	username := "nil"
            	if context.User != nil {
            		username = context.User.Username
            	}

            	body, _ := ioutil.ReadAll(r.Body)
				log.Error("Panic occurred while serving API v1 request", log.Fields{
					"error" : err, 
					"url" : r.URL.String(), 
					"user" : username,
					"body" : body,
					"stack" : string(debug.Stack()),
				})


				w.WriteHeader(http.StatusInternalServerError)
				w.Write(toJson(Error{
					Message : "An internal server error has occurred.", 
					Code : CodeInternalError,
				})) 
			}
		}()

		next(w,r,context)
	})
}

func createLoggerMiddleware() routes.Middleware {
	return routes.MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *routes.Context, next routes.HandlerFunc) {
		
		username := "nil"
    	if context.User != nil {
    		username = context.User.Username
    	}


		log.Debug("Serving API v1 request", log.Fields{
			"url" : r.URL.String(), 
			"user" : username,
		})
		
		next(w,r,context)
	})
}

func createUserIdentificationMiddleware(userDao *users.Dao) routes.Middleware {
	return routes.MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *routes.Context, next routes.HandlerFunc) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := userDao.GetByUsername(username)
			if err != nil {
				log.Error("Failed to fetch user from datastore", log.Fields{
					"username" : username,
					"error" : err,
				})

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(toJson(Error{
					Message : "An internal server error has occurred.", 
					Code : CodeInternalError,
				}))

				return
			}

			if user != nil {
				if user.VerifyPassword(password) {
					context.User = user
				}
			} else {
				// TODO
			}
		}

		next(w,r,context)
	})
}

func toJson(obj interface{}) []byte {
	objJson, _ := json.Marshal(obj)
	return objJson
}