package config_test

import (
	"context"
	"embed"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

//go:embed testdata/config/*
var configDir embed.FS

func TestEmbed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		FileSystem: &configDir,
		Directory:  "testdata/config",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}
