package entry

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YamlEntry struct {
	data map[string]any
}

func NewYaml(r io.Reader) (*YamlEntry, error) {
	data := make(map[string]any)

	err := yaml.NewDecoder(r).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &YamlEntry{
		data: data,
	}, nil
}

func (y *YamlEntry) Get(path string) (any, bool) {
	return get(y.data, path)
}
