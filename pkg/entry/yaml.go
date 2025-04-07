package entry

import (
	"context"
	"io"

	"gopkg.in/yaml.v3"
)

type YamlEntry struct {
	data map[string]any
}

func NewYaml(ctx context.Context, r io.Reader) (*YamlEntry, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	data := make(map[string]any)
	done := make(chan struct{})
	var err error

	go func() {
		defer close(done)

		err = yaml.NewDecoder(r).Decode(&data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}

	if err != nil {
		return nil, err
	}

	return &YamlEntry{
		data: data,
	}, nil
}

func (y *YamlEntry) Get(ctx context.Context, path string) (any, bool) {
	return getFromMap(y.data, path)
}
