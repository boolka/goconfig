//go:build goconfig_vault

package vault

import (
	"context"
	"errors"
	"io/fs"

	"github.com/boolka/goconfig/pkg/datamap"
	vaultApi "github.com/hashicorp/vault/api"
)

type VaultSource struct {
	client *vaultApi.Client
	data   map[string]any
}

func NewVaultSource(ctx context.Context, dirFs fs.ReadDirFS, fpath string, client any) (*VaultSource, error) {
	data, err := datamap.NewDataMapFromFile(ctx, dirFs, fpath)
	if err != nil {
		return nil, err
	}

	vaultClient, ok := client.(*vaultApi.Client)
	if !ok {
		return nil, errors.New("invalid vault client")
	}

	return &VaultSource{
		client: vaultClient,
		data:   data,
	}, nil
}

func (s *VaultSource) Get(ctx context.Context, path string) (any, bool) {
	var d string

	v, ok := datamap.GetByPath(s.data, path)
	if !ok {
		return nil, false
	}

	if d, ok = v.(string); !ok {
		return ErrInvalidPath, false
	}

	vaultMount, vaultPath, mapPath, err := parsePath(d)
	if err != nil {
		return err, false
	}

	secret, err := s.client.KVv2(vaultMount).Get(ctx, vaultPath)
	if err != nil {
		return err, false
	}

	if mapPath == "" {
		mapPath = path
	}

	return datamap.GetByPath(secret.Data, mapPath)
}

func (e *VaultSource) Client() *vaultApi.Client {
	return e.client
}
