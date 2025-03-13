package entry

import (
	"io"

	"github.com/pelletier/go-toml/v2"
)

type TomlEntry struct {
	data map[string]any
}

func NewToml(r io.Reader) (*TomlEntry, error) {
	data := make(map[string]any)

	err := toml.NewDecoder(r).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &TomlEntry{
		data: data,
	}, nil
}

func (e *TomlEntry) Get(path string) (any, bool) {
	return get(e.data, path)
}
