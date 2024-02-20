package providers

import (
	"context"
	"errors"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OpenID interface {
	AuthCodeURL(code string, opts... oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, token string, opts... oauth2.AuthCodeOption) (*oauth2.Token, error)
	UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error)
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}

type openID struct {
	*oidc.Provider
	oauth2.Config
}

func NewOpenID() *openID {
	return newOpenIDWithCredentials(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REDIRECT_URL"))
}

func newOpenIDWithCredentials(clientID, clientSecret, redirectURL string) *openID {

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		panic("failed to load oidc provider")
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &openID{
		provider,
		oauth2Config,
	}
}

func (p *openID) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return &oidc.IDToken{}, errors.New("no id token")
	}

	oidConfig := &oidc.Config{
		ClientID: p.ClientID,
	}

	return p.Verifier(oidConfig).Verify(ctx, rawIDToken)
}
