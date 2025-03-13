package entry

import (
	"encoding/json"
	"io"
)

type JsonEntry struct {
	data map[string]any
}

func NewJson(r io.Reader) (*JsonEntry, error) {
	data := make(map[string]any)

	err := json.NewDecoder(r).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &JsonEntry{
		data: data,
	}, nil
}

func (e *JsonEntry) Get(path string) (any, bool) {
	return get(e.data, path)
}
