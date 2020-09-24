package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims the struct holding the expected fields encoded in incoming jwt tokens
type Claims struct {
	Author string `json:"author"`
	jwt.StandardClaims
}
