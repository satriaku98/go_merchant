package auth

import "github.com/dgrijalva/jwt-go"

type MyClaims struct {
	jwt.StandardClaims
	CustomerId string `json:"customer_id"`
	Name       string `json:"name"`
}
