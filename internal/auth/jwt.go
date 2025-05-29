package auth

import (
	"errors"
	_ "os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	//jwtSecret       = []byte(os.Getenv("JWT_SECRET")) // panic early if empty
	jwtSecret = []byte("safojsaidfjsdfjsodfasfodij") // panic early if empty
)

// Claims carries the *minimum* you need in every request
type Claims struct {
	UserID  int    `json:"uid"`
	Role    string `json:"role"`
	Country string `json:"country"`
	jwt.RegisteredClaims
}

// GenerateToken signs a new 24-hour token
func GenerateToken(userID int, role, country string) (string, error) {
	claims := Claims{
		UserID:  userID,
		Role:    role,
		Country: country,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// ParseToken verifies signature *and* expiry, returning custom claims
func ParseToken(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims,
		func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if err != nil || !token.Valid {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}
