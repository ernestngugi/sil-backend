package providers

import (
	"context"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OpenID interface {
	AuthCodeURL(code string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, token string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error)
}

type openID struct {
	*oidc.Provider
	*oauth2.Config
}

func NewOpenID() *openID {
	return newOpenIDWithCredentials(os.Getenv("OAUTH_CLIENT_ID"), os.Getenv("OAUTH_CLIENT_SECRET"), os.Getenv("OAUTH_REDIRECT_URL"))
}

func newOpenIDWithCredentials(clientID, clientSecret, redirectURL string) *openID {

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		panic("failed to load oidc provider")
	}

	oauth2Config := &oauth2.Config{
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
