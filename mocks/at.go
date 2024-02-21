package mocks

import "github.com/ernestngugi/sil-backend/internal/model"

type mockATProvider struct{}

func NewMockATProvider() *mockATProvider {
	return &mockATProvider{}
}

func (m *mockATProvider) Send(request *model.ATRequest) error {
	return nil
}
