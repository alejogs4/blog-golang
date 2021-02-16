package userrepository

import (
	"database/sql"
	"errors"

	"github.com/alejogs4/blog/src/shared/infraestructure/database"
	"github.com/alejogs4/blog/src/user/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type PostgresUserRepository struct{}

func (repository PostgresUserRepository) Login(email, password string) (user.UserDTO, error) {
	result := database.PostgresDB.QueryRow(
		"SELECT id, firstname, lastname, email_verified, password FROM person WHERE email=$1",
		email,
	)
	var id string
	var firstname string
	var lastname string
	var emailVerified bool
	var userPassword string

	err := result.Scan(&id, &firstname, &lastname, &emailVerified, &userPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return user.UserDTO{}, user.ErrInvalidUserLogin
	}

	if err != nil {
		return user.UserDTO{}, err
	}

	if passErr := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); passErr != nil {
		return user.UserDTO{}, user.ErrInvalidUserLogin
	}

	domainUser, error := user.NewUser(id, firstname, lastname, email, password, emailVerified)
	if error != nil {
		return user.UserDTO{}, error
	}

	return user.ToDTO(domainUser), nil
}

func (repository PostgresUserRepository) Register(user user.User) error {
	_, err := database.PostgresDB.Exec(
		"INSERT INTO person(id, firstname, lastname, email, email_verified, password) VALUES($1, $2, $3, $4, $5, $6)",
		user.GetID(), user.GetFirstname(), user.GetLastname(), user.GetEmail(), user.GetEmailVerified(), user.GetPassword(),
	)

	return err
}
