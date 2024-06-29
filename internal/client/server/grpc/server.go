package grpc

import (
	"context"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/proto"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	auth   proto.AuthClient
	keeper proto.KeeperClient

	token string
}

func New(auth proto.AuthClient, keeper proto.KeeperClient) *Server {
	return &Server{
		auth:   auth,
		keeper: keeper,
	}
}

func (ths Server) Registration(ctx context.Context, dto server.RegistrationDTO) error {
	req := &proto.RegistrationRequest{
		Username: dto.Username,
		Password: dto.Password,
	}

	_, err := ths.auth.Registration(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (ths Server) Login(ctx context.Context, dto server.LoginDTO) error {
	req := &proto.LoginRequest{
		Username: dto.Username,
		Password: dto.Password,
	}

	resp, err := ths.auth.Login(ctx, req)
	if err != nil {
		return err
	}

	ths.token = resp.Token

	return nil
}

func (ths Server) Logout(_ context.Context) error {
	ths.token = ""
	return nil
}

func (ths Server) List(ctx context.Context) ([]server.PrivateDataListItemDTO, error) {
	resp, err := ths.keeper.List(ths.withToken(ctx), &proto.ListRequest{})
	if err != nil {
		return nil, err
	}

	result := make([]server.PrivateDataListItemDTO, len(resp.List))
	for i, item := range resp.List {
		result[i] = server.PrivateDataListItemDTO{
			BasePrivateDataDTO: server.BasePrivateDataDTO{
				ID:   item.Id,
				Type: server.PrivateDataType(item.Type),
				Name: item.Name,
				Meta: item.Meta,
			},
		}
	}

	return result, err
}

func (ths Server) Load(ctx context.Context, id uint32) (*server.PrivateDataDTO, error) {
	resp, err := ths.keeper.Load(ths.withToken(ctx), &proto.LoadRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &server.PrivateDataDTO{
		BasePrivateDataDTO: server.BasePrivateDataDTO{
			ID:   resp.Item.Id,
			Type: server.PrivateDataType(resp.Item.Type),
			Name: resp.Item.Name,
			Meta: resp.Item.Meta,
		},
		Data: resp.Item.Data,
	}, nil
}

func (ths Server) Store(ctx context.Context, dto server.PrivateDataDTO) error {
	req := &proto.StoreRequest{
		Item: &proto.PrivateData{
			Id:   dto.ID,
			Type: proto.DataType(dto.Type),
			Name: dto.Name,
			Meta: dto.Meta,
		},
	}
	_, err := ths.keeper.Store(ths.withToken(ctx), req)
	if err != nil {
		return err
	}

	return nil
}

// withToken добавляет JWT в метаданные запроса.
func (ths Server) withToken(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(
		ctx,
		metadata.New(map[string]string{
			"Token": ths.token,
		}),
	)
}
