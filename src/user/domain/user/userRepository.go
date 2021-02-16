package user

type UserRepository interface {
	Login(email, password string) (UserDTO, error)
	Register(user User) error
}
