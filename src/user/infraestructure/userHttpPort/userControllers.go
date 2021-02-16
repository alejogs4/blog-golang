package userhttpport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	usercommands "github.com/alejogs4/blog/src/user/application/userCommands"
	"github.com/alejogs4/blog/src/user/domain/user"
	userrepository "github.com/alejogs4/blog/src/user/infraestructure/userRepository"
)

var userCommand = usercommands.NewUserCommands(userrepository.PostgresUserRepository{})

func loginHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var loginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(request.Body).Decode(&loginInfo)
	if err != nil {
		httputils.DispatchNewHttpError(response, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user, err := userCommand.Login(loginInfo.Email, loginInfo.Password)
	if err != nil {
		httpError := mapUserErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	userResponse := map[string]map[string]interface{}{ // This sucks
		"data": {
			"user":    user,
			"message": "Ok",
			"token":   "",
		},
	}

	httputils.DispatchNewResponse(response, userResponse, http.StatusOK)
}

func registerHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var newUser struct {
		user.UserDTO
		Password string `json:"password"`
	}

	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		httpError := mapUserErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	userDTO, err := userCommand.Register(newUser.Email, newUser.Password, newUser.Firstname, newUser.Lastname)
	if err != nil {
		fmt.Println(err)
		httpError := mapUserErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(userDTO, "Ok"), http.StatusOK)
}
