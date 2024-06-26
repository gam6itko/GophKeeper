package server

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"log"
	"time"
)

var _ IServer = (*MockServer)(nil)

type MockServer struct{}

func (m *MockServer) Login(ctx context.Context, dto LoginDTO) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *MockServer) Logout(ctx context.Context) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *MockServer) Registration(ctx context.Context, dto RegistrationDTO) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *MockServer) List(_ context.Context) ([]PrivateDataListItemDTO, error) {
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
			Name: "TextDTO",
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
	enc := encrypt.NewAESCrypt()

	key := enc.KeyFit([]byte("123"))
	sign := "WTF"

	buff := bytes.NewBuffer([]byte("WTF"))
	encoder := gob.NewEncoder(buff)

	var result *PrivateDataDTO
	switch id {
	case 1:
		dto := LoginPassDTO{Sign: sign, Login: "user1", Password: "pass1"}
		err := encoder.Encode(dto)
		if err != nil {
			log.Fatal(err)
		}
		//b := make([]byte, 0, len(signTrait)+len(payload))
		//b = append(b, signTrait...)
		//b = append(b, payload...)
		data, err := enc.Encrypt(key, buff.Bytes())
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
		dto := TextDTO{Sign: sign, Text: "you shouldn't see this"}
		err := encoder.Encode(dto)
		if err != nil {
			log.Fatal(err)
		}
		data, err := enc.Encrypt(key, buff.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		result = &PrivateDataDTO{
			BasePrivateDataDTO: BasePrivateDataDTO{
				ID:   2,
				Name: "TextDTO",
				Type: TypeText,
				Meta: "this is text",
			},
			Data: data,
		}
	case 3:
		dto := BinaryDTO{Sign: sign, Data: []byte("123qweasdzxc_binary")}
		err := encoder.Encode(dto)
		if err != nil {
			log.Fatal(err)
		}
		data, err := enc.Encrypt(key, buff.Bytes())
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
		dto := BankCardDTO{
			Sign:    sign,
			Number:  "4929175965841786",
			Expires: "02/22",
			CVV:     "123",
		}
		err := encoder.Encode(dto)
		if err != nil {
			log.Fatal(err)
		}
		//b := make([]byte, 0, len(signTrait)+len(payload))
		//b = append(b, signTrait...)
		//b = append(b, payload...)
		data, err := enc.Encrypt(key, buff.Bytes())
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

	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case <-t.C:
		return result, nil
	}
}
func (m *MockServer) Store(ctx context.Context, dto PrivateDataDTO) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}
