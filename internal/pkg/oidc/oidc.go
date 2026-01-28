package oidc

import (
	"context"
	"fmt"

	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCProvider struct {
	Provider     *oidc.Provider
	IssuerURL    string
	OAuth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

func NewOIDCProvider(ctx context.Context, cfg config.KeycloakConfig, redirectURL string) (*OIDCProvider, error) {
	issuer := fmt.Sprintf("%s/realms/%s", cfg.URL, cfg.Realm)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &OIDCProvider{
		Provider:     provider,
		IssuerURL:    issuer,
		OAuth2Config: oauth2Config,
		Verifier:     verifier,
	}, nil
}
