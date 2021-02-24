package userhttpport_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	integrationtest "github.com/alejogs4/blog/src/shared/infraestructure/integrationTest"
	"github.com/alejogs4/blog/src/shared/infraestructure/token"
	"github.com/alejogs4/blog/src/user/domain/user"
	userhttpport "github.com/alejogs4/blog/src/user/infraestructure/userHttpPort"
	_ "github.com/lib/pq"
)

func TestMain(t *testing.M) {
	os.Exit(integrationtest.SetupDatabaseForTesting(t))
}

func TestRegisterLoginIntegration(t *testing.T) {
	t.Run("Should register a new user meanwhile is right data and after allow user login", func(t *testing.T) {
		t.Parallel()

		newUser, _ := user.NewUser("id", "Jose", "Velez", "josevelez99@gmail.com", "123456", false)
		request, response, registerRoute := prepareRegisterRequest(newUser)
		registerRoute(response, request)

		if response.Code != http.StatusCreated {
			t.Errorf("Error: Expected status %d, Received status %d", http.StatusCreated, response.Code)
		}

		var receivedResponse struct {
			Data    user.UserDTO `json:"data"`
			Message string       `json:"message"`
		}

		json.NewDecoder(response.Body).Decode(&receivedResponse)
		if receivedResponse.Data.Email != newUser.GetEmail() {
			t.Errorf("Error: Expected user email %v, Received user email %v", newUser.GetEmail(), receivedResponse.Data.Email)
		}

		loginRequest, loginResponse, loginRoute := prepareLoginRequest(newUser.GetEmail(), newUser.GetPassword())
		loginRoute(loginResponse, loginRequest)

		if loginResponse.Code != http.StatusOK {
			t.Errorf("Error: Expected status %d, Received status %d", http.StatusOK, loginResponse.Code)
		}

		var receivedLoginResponse struct {
			Data    userhttpport.LoginResponse `json:"data"`
			Message string                     `json:"message"`
		}

		json.NewDecoder(loginResponse.Body).Decode(&receivedLoginResponse)
		if receivedLoginResponse.Data.User.Email != newUser.GetEmail() {
			t.Errorf("Error: Expected user email %v, Received user email %v", newUser.GetEmail(), receivedLoginResponse.Data.User.Email)
		}

		_, err := token.ValidateToken(receivedLoginResponse.Data.Token)
		if err != nil {
			t.Errorf("Error: Sended token is invalid - %s", err)
		}
	})

	t.Run("Should send a 400 error if user email or password are wrong", func(t *testing.T) {
		t.Parallel()

		newUser, _ := user.NewUser("id", "Alejandro", "Garcia", "alejogs5@gmail.com", "123456", false)
		request, response, registerRoute := prepareRegisterRequest(newUser)
		registerRoute(response, request)

		if response.Code != http.StatusCreated {
			t.Errorf("Error: Expected status %d, Received status %d", http.StatusCreated, response.Code)
		}

		loginRequest, loginResponse, loginRoute := prepareLoginRequest("wrong-email@gmail.com", newUser.GetPassword())
		loginRoute(loginResponse, loginRequest)

		if loginResponse.Code != http.StatusBadRequest {
			t.Errorf("Error: Expected status %d, Received status %d", http.StatusBadRequest, loginResponse.Code)
		}

		loginRequest, loginResponse, loginRoute = prepareLoginRequest(newUser.GetEmail(), "wrong-password")
		loginRoute(loginResponse, loginRequest)

		if loginResponse.Code != http.StatusBadRequest {
			t.Errorf("Error: Expected status %d, Received status %d", http.StatusBadRequest, loginResponse.Code)
		}
	})
}
