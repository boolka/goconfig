package vault

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

type TokenAuth struct {
	Token string
}

func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{
		Token: token,
	}
}

func (a *TokenAuth) Login(_ context.Context, client *vault.Client) (*vault.Secret, error) {
	return &vault.Secret{
		Auth: &vault.SecretAuth{
			ClientToken: a.Token,
		},
	}, nil
}
