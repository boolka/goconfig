package entry

import (
	"context"
	"os"
)

type EnvEntry struct {
	entry Entry
}

func NewEnv(entry Entry) (*EnvEntry, error) {
	return &EnvEntry{
		entry: entry,
	}, nil
}

func (e *EnvEntry) Get(ctx context.Context, path string) (any, bool) {
	v, ok := e.entry.Get(ctx, path)
	if !ok {
		return nil, false
	}

	vString, ok := v.(string)
	if !ok {
		return nil, false
	}

	return os.LookupEnv(vString)
}
