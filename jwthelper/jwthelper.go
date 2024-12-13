package jwthelper

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTService defines the structure for the JWT service
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new instance of the JWTService with a given secret key
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken generates a new JWT token with the given user information and expiration time (in hours)
func (s *JWTService) GenerateToken(username string, isAdmin bool, expirationHours int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":  username,
		"admin": isAdmin,
		"exp":   time.Now().Add(time.Duration(expirationHours) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		log.Printf("Error signing the token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates the JWT token
func (s *JWTService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		log.Printf("Error parsing the token: %v", err)
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return map[string]interface{}{
			"user":  claims["user"],
			"admin": claims["admin"],
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
