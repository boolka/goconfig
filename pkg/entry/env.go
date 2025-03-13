package entry

import (
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

func (e *EnvEntry) Get(path string) (any, bool) {
	v, ok := e.entry.Get(path)

	if !ok {
		return nil, false
	}

	vString, ok := v.(string)

	if !ok {
		return nil, false
	}

	v = os.Getenv(vString)

	if v == "" {
		return nil, false
	}

	return v, true
}
