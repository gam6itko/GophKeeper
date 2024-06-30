package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64 `json:"uid"`
}

type Issuer struct {
	options options
}

func NewIssuer(opts ...IssuerOption) *Issuer {
	o := options{
		key:       []byte("secret_key"),
		expiresIn: time.Hour,
		issuer:    "server",
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
				Issuer: ths.options.issuer,
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

func (ths Issuer) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return ths.options.key, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
