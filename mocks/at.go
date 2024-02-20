package mocks

type mockATProvider struct{} 

func NewMockATProvider() *mockATProvider {
	return &mockATProvider{}
}

func (m *mockATProvider) Send() error {
	return nil
}
