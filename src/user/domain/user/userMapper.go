package user

// ToDTO map a raw user into a UserDTO
func ToDTO(rawUser User) UserDTO {
	return UserDTO{
		ID:            rawUser.GetID(),
		Firstname:     rawUser.GetFirstname(),
		Lastname:      rawUser.GetLastname(),
		Email:         rawUser.GetEmail(),
		EmailVerified: rawUser.GetEmailVerified(),
	}
}
