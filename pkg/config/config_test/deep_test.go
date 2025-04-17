package config_test

import (
	"context"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestDeepConfig(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "config_file_env_custom")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/deep",
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := cfg.Get(ctx, "deep.deep"); !ok {
		t.Fatal("deep.deep", ok)
	}

	if v, ok := cfg.Get(ctx, "deep.deep.custom"); !ok || v != "config_file_env_custom" {
		t.Fatal(v, ok)
	}
}
