package jwt

import (
	"context"
	"github.com/gam6itko/goph-keeper/internal/server/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const Header = "Authorization"

type Interceptor struct {
	issuer *jwt.Issuer
}

func New(issuer *jwt.Issuer) *Interceptor {
	return &Interceptor{
		issuer,
	}
}

func (ths Interceptor) Intercept(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, `Token not found`)
	}
	values := md.Get(Header)
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, `Token not found`)
	}

	//todo jwt validate

	return handler(ctx, req)
}
