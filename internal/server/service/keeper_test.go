package service

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gam6itko/goph-keeper/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestKeeperImpl_List(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.
			ExpectQuery("SELECT (.+) FROM `user_data`").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "type", "name", "meta"}).
					AddRow("1", "1", "name", "meta"),
			)

		svc := NewKeeperImpl(db)
		req := proto.ListRequest{}
		ctx := metadata.NewIncomingContext(
			context.TODO(),
			metadata.New(map[string]string{"UserID": "1"}),
		)
		reps, err := svc.List(ctx, &req)
		require.NoError(t, err)
		require.NotNil(t, reps)
		require.Len(t, reps.List, 1)
	})

}
