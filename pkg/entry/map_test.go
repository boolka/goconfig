package entry

import "testing"

func TestMapGetter(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"field1": 1,
	}

	if v, ok := get(m, "field1"); !ok || v != 1 {
		t.Fatal(v, ok)
	}
}

func TestNestedMapGetter(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"field1": map[string]any{
			"field2": map[string]any{
				"field3": 1,
			},
		},
	}

	if v, ok := get(m, "field1.field2.field3"); !ok || v != 1 {
		t.Fatal(v, ok)
	}
}

func TestEmptyMap(t *testing.T) {
	t.Parallel()

	m := map[string]any{}

	if v, ok := get(m, "field1.field2.field3"); ok || v != nil {
		t.Fatal(v, ok)
	}
}

func TestNilMapGetter(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"field1": map[string]any{
			"field2": map[string]any{
				"field3": nil,
			},
		},
	}

	if v, ok := get(m, "field1.field2.field3"); !ok || v != nil {
		t.Fatal(v, ok)
	}
}
