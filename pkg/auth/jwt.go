package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type JwtHandler struct {
	authToken    []byte
	refreshToken []byte
}

func NewJwtHandler() *JwtHandler {
	authToken := os.Getenv("JWT_AUTH_TOKEN")
	if authToken == "" {
		panic("JWT_AUTH_TOKEN env variable not set")
	}
	refreshToken := os.Getenv("JWT_REFRESH_TOKEN")
	if refreshToken == "" {
		panic("JWT_REFRESH_TOKEN env variable not set")
	}
	return &JwtHandler{
		authToken:    []byte(authToken),
		refreshToken: []byte(refreshToken),
	}

}

func (j *JwtHandler) TrackAuthToken(tokenString string, claim jwt.Claims) (jwt.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return j.authToken, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claim, nil
}

func (j *JwtHandler) GenAuthToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(j.authToken)
}
