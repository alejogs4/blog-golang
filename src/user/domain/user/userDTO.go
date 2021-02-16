package user

type UserDTO struct {
	ID            string `json:"id,omitempty"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}
