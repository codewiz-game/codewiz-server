package views

import (
	"github.com/crob1140/codewiz-server/models"
	"encoding/gob"
	"net/http"
)

type custom404Handler struct{}

func init() {
	// Register the field errors type to allow them to be 
	// encoded as flash messages in the session cookie.
	gob.Register(make(models.ValidationErrors))
}

// 404 is a special-case error, since the handler function cannot be called from a handler.
// In order to deal with this, we need to create a custom http.Handler that will catch
// all routes that are not handled by the main router.
func (*custom404Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	render(w, "404.html", nil)
}

func custom500Handler(w http.ResponseWriter, r *http.Request) {

}

func custom401Handler(w http.ResponseWriter, r *http.Request) {

}