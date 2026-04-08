package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gam6itko/goph-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

type KeeperImpl struct {
	proto.UnimplementedKeeperServer
	db *sql.DB
}

func NewKeeperImpl(db *sql.DB) *KeeperImpl {
	return &KeeperImpl{db: db}
}

func (ths KeeperImpl) List(ctx context.Context, req *proto.ListRequest) (*proto.ListResponse, error) {
	userID, ok := ths.extractUserID(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "metadata parse fail")
	}

	rows, err := ths.db.QueryContext(
		ctx,
		"SELECT `id`,  `type`, `name`, `meta` FROM `user_data` WHERE `user_id` = ?",
		userID,
	)
	if err != nil {
		log.Printf("ERROR. %s", err)
		errMessage := fmt.Sprintf("query rows error: %s", err)
		return nil, status.Error(codes.Internal, errMessage)
	}
	if rows.Err() != nil {
		log.Printf("ERROR. %s", rows.Err())
		return nil, status.Error(codes.Internal, rows.Err().Error())
	}

	items := make([]*proto.PrivateData, 0)
	for rows.Next() {
		var (
			id   uint32
			tp   uint32
			name string
			meta string
		)
		if err = rows.Scan(&id, &tp, &name, &meta); err != nil {
			log.Printf("ERROR. %s", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		items = append(items, &proto.PrivateData{
			Id:   id,
			Type: proto.DataType(tp),
			Name: name,
			Meta: meta,
		})
	}

	return &proto.ListResponse{List: items}, nil
}
func (ths KeeperImpl) Load(ctx context.Context, req *proto.LoadRequest) (*proto.LoadResponse, error) {
	userID, ok := ths.extractUserID(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "metadata parse fail")
	}

	row := ths.db.QueryRowContext(
		ctx,
		"SELECT `id`,  `type`, `name`, `meta`, `data` FROM `user_data` WHERE `user_id` = ? AND `id` = ?",
		userID, req.Id,
	)
	if row.Err() != nil {
		log.Printf("ERROR. %s", row.Err())
		errMessage := fmt.Sprintf("query row error: %s", row.Err())
		return nil, status.Error(codes.Internal, errMessage)
	}

	var (
		id   uint32
		tp   uint32
		name string
		meta string
		data []byte
	)
	if err := row.Scan(&id, &tp, &name, &meta, &data); err != nil {
		log.Printf("ERROR. %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.LoadResponse{
		Item: &proto.PrivateData{
			Id:   id,
			Type: proto.DataType(tp),
			Name: name,
			Meta: meta,
			Data: data,
		},
	}, nil
}

// Store - сохраняет данные в БД.
// В данный момент можно сохранить множество записей с одним и мем же именем.
func (ths KeeperImpl) Store(ctx context.Context, req *proto.StoreRequest) (*proto.StoreResponse, error) {
	userID, ok := ths.extractUserID(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "metadata parse fail")
	}

	r, err := ths.db.ExecContext(
		ctx,
		"INSERT INTO `user_data` (`user_id`, `type`, `name`, `meta`, `data`) VALUE (?,?,?,?,?)",
		userID,
		req.Item.Type,
		req.Item.Name,
		req.Item.Meta,
		req.Item.Data,
	)
	if err != nil {
		log.Printf("ERROR. %s", err)
		errMessage := fmt.Sprintf("query row error: %s", err)
		return nil, status.Error(codes.Internal, errMessage)
	}

	id, err := r.LastInsertId()
	if err != nil {
		log.Printf("ERROR. %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.StoreResponse{
		Id: uint32(id),
	}, nil
}

func (ths KeeperImpl) extractUserID(ctx context.Context) (uint64, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, false
	}

	userIDList := md.Get("UserID")
	if len(userIDList) == 0 {
		return 0, false
	}

	userID, err := strconv.ParseUint(userIDList[0], 10, 32)
	if err != nil {
		log.Printf("ERROR. %s", err)
		return 0, false
	}

	return userID, true
}
