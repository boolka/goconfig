package entry_test

import (
	"context"
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/entry"
)

func TestEnv(t *testing.T) {
	ctx := context.Background()

	f, err := os.Open("./testdata/env.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := entry.NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	envEntry, err := entry.NewEnv(tomlEntry)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := envEntry.Get(ctx, "custom"); ok {
		t.Fatal(v, ok)
	}

	if v, ok := envEntry.Get(ctx, "obj.custom"); ok {
		t.Fatal(v, ok)
	}

	t.Setenv("CUSTOM_ENV", "variable1234")
	t.Setenv("CUSTOM_ENV_1", "variable4321")

	if v, ok := envEntry.Get(ctx, "custom"); !ok || v != "variable1234" {
		t.Fatal(v, ok)
	}

	if v, ok := envEntry.Get(ctx, "obj.custom"); !ok || v != "variable4321" {
		t.Fatal(v, ok)
	}

	t.Setenv("CUSTOM_ENV", "new_variable1234")
	t.Setenv("CUSTOM_ENV_1", "new_variable4321")

	if v, ok := envEntry.Get(ctx, "custom"); !ok || v != "new_variable1234" {
		t.Fatal(v, ok)
	}

	if v, ok := envEntry.Get(ctx, "obj.custom"); !ok || v != "new_variable4321" {
		t.Fatal(v, ok)
	}

	if v, ok := envEntry.Get(ctx, "empty_field"); ok {
		t.Fatal(v, ok)
	}
}
