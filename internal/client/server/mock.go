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

func (m *MockServer) Logout(ctx context.Context) error {
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

func (m *MockServer) List(ctx context.Context) ([]PrivateDataListItemDTO, error) {
	loginPass := PrivateDataListItemDTO{
		ID:   1,
		Name: "LoginPass",
		Type: TypeLoginPass,
		Meta: "this is login pass",
	}
	text := PrivateDataListItemDTO{
		ID:   2,
		Name: "Text",
		Type: TypeText,
		Meta: "this is text",
	}
	binary := PrivateDataListItemDTO{
		ID:   3,
		Name: "Binary",
		Type: TypeBinary,
		Meta: "this is binary",
	}
	bankCard := PrivateDataListItemDTO{
		ID:   4,
		Name: "BankCard",
		Type: TypeBankCard,
		Meta: "this is bank card",
	}
	return []PrivateDataListItemDTO{
		loginPass,
		text,
		binary,
		bankCard,
	}, nil
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
