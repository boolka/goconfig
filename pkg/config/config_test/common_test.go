package config_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
	"golang.org/x/sync/errgroup"
)

func TestDefaultConfigDir(t *testing.T) {
	t.Parallel()

	_, err := config.New(context.Background(), config.Options{
		Directory: "testdata/config",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigLogger(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/config",
		Logger:    logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}

func TestEnvConfigPath(t *testing.T) {
	t.Setenv("GO_CONFIG_PATH", "testdata/config")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}

func TestCascadeConfig(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "config_file_env_custom")
	t.Setenv("TEST_FILE_VAULT", "config_file_vault_custom")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory:  "testdata/config",
		Instance:   "1",
		Deployment: "testing",
		Hostname:   "host-name",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "default-1"); !ok || v != "default-1.json" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "testing"); !ok || v != "testing.yaml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "testing-1"); !ok || v != "testing-1.yaml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "host-name"); !ok || v != "host-name.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "host-name-1"); !ok || v != "host-name-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "host-name-testing"); !ok || v != "host-name-testing.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "host-name-testing-1"); !ok || v != "host-name-testing-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "local"); !ok || v != "local.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "local-1"); !ok || v != "local-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "local-testing"); !ok || v != "local-testing.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "local-testing-1"); !ok || v != "local-testing-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "env"); !ok || v != "config_file_env_custom" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "vault"); !ok || v != "config_file_vault_custom" {
		t.Fatal(v, ok)
	}
}

func TestMustGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic")
		} else {
			if err != "path not_exist not found" {
				t.Fatal(err)
			}
		}
	}()

	cfg, err := config.New(context.Background(), config.Options{
		Directory: "testdata/config",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v := cfg.MustGet(ctx, "default"); v != "default.json" {
		t.Fatal(err)
	}

	cfg.MustGet(ctx, "not_exist")
}

func TestSkipEmptyEnv(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/deep",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "deep.deep.custom"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}

func TestConfigUninitialized(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Fatal(err)
		}
	}()

	var cfg *config.Config

	cfg.Get(context.Background(), "field")
}

func TestEmptyDir(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	_, err := config.New(ctx, config.Options{
		Directory: "./testdata/empty",
	})
	if err != config.ErrEmptyDir {
		t.Fatal(err)
	}
}

func TestSpoofEnvValue(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "initial_value")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory:  "testdata/config",
		Instance:   "1",
		Deployment: "testing",
		Hostname:   "host-name",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("TEST_FILE_ENV", "spoofed_value")

	if v, ok := cfg.Get(ctx, "env"); !ok || v != "spoofed_value" {
		t.Fatal(v, ok)
	}
}

func TestNilValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/nil",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "empty"); !ok || v != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "empty.empty"); ok || v != nil {
		t.Fatal(err)
	}
}

func TestNewContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := config.New(ctx, config.Options{
		Directory: "testdata/config",
	})

	if err != context.Canceled {
		t.Fatal(err)
	}
}

func TestGetContextCanceled(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(context.Background(), config.Options{
		Directory: "testdata/config",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	v, ok := cfg.Get(ctx, "default")

	if ok || v != context.Canceled {
		t.Fatal(v, ok)
	}
}

func TestMustGetContextCanceled(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic")
		} else {
			err := err.(error)

			errors.Is(err, context.Canceled)
		}
	}()

	cfg, err := config.New(context.Background(), config.Options{
		Directory: "testdata/config",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg.MustGet(ctx, "default")
}

func TestMultiSource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(context.Background(), config.Options{
		Directory: "testdata/multi_source",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "toml"); !ok || v != "default.toml" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "json"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "yaml"); !ok || v != "default.yaml" {
		t.Fatal(v, ok)
	}
}

func TestConcurrent(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "config_file_env_custom")
	t.Setenv("TEST_FILE_VAULT", "config_file_vault_custom")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/config",
	})
	if err != nil {
		t.Fatal(err)
	}

	keys := []string{
		"default",
		"default-1",
		"testing",
		"testing-1",
		"host-name",
		"host-name-1",
		"host-name-testing",
		"host-name-testing-1",
		"local",
		"local-1",
		"local-testing",
		"local-testing-1",
		"env",
		"vault",
	}

	values := map[string]string{
		"default":             "default.json",
		"default-1":           "default-1.json",
		"testing":             "testing.yaml",
		"testing-1":           "testing-1.yaml",
		"host-name":           "host-name.toml",
		"host-name-1":         "host-name-1.toml",
		"host-name-testing":   "host-name-testing.toml",
		"host-name-testing-1": "host-name-testing-1.toml",
		"local":               "local.toml",
		"local-1":             "local-1.toml",
		"local-testing":       "local-testing.toml",
		"local-testing-1":     "local-testing-1.toml",
		"env":                 "config_file_env_custom",
		"vault":               "config_file_vault_custom",
	}

	eg := errgroup.Group{}

	for range 1000 {
		eg.Go(func() error {
			k := keys[rand.IntN(len(keys))]

			if cfg.MustGet(ctx, k) != values[k] {
				return fmt.Errorf("Got: %s, expected: %s", cfg.MustGet(ctx, k), values[k])
			}

			return nil
		})

		eg.Go(func() error {
			k := keys[rand.IntN(len(keys))]

			if v, ok := cfg.Get(ctx, k); !ok || v != values[k] {
				return fmt.Errorf("Got: %s, expected: %s", v, values[k])
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
