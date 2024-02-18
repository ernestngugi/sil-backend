package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/ernestngugi/sil-backend/providers"
	"golang.org/x/oauth2"
)

const tokenHeader = "X-SIL-TOKEN"

type Authenticator interface {
	TokenFromRequest(ctx context.Context, request *http.Request) (*oidc.UserInfo, error)
}

type authenticator struct {
	oidcProvider providers.OpenID
}

func NewAuthenticator(oidcProvider providers.OpenID) Authenticator {
	return &authenticator{
		oidcProvider: oidcProvider,
	}
}

func (a *authenticator) TokenFromRequest(ctx context.Context, request *http.Request) (*oidc.UserInfo, error) {

	token := request.Header.Get(tokenHeader)

	if token == "" {
		return &oidc.UserInfo{}, errors.New("token not provided")
	}

	authToken, err := a.oidcProvider.Exchange(ctx, token)
	if err != nil {
		return &oidc.UserInfo{}, err
	}

	_, err = a.oidcProvider.VerifyIDToken(ctx, authToken)
	if err != nil {
		return &oidc.UserInfo{}, err
	}

	userInfo, err := a.oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(authToken))
	if err != nil {
		return &oidc.UserInfo{}, errors.New("failed to get user info")
	}

	return userInfo, nil
}
