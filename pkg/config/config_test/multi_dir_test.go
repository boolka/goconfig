package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestMultiConfigDir(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/deep" + string(os.PathListSeparator) + "testdata/config",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "deep.deep.custom"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}
