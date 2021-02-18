package user

import (
	"strings"

	fieldutils "github.com/alejogs4/blog/src/shared/domain/fieldsutils"
)

const MinPasswordLength = 6

type User struct {
	id            string
	firstname     string
	lastname      string
	email         string
	emailVerified bool
	password      string
}

func NewUser(id, firstname, lastname, email, password string, emailVerified bool) (User, error) {
	normalizedID := fieldutils.NormalizedStringField(id)
	normalizedFirstname := fieldutils.NormalizedStringField(firstname)
	normalizedLastname := strings.TrimSpace(lastname)
	normalizedEmail := fieldutils.NormalizedStringField(email)

	if normalizedID == "" || normalizedFirstname == "" || normalizedLastname == "" || normalizedEmail == "" || password == "" {
		return User{}, ErrBadUserData
	}

	if len(password) < MinPasswordLength {
		return User{}, ErrTooShortUserPassword
	}

	return User{
		id:            normalizedID,
		firstname:     normalizedFirstname,
		lastname:      normalizedLastname,
		email:         normalizedEmail,
		password:      password,
		emailVerified: emailVerified,
	}, nil
}

// GetID .
func (user User) GetID() string {
	return user.id
}

// GetFirstname .
func (user User) GetFirstname() string {
	return user.firstname
}

// GetLastname .
func (user User) GetLastname() string {
	return user.lastname
}

// GetEmail .
func (user User) GetEmail() string {
	return user.email
}

// GetEmailVerified .
func (user User) GetEmailVerified() bool {
	return user.emailVerified
}

// GetPassword .
func (user User) GetPassword() string {
	return user.password
}
