package jwt

import (
	"fmt"
	"time"

	djwt "github.com/dgrijalva/jwt-go"
)

// JWT JWT
type JWT struct {
	secretKey string
}

// Token Token
type Token struct {
	djwt.StandardClaims
}

// New New
func New(secretKey string) *JWT {
	if len(secretKey) == 0 {
		return nil
	}

	return &JWT{
		secretKey: secretKey,
	}
}

// Verify Verify
func (e *JWT) Verify(accessToken string) error {
	token, err := djwt.Parse(accessToken, func(token *djwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*djwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(e.secretKey), nil
	})

	if err != nil || !token.Valid {
		return err
	}

	_, ok := token.Claims.(djwt.MapClaims)

	if !ok {
		return err
	}

	return nil
}

// Generate Generate
func (e *JWT) Generate() (string, error) {
	claims := Token{
		djwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := djwt.NewWithClaims(djwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(e.secretKey))
}
