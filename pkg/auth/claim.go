package auth

import "github.com/golang-jwt/jwt"

type UserClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}
