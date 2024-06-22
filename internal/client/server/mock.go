package server

import (
	"context"
	"time"
)

var _ IServer = (*MockServer)(nil)

type MockServer struct{}

func (m *MockServer) Login(ctx context.Context, dto LoginDTO) error {
	t := time.NewTicker(2 * time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *MockServer) Registration(ctx context.Context, dto RegistrationDTO) error {
	t := time.NewTicker(2 * time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *MockServer) Load(ctx context.Context, id uint32) (*PrivateDataDTO, error) {
	result := &PrivateDataDTO{}

	t := time.NewTicker(2 * time.Second)
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case <-t.C:
		return result, nil
	}
}
func (m *MockServer) Store(ctx context.Context, dto PrivateDataDTO) error {
	t := time.NewTicker(2 * time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}
