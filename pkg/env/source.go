package env

import (
	"context"
	"io/fs"
	"os"

	"github.com/boolka/goconfig/pkg/datamap"
)

type EnvSource struct {
	data map[string]any
}

func NewEnvSource(ctx context.Context, dirFs fs.ReadDirFS, fpath string) (*EnvSource, error) {
	data, err := datamap.NewDataMapFromFile(ctx, dirFs, fpath)
	if err != nil {
		return nil, err
	}

	return &EnvSource{
		data: data,
	}, nil
}

func (s *EnvSource) Get(_ context.Context, path string) (any, bool) {
	v, ok := datamap.GetByPath(s.data, path)
	if !ok {
		return nil, false
	}

	vString, ok := v.(string)
	if !ok {
		return nil, false
	}

	return os.LookupEnv(vString)
}
