package userhttpport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/token"
	usercommands "github.com/alejogs4/blog/src/user/application/userCommands"
	"github.com/alejogs4/blog/src/user/domain/user"
)

type LoginResponse struct {
	User  user.UserDTO `json:"user"`
	Token string       `json:"token"`
}

type userControllers struct {
	userCommand usercommands.UserCommands
}

func NewUserController(userRepository user.UserRepository) userControllers {
	return userControllers{userCommand: usercommands.NewUserCommands(userRepository)}
}

func (controller userControllers) LoginHandler(response http.ResponseWriter, request *http.Request) {
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

	user, err := controller.userCommand.Login(loginInfo.Email, loginInfo.Password)
	if err != nil {
		httpError := mapUserErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	userToken, err := token.GenerateToken(user)
	if err != nil {
		httpError := mapUserErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	userResponse := LoginResponse{User: user, Token: userToken}
	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(userResponse, "Ok"), http.StatusOK)
}

func (controller userControllers) RegisterHandler(response http.ResponseWriter, request *http.Request) {
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

	userDTO, err := controller.userCommand.Register(newUser.Email, newUser.Password, newUser.Firstname, newUser.Lastname)
	if err != nil {
		fmt.Println(err)
		httpError := mapUserErrorToHttpError(err)
		httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
		return
	}

	httputils.DispatchNewResponse(response, httputils.WrapAPIResponse(userDTO, "Ok"), http.StatusCreated)
}
