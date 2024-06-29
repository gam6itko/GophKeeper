package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type options struct {
	key       []byte
	expiresIn time.Duration
}

type IssuerOption func(*options)

func WithKey(key []byte) IssuerOption {
	if key == nil {
		panic("key must not be nil")
	}

	return func(o *options) {
		o.key = key
	}
}

func WithExpiresIn(d time.Duration) IssuerOption {
	return func(o *options) {
		o.expiresIn = d
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

type Issuer struct {
	options options
}

func NewIssuer(opts ...IssuerOption) *Issuer {
	o := options{
		key:       []byte("secret_key"),
		expiresIn: time.Hour,
	}
	for _, opt := range opts {
		opt(&o)
	}
	i := &Issuer{
		options: o,
	}

	return i
}

func (ths Issuer) Issue(userID uint64) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: &jwt.NumericDate{
					Time: time.Now().UTC(),
				},
				ExpiresAt: &jwt.NumericDate{
					Time: time.Now().Add(ths.options.expiresIn).UTC(),
				},
			},
			UserID: userID,
		},
	)

	// создаём строку токена
	tokenString, err = token.SignedString(ths.options.key)
	return
}

func (ths Issuer) Parse(tokenString string) (uint64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return ths.options.key, nil
		},
	)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
