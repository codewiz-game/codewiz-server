package v1

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

// TODO: set up a test userDao

func TestGetUser_Existing(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getUserHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := `{"user" : true}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}