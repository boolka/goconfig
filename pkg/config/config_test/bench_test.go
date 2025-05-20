package config_test

import (
	"context"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func BenchmarkConfig(b *testing.B) {
	b.Setenv("TEST_FILE_ENV", "config_file_env_custom")

	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		cfg, err := config.New(ctx, config.Options{
			Directory:  "testdata/config",
			Instance:   "1",
			Deployment: "testing",
			Hostname:   "host-name",
		})
		if err != nil {
			b.Fatal(err)
		}

		if v, ok := cfg.Get(ctx, "default"); !ok || v != "default.json" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "default-1"); !ok || v != "default-1.json" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "testing"); !ok || v != "testing.yaml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "testing-1"); !ok || v != "testing-1.yaml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "host-name"); !ok || v != "host-name.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "host-name-1"); !ok || v != "host-name-1.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "host-name-testing"); !ok || v != "host-name-testing.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "host-name-testing-1"); !ok || v != "host-name-testing-1.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "local"); !ok || v != "local.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "local-1"); !ok || v != "local-1.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "local-testing"); !ok || v != "local-testing.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "local-testing-1"); !ok || v != "local-testing-1.toml" {
			b.Fatal(v, ok)
		}

		if v, ok := cfg.Get(ctx, "env"); !ok || v != "config_file_env_custom" {
			b.Fatal(v, ok)
		}
	}
}
