package entry

import (
	"context"
	"encoding/json"
	"io"
)

type JsonEntry struct {
	data map[string]any
}

func NewJson(ctx context.Context, r io.Reader) (*JsonEntry, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	data := make(map[string]any)
	done := make(chan struct{})
	var err error

	go func() {
		defer close(done)

		err = json.NewDecoder(r).Decode(&data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}

	if err != nil {
		return nil, err
	}

	return &JsonEntry{
		data: data,
	}, nil
}

func (e *JsonEntry) Get(_ context.Context, path string) (any, bool) {
	return getFromMap(e.data, path)
}
