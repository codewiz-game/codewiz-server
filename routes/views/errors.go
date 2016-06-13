package views

import (
	"encoding/gob"
	"net/http"
)

type custom404Handler struct{}

type validationErrors struct {
	FormErrors []string
	FieldErrors map[string][]string
}


func init() {
	// Register the validation errors type to allow them to be encoded to the session cookie.
	gob.Register(validationErrors{})
}

func newValidationErrors(formErrors []string, fieldErrors map[string][]string) *validationErrors {
	return &validationErrors{FormErrors : formErrors, FieldErrors : fieldErrors}
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