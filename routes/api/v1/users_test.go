package v1

import (
    "os"
	"testing"
    paths "path"
    "strings"
	"net/http"
    "encoding/base64"
	"net/http/httptest"
    "github.com/crob1140/codewiz/config"
    "github.com/crob1140/codewiz/config/keys"
    "github.com/crob1140/codewiz/datastore"
    "github.com/crob1140/codewiz/models/users"
    "github.com/crob1140/codewiz/routes"
    _ "github.com/mattn/go-sqlite3"
)

const (
    testAPIPath = "/api/v1"
)

var (
    testUser = users.NewUser("TestUser", "testpassword", "test@test.com")
    testAdmin = users.NewUser("TestAdmin", "testpassword", "test@test.com")
    testRouter *routes.Router   
)

func TestMain(m *testing.M) {
    testRouter = createTestRouter(testAPIPath)
    retCode := m.Run()
    os.Exit(retCode)
}

func TestGetUser_AsVisitor(t *testing.T) {
    
    requestPath := paths.Join(testAPIPath, "/users/1")
    request := createTestRequest("GET", requestPath, "")
    writer := httptest.NewRecorder()

    testRouter.ServeHTTP(writer, request)

    if status := writer.Code; status != http.StatusUnauthorized {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := string(toJson(Error{
        Message : "Current user does not have permission to access this resource.",
        Code : CodeOwnerOnly,
    }))

    if writer.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            writer.Body.String(), expected)
    }
}

func TestGetUser_AsAdmin(t *testing.T) {
    // TODO: need to implement admin role
}

func TestGetUser_AsDifferentUser(t *testing.T) {
    requestPath := paths.Join(testAPIPath, "/users/1")
    request := createTestRequest("GET", requestPath, "")

    encodedAuthDetails := base64.StdEncoding.EncodeToString([]byte("OtherUser:otherpassword"))
    request.Header["Authorization"] = []string{"Basic " + encodedAuthDetails}

    writer := httptest.NewRecorder()
    testRouter.ServeHTTP(writer, request)

    if status := writer.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := string(toJson(User{
        Username : "TestUser",
        Email : "test@test.com",
    }))

    if writer.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            writer.Body.String(), expected)
    } 
}

func TestGetUser_AsMatchingUser(t *testing.T) {
    requestPath := paths.Join(testAPIPath, "/users/1")
    request := createTestRequest("GET", requestPath, "")

    encodedAuthDetails := base64.StdEncoding.EncodeToString([]byte("TestUser:testpassword"))
    request.Header["Authorization"] = []string{"Basic " + encodedAuthDetails}

    writer := httptest.NewRecorder()
    testRouter.ServeHTTP(writer, request)

    if writer.Code != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := string(toJson(User{
        Username : "TestUser",
        Email : "test@test.com",
    }))

    if writer.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            writer.Body.String(), expected)
    }
}

func TestGetUser_UserDoesNotExist(t *testing.T) {

}

func createTestRouter(apiPath string) *routes.Router {
    ds, err := datastore.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
    if err != nil {
        panic(err)
    }

    migrationsPath := config.GetString(keys.DatabaseMigrationsPath)
    if migrationsPath == "" {
        panic("Migrations path has not been defined.")
    }

    errs, ok := ds.UpSync(migrationsPath)
    if !ok {
        errMsg := ""
        for _, err := range errs {
            errMsg = errMsg + err.Error()
        }

        panic(errMsg)
    }
    
    dao := users.NewDao(ds)
    
    err = dao.Insert(testUser)
    if err != nil {
        panic(err)
    }

    return NewRouter(apiPath, dao) 
}

func createTestRequest(method string, path string, body string) *http.Request {
    reader := strings.NewReader(body) // Empty body for all requests
    request, _ := http.NewRequest(method, path, reader)
    return request
}