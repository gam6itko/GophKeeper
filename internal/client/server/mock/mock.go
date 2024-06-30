package mock

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"log"
	"time"
)

var _ server.IServer = (*Server)(nil)

// Server структура для разработки и отладки.
type Server struct{}

func (m *Server) Login(ctx context.Context, dto server.LoginDTO) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *Server) Logout(ctx context.Context) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *Server) Registration(ctx context.Context, dto server.RegistrationDTO) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

func (m *Server) List(_ context.Context) ([]server.PrivateDataListItemDTO, error) {
	loginPass := server.PrivateDataListItemDTO{
		BasePrivateDataDTO: server.BasePrivateDataDTO{
			ID:   1,
			Name: "LoginPass",
			Type: server.TypeLoginPass,
			Meta: "this is login pass",
		},
	}
	text := server.PrivateDataListItemDTO{
		BasePrivateDataDTO: server.BasePrivateDataDTO{
			ID:   2,
			Name: "TextDTO",
			Type: server.TypeText,
			Meta: "this is text",
		},
	}
	binary := server.PrivateDataListItemDTO{
		BasePrivateDataDTO: server.BasePrivateDataDTO{
			ID:   3,
			Name: "Binary",
			Type: server.TypeBinary,
			Meta: "this is binary",
		},
	}
	bankCard := server.PrivateDataListItemDTO{
		BasePrivateDataDTO: server.BasePrivateDataDTO{
			ID:   4,
			Name: "BankCard",
			Type: server.TypeBankCard,
			Meta: "this is bank card",
		},
	}

	return []server.PrivateDataListItemDTO{
		loginPass,
		text,
		binary,
		bankCard,
	}, nil
}

func (m *Server) Load(ctx context.Context, id uint32) (*server.PrivateDataDTO, error) {
	enc := encrypt.NewAESCrypt()

	key := enc.KeyFit([]byte("123"))

	buff := bytes.NewBuffer([]byte("WTF"))
	encoder := gob.NewEncoder(buff)

	var result *server.PrivateDataDTO
	switch id {
	case 1:
		dto := server.LoginPassDTO{Login: "user1", Password: "pass1"}
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
		result = &server.PrivateDataDTO{
			BasePrivateDataDTO: server.BasePrivateDataDTO{
				ID:   1,
				Name: "LoginPass",
				Type: server.TypeLoginPass,
				Meta: "this is login pass",
			},
			Data: data,
		}
	case 2:
		dto := server.TextDTO{Text: "you shouldn't see this"}
		err := encoder.Encode(dto)
		if err != nil {
			log.Fatal(err)
		}
		data, err := enc.Encrypt(key, buff.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		result = &server.PrivateDataDTO{
			BasePrivateDataDTO: server.BasePrivateDataDTO{
				ID:   2,
				Name: "TextDTO",
				Type: server.TypeText,
				Meta: "this is text",
			},
			Data: data,
		}
	case 3:
		dto := server.BinaryDTO{Data: []byte("123qweasdzxc_binary")}
		err := encoder.Encode(dto)
		if err != nil {
			log.Fatal(err)
		}
		data, err := enc.Encrypt(key, buff.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		result = &server.PrivateDataDTO{
			BasePrivateDataDTO: server.BasePrivateDataDTO{
				ID:   3,
				Name: "Binary",
				Type: server.TypeBinary,
				Meta: "this is binary",
			},
			Data: data,
		}
	case 4:
		dto := server.BankCardDTO{
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
		result = &server.PrivateDataDTO{
			BasePrivateDataDTO: server.BasePrivateDataDTO{
				ID:   4,
				Name: "BankCard",
				Type: server.TypeBankCard,
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
func (m *Server) Store(ctx context.Context, dto server.PrivateDataDTO) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}
