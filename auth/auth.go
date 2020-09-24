package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var tag = "[Switch_Auth]"

// isValidJwtToken checks the structure of the jwt token is present
// this function is necessary as passing a non-jwt string to jwt package panics the whole
// thread which is awful but we have to live by it anyway
func isValidJwtToken(token string) bool {
	return !(len(token) <= 0 || strings.Count(token, ".") != 2)
}

// getjwtToken extracts the jwt token from the request header
func getjwtToken(headers http.Header) string {
	authStringSlices := strings.Split(headers.Get("Authorization"), "bearer ")

	if len(authStringSlices) < 2 {
		return ""
	}
	return authStringSlices[1]
}

// verifyTokenSignature verfies the token passed is a valid jwt token signed by the same secret
func verifyTokenSignature(token string) bool {
	var jwtKey = []byte(jwtSecret)
	claims := &Claims{}

	if !isValidJwtToken(token) {
		fmt.Println(fmt.Sprintf("%s invalid jwt token format, auth rejected", tag))
		return false
	}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !tkn.Valid {
		fmt.Println(fmt.Sprintf("%s invalid token, auth rejected", tag))
		return false
	}

	fmt.Println(fmt.Sprintf("%s token verified for author %s", tag, claims.Author))
	return true
}

// Authenticate verifies the caller comes trusted  authority using jwt
func Authenticate(headers http.Header) bool {
	fmt.Println(fmt.Sprintf("%s Started Authentication Flow ....", tag))
	return verifyTokenSignature(getjwtToken(headers))
}
