package datamap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type decoder interface {
	Decode(v any) error
}

func NewDataMapFromFile(ctx context.Context, dirFs fs.ReadDirFS, fpath string) (map[string]any, error) {
	f, err := dirFs.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var d decoder

	switch filepath.Ext(fpath) {
	case ".json":
		d = json.NewDecoder(f)
	case ".toml":
		d = toml.NewDecoder(f)
	case ".yaml", ".yml":
		d = yaml.NewDecoder(f)
	default:
		return nil, ErrUnknownFileSource
	}

	var data map[string]any
	done := make(chan error)

	go func() {
		defer close(done)

		done <- d.Decode(&data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("decode file error: %w", err)
		}
	}

	return data, nil
}
