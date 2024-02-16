package providers

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OpenID struct {
	*oidc.Provider
	oauth2.Config
}

func NewOpenID(clientID, clientSecret, redirectURL string) *OpenID {

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

	return &OpenID{
		provider,
		oauth2Config,
	}
}

func (p *OpenID) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return &oidc.IDToken{}, errors.New("no id token")
	}

	oidConfig := &oidc.Config{
		ClientID: p.ClientID,
	}

	return p.Verifier(oidConfig).Verify(ctx, rawIDToken)
}
