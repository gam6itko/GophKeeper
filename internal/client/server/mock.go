package server

import (
	"context"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"log"
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
		BasePrivateDataDTO: BasePrivateDataDTO{
			ID:   1,
			Name: "LoginPass",
			Type: TypeLoginPass,
			Meta: "this is login pass",
		},
	}
	text := PrivateDataListItemDTO{
		BasePrivateDataDTO: BasePrivateDataDTO{
			ID:   2,
			Name: "Text",
			Type: TypeText,
			Meta: "this is text",
		},
	}
	binary := PrivateDataListItemDTO{
		BasePrivateDataDTO: BasePrivateDataDTO{
			ID:   3,
			Name: "Binary",
			Type: TypeBinary,
			Meta: "this is binary",
		},
	}
	bankCard := PrivateDataListItemDTO{
		BasePrivateDataDTO: BasePrivateDataDTO{
			ID:   4,
			Name: "BankCard",
			Type: TypeBankCard,
			Meta: "this is bank card",
		},
	}

	return []PrivateDataListItemDTO{
		loginPass,
		text,
		binary,
		bankCard,
	}, nil
}

func (m *MockServer) Load(ctx context.Context, id uint32) (*PrivateDataDTO, error) {
	c := encrypt.NewAESCrypt()

	key := c.KeyFit([]byte("123"))
	sign := "WTF"

	var result *PrivateDataDTO
	switch id {
	case 1:
		data, err := c.Encrypt(key, []byte(sign+"_login_password"))
		if err != nil {
			log.Fatal(err)
		}
		result = &PrivateDataDTO{
			BasePrivateDataDTO: BasePrivateDataDTO{
				ID:   1,
				Name: "LoginPass",
				Type: TypeLoginPass,
				Meta: "this is login pass",
			},
			Data: data,
		}
	case 2:
		data, err := c.Encrypt(key, []byte(sign+"_text"))
		if err != nil {
			log.Fatal(err)
		}
		result = &PrivateDataDTO{
			BasePrivateDataDTO: BasePrivateDataDTO{
				ID:   2,
				Name: "Text",
				Type: TypeText,
				Meta: "this is text",
			},
			Data: data,
		}
	case 3:
		data, err := c.Encrypt(key, []byte(sign+"_binary"))
		if err != nil {
			log.Fatal(err)
		}
		result = &PrivateDataDTO{
			BasePrivateDataDTO: BasePrivateDataDTO{
				ID:   3,
				Name: "Binary",
				Type: TypeBinary,
				Meta: "this is binary",
			},
			Data: data,
		}
	case 4:
		data, err := c.Encrypt(key, []byte(sign+"_bank_card"))
		if err != nil {
			log.Fatal(err)
		}
		result = &PrivateDataDTO{
			BasePrivateDataDTO: BasePrivateDataDTO{
				ID:   4,
				Name: "BankCard",
				Type: TypeBankCard,
				Meta: "this is bank card",
			},
			Data: data,
		}
	}

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
