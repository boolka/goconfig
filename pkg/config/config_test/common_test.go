package config_test

import (
	"context"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestDefaultConfigDir(t *testing.T) {
	t.Chdir("testdata")

	_, err := config.New(config.Options{})

	if err != nil {
		t.Fatal(err)
	}
}

func TestEnvConfigPath(t *testing.T) {
	t.Setenv("GO_CONFIG_PATH", "testdata/config")

	cfg, err := config.New(config.Options{})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}

func TestCascadeConfig(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "config_file_env_custom")

	cfg, err := config.New(config.Options{
		Directory:  "testdata/config",
		Instance:   "1",
		Deployment: "testing",
		Hostname:   "host-name",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("default-1"); !ok || v != "default-1.json" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("testing"); !ok || v != "testing.yaml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("testing-1"); !ok || v != "testing-1.yaml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("host-name"); !ok || v != "host-name.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("host-name-1"); !ok || v != "host-name-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("host-name-testing"); !ok || v != "host-name-testing.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("host-name-testing-1"); !ok || v != "host-name-testing-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("local"); !ok || v != "local.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("local-1"); !ok || v != "local-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("local-testing"); !ok || v != "local-testing.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("local-testing-1"); !ok || v != "local-testing-1.toml" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("env"); !ok || v != "config_file_env_custom" {
		t.Fatal(v, ok)
	}
}

func TestSkipEmptyEnv(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "testdata/deep",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("deep.deep.custom"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}

func TestConfigUninitialized(t *testing.T) {
	t.Parallel()

	var cfg *config.Config

	err, ok := cfg.Get("field")

	if ok || err != config.ErrUninitialized {
		t.Fatal(err, ok)
	}
}

func TestEmptyDir(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/empty",
	})

	if err != nil {
		t.Fatal(err)
	}

	v, ok := cfg.Get("empty")

	if ok || v != nil {
		t.Fatal(v, ok)
	}
}

func TestSpoofEnvValue(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "initial_value")

	cfg, err := config.New(config.Options{
		Directory:  "testdata/config",
		Instance:   "1",
		Deployment: "testing",
		Hostname:   "host-name",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("TEST_FILE_ENV", "spoofed_value")

	if v, ok := cfg.Get("env"); !ok || v != "spoofed_value" {
		t.Fatal(v, ok)
	}
}

func TestNilValue(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "testdata/nil",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("empty"); !ok || v != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("empty.empty"); ok || v != nil {
		t.Fatal(err)
	}
}

func TestContext(t *testing.T) {
	t.Chdir("testdata")

	cfg, err := config.New(config.Options{})

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	v, ok := cfg.GetContext(ctx, "default")

	if !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}

func TestContextErr(t *testing.T) {
	t.Chdir("testdata")

	cfg, err := config.New(config.Options{})

	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	v, ok := cfg.GetContext(ctx, "default")

	if ok || v != config.ErrContextDone {
		t.Fatal(v, ok)
	}
}
