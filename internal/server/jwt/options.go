package jwt

import "time"

type options struct {
	key       []byte
	expiresIn time.Duration
	issuer    string
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

func WithIssuer(issuer string) IssuerOption {
	return func(o *options) {
		o.issuer = issuer
	}
}
