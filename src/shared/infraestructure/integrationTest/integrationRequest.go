package integrationtest

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	userrepository "github.com/alejogs4/blog/src/user/infraestructure/userRepository"
)

func PrepareLoginRequest(email, password string, testDatabase *sql.DB) (*http.Request, *httptest.ResponseRecorder, http.HandlerFunc) {
	userController := userhttpport.NewUserController(userrepository.NewUserRepository(testDatabase))
	loginBody := []byte(fmt.Sprintf(
		`{"email": "%v", "password": "%v"}`,
		email, password,
	))

	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(loginBody))
	loginResponse := httptest.NewRecorder()
	loginRoute := middleware.Chain(userController.LoginHandler, httputils.Verb(http.MethodPost))

	return loginRequest, loginResponse, loginRoute
}
