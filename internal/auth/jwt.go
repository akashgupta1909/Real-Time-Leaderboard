package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT creates a new JWT token with the provided userId and signs it with the secret key.
// It returns the signed token as a string or an error if the signing fails.
func GenerateJWT(userId string) (string, error) {

	jwtExpiration := os.Getenv("JWT_EXPIRATION")
	if jwtExpiration == "" {
		jwtExpiration = "24" // Default expiration time in hours
	}
	parsedExpiration, err := strconv.Atoi(jwtExpiration)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Duration(parsedExpiration) * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateJWT parses and validates the provided JWT token string using the secret key.
// It returns the parsed token and any error encountered during validation.
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
