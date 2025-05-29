package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	_ "os"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	//jwtSecret       = []byte(os.Getenv("JWT_SECRET")) // panic early if empty
	jwtSecret = []byte("safojsaidfjsdfjsodfasfodij") // panic early if empty
)

// Claims carries the *minimum* you need in every request
type Claims struct {
	Plaintext string `json:"token"`
	UserID    int    `json:"uid"`
	Hash      []byte `json:"-"`
	Role      string `json:"role"`
	Country   string `json:"country"`
	jwt.RegisteredClaims
}

// GenerateToken signs a new 24-hour token
func GenerateToken(userID int, role, country string) (*Claims, error) {
	claims := Claims{
		UserID:  userID,
		Role:    role,
		Country: country,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	claims.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum256([]byte(claims.Plaintext))
	claims.Hash = hash[:]
	return &claims, nil

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
