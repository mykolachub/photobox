package utils

import (
	"errors"
	"photobox-user/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(id, email, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.JWTClainsUserId:     id,
		config.JWTClainsUserEmail:  email,
		config.JWTClainsExpiration: time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func ParseAndValidateToken(accessToken, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
