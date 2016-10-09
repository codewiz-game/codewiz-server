package routes

import (
	"github.com/crob1140/codewiz/models/users"
	"net/http"
	"net/http/httptest"
	paths "path"
	"strings"
	"testing"
)

const testRouterPath = "/"
const testRoutePath = "test/route"


func TestMiddlewareUsedToPreventHandler(t *testing.T) {
	router := createTestRouter()

	var middlewareReached bool
	filterMiddleware := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
		middlewareReached = true
	})
	router.Use(filterMiddleware)

	var handlerReached bool
	handler := func(w http.ResponseWriter, r *http.Request, context *Context) {
		handlerReached = true
	}

	route := createTestRoute(router)
	route.HandlerFunc(handler)

	writer := httptest.NewRecorder()
	request := createTestRequest()
	router.ServeHTTP(writer, request)

	if !middlewareReached {
		t.Errorf("Filter middleware was not called")
	}

	if handlerReached {
		t.Errorf("Handler was still called despite filter middleware")
	}
}

func TestMiddlewareShareSameContextInstance(t *testing.T) {
	router := createTestRouter()

	testUser := users.NewUser("TestUsername", "testpassword", "test@test.com")
	authMiddleware := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
		context.User = testUser
		next(w, r, context)
	})
	router.Use(authMiddleware)

	var userMatches bool
	handler := func(w http.ResponseWriter, r *http.Request, context *Context) {
		userMatches = (context.User == testUser)
	}

	route := createTestRoute(router)
	route.HandlerFunc(handler)

	writer := httptest.NewRecorder()
	request := createTestRequest()
	router.ServeHTTP(writer, request)

	if !userMatches {
		t.Errorf("Context was not passed as expected from middleware to handler")
	}
}

func TestRouterCreatedWithNoMiddleware(t *testing.T) {
	router := createTestRouter()

	var handlerReached bool
	handler := func(w http.ResponseWriter, r *http.Request, context *Context) {
		handlerReached = true
	}

	route := createTestRoute(router)
	route.HandlerFunc(handler)

	writer := httptest.NewRecorder()
	request := createTestRequest()
	router.ServeHTTP(writer, request)

	if !handlerReached {
		t.Errorf("Handler could not be reached without middleware")
	}
}

func TestReplacingHandlerDoesNotAffectMiddleware(t *testing.T) {
	router := createTestRouter()

	var middlewareReached bool
	middleware := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
		middlewareReached = true
		next(w, r, context)
	})
	router.Use(middleware)

	var firstHandlerReached bool
	firstHandler := func(w http.ResponseWriter, r *http.Request, context *Context) {
		firstHandlerReached = true
	}

	var secondHandlerReached bool
	secondHandler := func(w http.ResponseWriter, r *http.Request, context *Context) {
		secondHandlerReached = true
	}

	route := createTestRoute(router)
	route.HandlerFunc(firstHandler)
	route.HandlerFunc(secondHandler)

	writer := httptest.NewRecorder()
	request := createTestRequest()
	router.ServeHTTP(writer, request)

	if !middlewareReached {
		t.Errorf("Middleware was not called")
	}

	if firstHandlerReached {
		t.Errorf("First handler was still called")
	}

	if !secondHandlerReached {
		t.Errorf("Second handler was not called")
	}
}

func TestMiddlewareExecuteInOrderOfAddition(t *testing.T) {
	router := createTestRouter()
	epoch := 0

	firstMiddlewareEpoch := -1
	firstMiddleware := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
		firstMiddlewareEpoch = epoch
		epoch++
		next(w, r, context)
	})
	router.Use(firstMiddleware)

	secondMiddlewareEpoch := -1
	secondMiddleware := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, context *Context, next HandlerFunc) {
		secondMiddlewareEpoch = epoch
		epoch++
		next(w, r, context)
	})
	router.Use(secondMiddleware)


	handlerEpoch := -1
	handler := func(w http.ResponseWriter, r *http.Request, context *Context) {
		handlerEpoch = epoch
		epoch++
	}

	route := createTestRoute(router)
	route.HandlerFunc(handler)

	writer := httptest.NewRecorder()
	request := createTestRequest()
	router.ServeHTTP(writer, request)

	if firstMiddlewareEpoch == -1 {
		t.Errorf("First middleware was not called")
	}

	if secondMiddlewareEpoch == -1 {
		t.Errorf("Second middleware was not called")
	}

	if handlerEpoch == -1 {
		t.Errorf("Handler was not called")
	}

	if secondMiddlewareEpoch < firstMiddlewareEpoch {
		t.Errorf("First middleware was executed before second middleware")
	}

	if handlerEpoch < secondMiddlewareEpoch {
		t.Errorf("Handler was executed before second middleware")
	}
}

func createTestRouter() *Router {
	return NewRouter(testRouterPath)
}

func createTestRoute(router *Router) *Route {
	return router.Path(testRoutePath).Methods("GET")
}

func createTestRequest() *http.Request {
	requestPath := paths.Join(testRouterPath, testRoutePath)
	requestReader := strings.NewReader("") // Empty body for all requests
	request, _ := http.NewRequest("GET", requestPath, requestReader)
	return request
}
