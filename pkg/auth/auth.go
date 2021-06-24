package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// JWTExpiration describes how long the token should be valid, in hours
	JWTExpiration = time.Duration(time.Hour * 48)
)

// NewJWT creates a new JWT for the specfied user, based on the
// provided secret
func NewJWT(user string, expTime int64, hexKey string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	context := jwt.MapClaims{}
	userContext := jwt.MapClaims{}
	userContext["name"] = user
	context["user"] = user

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["aud"] = "jitsi"
	claims["iss"] = "openhouse_client"
	claims["sub"] = "localhost:3000"
	claims["room"] = "*"
	claims["exp"] = expTime
	claims["nbf"] = time.Date(2021, 01, 01, 12, 0, 0, 0, time.UTC).Unix()
	claims["context"] = context

	t, err := token.SignedString([]byte(hexKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// GetUserFromJWT checks if the auth token is valid
func GetUserFromJWT(token interface{}) string {
	user := token.(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["name"].(string)
}
