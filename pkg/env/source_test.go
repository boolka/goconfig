package env_test

import (
	"context"
	"io/fs"
	"os"
	"testing"

	envEntry "github.com/boolka/goconfig/pkg/env"
)

func TestEnvSource(t *testing.T) {
	ctx := context.Background()

	envSource, err := envEntry.NewEnvSource(ctx, os.DirFS("testdata").(fs.ReadDirFS), "env.toml")
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := envSource.Get(ctx, "custom"); ok {
		t.Fatal(v, ok)
	}

	if v, ok := envSource.Get(ctx, "obj.custom"); ok {
		t.Fatal(v, ok)
	}

	t.Setenv("CUSTOM_ENV", "variable1234")
	t.Setenv("CUSTOM_ENV_1", "variable4321")

	if v, ok := envSource.Get(ctx, "custom"); !ok || v != "variable1234" {
		t.Fatal(v, ok)
	}

	if v, ok := envSource.Get(ctx, "obj.custom"); !ok || v != "variable4321" {
		t.Fatal(v, ok)
	}

	t.Setenv("CUSTOM_ENV", "new_variable1234")
	t.Setenv("CUSTOM_ENV_1", "new_variable4321")

	if v, ok := envSource.Get(ctx, "custom"); !ok || v != "new_variable1234" {
		t.Fatal(v, ok)
	}

	if v, ok := envSource.Get(ctx, "obj.custom"); !ok || v != "new_variable4321" {
		t.Fatal(v, ok)
	}

	if v, ok := envSource.Get(ctx, "empty_field"); ok {
		t.Fatal(v, ok)
	}
}
