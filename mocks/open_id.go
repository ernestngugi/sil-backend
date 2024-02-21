package mocks

import (
	"context"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type mockOpenID struct {
	AuthToken *oauth2.Token
	User      *oidc.UserInfo
}

func NewMockOpenID() *mockOpenID {
	return &mockOpenID{}
}

func (m *mockOpenID) AuthCodeURL(code string, opts ...oauth2.AuthCodeOption) string {
	return ""
}

func (m *mockOpenID) Exchange(ctx context.Context, token string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return m.AuthToken, nil
}

func (m *mockOpenID) UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error) {
	return m.User, nil
}
