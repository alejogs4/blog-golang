package token

import (
	"errors"
	"time"

	"github.com/alejogs4/blog/src/user/domain/user"
	"github.com/dgrijalva/jwt-go"
)

type claimer struct {
	user.UserDTO
	jwt.StandardClaims
}

func GenerateToken(payload user.UserDTO) (string, error) {
	claimer := claimer{
		UserDTO: user.UserDTO{
			ID:            payload.ID,
			Firstname:     payload.Firstname,
			Lastname:      payload.Lastname,
			Email:         payload.Email,
			EmailVerified: payload.EmailVerified,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Alejandro",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claimer)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(token string) (user.UserDTO, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &claimer{}, verifyFunction)
	if err != nil {
		return user.UserDTO{}, err
	}

	if !parsedToken.Valid {
		return user.UserDTO{}, user.ErrInvalidUser
	}

	claim, ok := parsedToken.Claims.(*claimer)
	if !ok {
		return user.UserDTO{}, errors.New("It was not possible get user from token")
	}

	return claim.UserDTO, nil
}

func verifyFunction(t *jwt.Token) (interface{}, error) {
	return verifyKey, nil
}
