package user

// UserRepository defines how the methods that interact with users will be used
type UserRepository interface {
	Login(email, password string) (UserDTO, error)
	Register(user User) error
}
