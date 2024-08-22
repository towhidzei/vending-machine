package user

import (
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Id uint `json:"id"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (m *JWTManager) Generate(u *domain.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().UTC().Add(m.tokenDuration).Unix(),
		},
		Id: u.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Wrapf(domain.ErrInvalidToken, "unexpected token signing method: %v", t.Method)
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.Wrap(domain.ErrInvalidToken, "invalid token claims")
	}

	return claims, nil
}
