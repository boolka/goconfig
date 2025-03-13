package entry

import "strings"

func get(data map[string]any, path string) (any, bool) {
	cuts := strings.Split(path, ".")
	deep := len(cuts) - 1

	var value any = data

	for i, cut := range cuts {
		v, ok := value.(map[string]any)

		if ok {
			value, ok = v[cut]

			if ok && i == deep {
				return value, true
			}
		}
	}

	return value, false
}
