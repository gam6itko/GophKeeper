package jwt

import (
	"context"
	"github.com/gam6itko/goph-keeper/internal/server/jwt"
	"github.com/gam6itko/goph-keeper/internal/server/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

const Header = "Token"

type Interceptor struct {
	issuer *jwt.Issuer
}

func New(issuer *jwt.Issuer) *Interceptor {
	return &Interceptor{
		issuer,
	}
}

func (ths Interceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Login, Registration неу требует аутентификации по Token.
	switch info.Server.(type) {
	case service.AuthServerImpl:
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, `Token not found`)
	}
	token := md.Get(Header)
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, `Token not found`)
	}

	userID, err := ths.issuer.Parse(token[0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, `Failed to extract user ID`)
	}

	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Join(md, metadata.MD{
			"UserID": []string{
				strconv.FormatUint(userID, 10),
			},
		}),
	)

	return handler(ctx, req)
}
