package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func SetSecret(secret string) {
	secretKey = []byte(secret)
}

type Claims struct {
	UserID        string   `json:"user_id"`
	Email         string   `json:"email"`
	Name          string   `json:"name,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	KeycloakToken string   `json:"keycloak_token,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, email string, name string, roles []string, keycloakToken ...string) (string, error) {
	var kcToken string
	if len(keycloakToken) > 0 {
		kcToken = keycloakToken[0]
	}

	claims := &Claims{
		UserID:        userID,
		Email:         email,
		Name:          name,
		Roles:         roles,
		KeycloakToken: kcToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
