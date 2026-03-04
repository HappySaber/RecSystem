package middleware

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	//AppID int    `json:"app_id"`
	jwt.RegisteredClaims
}
