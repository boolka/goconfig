package entry

import (
	"context"
	"fmt"
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
	done := make(chan error)

	go func() {
		defer close(done)

		done <- toml.NewDecoder(r).Decode(&data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("NewToml err: %w", err)
		}
	}

	return &TomlEntry{
		data: data,
	}, nil
}

func (e *TomlEntry) Get(_ context.Context, path string) (any, bool) {
	return getFromMap(e.data, path)
}
