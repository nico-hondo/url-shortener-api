package pkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	apperror "github.com/nico-hondo/url-shortener/internal/error"
)

type JwtClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func NewJwtClaims(id int) *JwtClaims {
	now := time.Now()
	jti := uuid.NewString()

	return &JwtClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}

func (jc *JwtClaims) GenToken() (string, error) {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		return "", apperror.ErrMissingKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jc)

	//error handling jika kunci tidak ada
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func (jc *JwtClaims) DecodeToken(token string) error {
	jwtToken, err := jwt.ParseWithClaims(token, jc, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return jwt.ErrTokenExpired
	}

	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}

	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
