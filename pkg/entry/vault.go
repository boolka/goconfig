package entry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

var ErrVaultInvalidPath = errors.New("invalid vault path")
var ErrVaultUnauthorized = errors.New("vault unauthorized")

func parsePath(cfgPath string) (string, string, string, error) {
	sepPath := strings.Split(cfgPath, ",")

	switch len(sepPath) {
	case 2:
		return sepPath[0], sepPath[1], "", nil
	case 3:
		return sepPath[0], sepPath[1], sepPath[2], nil
	}

	return "", "", "", ErrVaultInvalidPath
}

type VaultEntry struct {
	client *vault.Client
	entry  Entry
}

func NewVault(ctx context.Context, entry Entry, client *vault.Client, auth vault.AuthMethod) (*VaultEntry, error) {
	if client.Token() == "" {
		if auth == nil {
			return nil, ErrVaultUnauthorized
		}

		r, err := client.Auth().Login(ctx, auth)
		if err != nil {
			return nil, fmt.Errorf("vault auth error: %w", err)
		}
		if r == nil {
			return nil, errors.New("no auth info was returned after login")
		}
	}

	return &VaultEntry{
		client: client,
		entry:  entry,
	}, nil
}

func (e *VaultEntry) Get(ctx context.Context, path string) (any, bool) {
	var data string

	v, ok := e.entry.Get(ctx, path)

	if !ok {
		return nil, false
	}

	if data, ok = v.(string); !ok {
		return ErrVaultInvalidPath, false
	}

	vaultMount, vaultPath, mapPath, err := parsePath(data)
	if err != nil {
		return err, false
	}

	secret, err := e.client.KVv2(vaultMount).Get(ctx, vaultPath)
	if err != nil {
		return err, false
	}

	if mapPath == "" {
		mapPath = path
	}

	return getFromMap(secret.Data, mapPath)
}

func (e *VaultEntry) Client() *vault.Client {
	return e.client
}
