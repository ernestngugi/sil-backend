package mocks

type mockOpenID struct{}

func NewMockOpenID() *mockOpenID {
	return &mockOpenID{}
}
