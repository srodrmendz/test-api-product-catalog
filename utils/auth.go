package utils

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/srodrmendz/api-auth/model"
)

// Method used to validate jwt token and obtain claims
func GetClaimsFromToken(token string, secretKey string) (*model.JWTClaim, error) {
	tkn, err := jwt.Parse(token, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	claim, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	email, ok := (claim)["email"].(string)
	if !ok {
		return nil, errors.New("invalid token")
	}

	userName, ok := (claim)["username"].(string)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return &model.JWTClaim{
		Email:    email,
		Username: userName,
	}, nil
}
