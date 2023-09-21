package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID uint // Add additional claims as needed
	jwt.StandardClaims
}
