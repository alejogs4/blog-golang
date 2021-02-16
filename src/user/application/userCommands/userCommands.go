package usercommands

import (
	"github.com/alejogs4/blog/src/user/domain/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCommands struct {
	userRepository user.UserRepository
}

func NewUserCommands(userRepository user.UserRepository) UserCommands {
	return UserCommands{userRepository}
}

func (userCommand UserCommands) Login(email, password string) (user.UserDTO, error) {
	loggedUser, err := userCommand.userRepository.Login(email, password)
	if err != nil {
		return user.UserDTO{}, err
	}

	return loggedUser, nil
}

func (userCommands UserCommands) Register(email, password, firstname, lastname string) (user.UserDTO, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return user.UserDTO{}, err
	}

	userID := uuid.New().String()
	newUser, err := user.NewUser(userID, firstname, lastname, email, string(encryptedPassword), false)

	if err != nil {
		return user.UserDTO{}, err
	}

	err = userCommands.userRepository.Register(newUser)
	if err != nil {
		return user.UserDTO{}, err
	}

	return user.ToDTO(newUser), nil
}
