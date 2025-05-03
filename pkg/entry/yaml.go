package entry

import (
	"context"
	"fmt"
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
	done := make(chan error)

	go func() {
		defer close(done)

		done <- yaml.NewDecoder(r).Decode(&data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("NewYaml err: %w", err)
		}
	}

	return &YamlEntry{
		data: data,
	}, nil
}

func (e *YamlEntry) Get(ctx context.Context, path string) (any, bool) {
	return getFromMap(e.data, path)
}
