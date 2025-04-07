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

const trimChars = "\t\r\n\x20"

func parsePath(cfgPath string) (mount string, secret string, key string, err error) {
	sepPath := strings.Split(cfgPath, ",")

	switch len(sepPath) {
	case 2:
		return strings.Trim(sepPath[0], trimChars), strings.Trim(sepPath[1], trimChars), "", nil
	case 3:
		return strings.Trim(sepPath[0], trimChars), strings.Trim(sepPath[1], trimChars), strings.Trim(sepPath[2], trimChars), nil
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
