package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	// "go/token"
	"os"
	"project_management/api/constants"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
}

func GenerateJwtToken(userID string, audience string) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			Audience:  audience,
			IssuedAt:  time.Now().UTC().Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(os.Getenv("SecretKey")))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// niche na function ma always accesstoken avse request header mathi.
func VerifyToken(ctx context.Context, tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SecretKey")), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if !token.Valid {
		fmt.Println(constants.INVALID_TOKEN)
		return "", errors.New(constants.INVALID_TOKEN)
	}

	userId := claims.Subject

	return userId, nil
}
