package server

import (
	"context"
	"time"
)

type MockLoginHandler struct{}

func (m *MockLoginHandler) Login(ctx context.Context, username string, password string) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}

type MockRegistrationHandler struct{}

func (m *MockRegistrationHandler) Register(ctx context.Context, username string, password string) error {
	t := time.NewTicker(time.Second)
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-t.C:
		return nil
	}
}
