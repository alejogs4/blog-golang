package userhttpport_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
	"github.com/alejogs4/blog/src/user/domain/user"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	userrepository "github.com/alejogs4/blog/src/user/infraestructure/userRepository"
)

func prepareRegisterRequest(newUser user.User) (*http.Request, *httptest.ResponseRecorder, http.HandlerFunc) {
	userController := userhttpport.NewUserController(userrepository.NewUserRepository(testDatabase))

	userBody := []byte(fmt.Sprintf(
		`{"email": "%v", "firstname": "%v", "lastname": "%v", "password": "%v"}`,
		newUser.GetEmail(), newUser.GetFirstname(), newUser.GetLastname(), newUser.GetPassword(),
	))

	request := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(userBody))
	response := httptest.NewRecorder()

	registerRoute := middleware.Chain(userController.RegisterHandler, httputils.Verb(http.MethodPost))
	return request, response, registerRoute
}
