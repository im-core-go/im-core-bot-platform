package auth

import "github.com/golang-jwt/jwt"

type UserClaim struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}
