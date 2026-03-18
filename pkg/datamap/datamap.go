package datamap

import (
	"strings"

	"github.com/boolka/goconfig/pkg/normalization"
)

func GetByPath(data map[string]any, path string) (any, bool) {
	cuts := strings.Split(path, ".")
	deep := len(cuts) - 1

	var value any = data

	for i, cut := range cuts {
		if v, ok := value.(map[string]any); ok {
			value, ok = v[cut]

			if ok && i == deep {
				return normalization.Number(value), true
			}
		}
	}

	return nil, false
}
