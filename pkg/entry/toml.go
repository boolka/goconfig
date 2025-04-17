package entry

import (
	"context"
	"io"

	"github.com/pelletier/go-toml/v2"
)

type TomlEntry struct {
	data map[string]any
}

func NewToml(ctx context.Context, r io.Reader) (*TomlEntry, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	data := make(map[string]any)
	done := make(chan struct{})
	var err error

	go func() {
		defer close(done)

		err = toml.NewDecoder(r).Decode(&data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}

	if err != nil {
		return nil, err
	}

	return &TomlEntry{
		data: data,
	}, nil
}

func (e *TomlEntry) Get(_ context.Context, path string) (any, bool) {
	return getFromMap(e.data, path)
}
