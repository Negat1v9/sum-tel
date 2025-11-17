package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNotValid = errors.New("not valid data")
)

const (
	expirationTime = time.Hour * 24
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// Generate new JWT token, expTime added in side
func GenerateJwtToken(claims *Claims, secretKey []byte) (string, error) {

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}

// Checking the token and returning userID from it
func JwtClaimsFromToken(tokenStr string, secretKey []byte) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	switch {
	case err != nil:
		return nil, ErrNotValid

	case !token.Valid:
		return nil, ErrNotValid
	}

	return claims, nil
}
