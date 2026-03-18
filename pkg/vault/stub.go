//go:build !goconfig_vault

package vault

import (
	"context"
	"io/fs"

	"github.com/boolka/goconfig/pkg/datamap"
)

type VaultSource struct {
	data map[string]any
}

func NewVaultSource(ctx context.Context, dirFs fs.ReadDirFS, fpath string, _ any) (*VaultSource, error) {
	data, err := datamap.NewDataMapFromFile(ctx, dirFs, fpath)
	if err != nil {
		return nil, err
	}

	return &VaultSource{
		data: data,
	}, nil
}

func (s *VaultSource) Get(_ context.Context, path string) (any, bool) {
	return datamap.GetByPath(s.data, path)
}

func (e *VaultSource) Client() any {
	return nil
}
