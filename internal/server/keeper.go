package server

import (
	"context"
	"database/sql"
	"github.com/gam6itko/goph-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KeeperImpl struct {
	proto.UnimplementedKeeperServer
	db *sql.DB
}

func NewKeeperImpl(db *sql.DB) *KeeperImpl {
	return &KeeperImpl{db: db}
}

func (KeeperImpl) List(ctx context.Context, req *proto.ListRequest) (*proto.ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (KeeperImpl) Load(ctx context.Context, req *proto.LoadRequest) (*proto.LoadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Load not implemented")
}
func (KeeperImpl) Store(ctx context.Context, req *proto.StoreRequest) (*proto.StoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
