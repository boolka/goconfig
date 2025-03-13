package config_test

import (
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

// No instance & deployment
func TestPrecedence(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "testdata/precedence",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("custom"); !ok || v != "default.toml" {
		t.Fatal(v, ok)
	}
}

func TestInstanceEnvPrecedence(t *testing.T) {
	t.Setenv("GO_INSTANCE", "1")

	cfg, err := config.New(config.Options{
		Directory: "testdata/precedence",
		Instance:  "2",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("custom"); !ok || v != "default-2.toml" {
		t.Fatal(v, ok)
	}
}

func TestDeploymentEnvPrecedence(t *testing.T) {
	t.Setenv("GO_DEPLOYMENT", "testing")

	cfg, err := config.New(config.Options{
		Directory:  "testdata/precedence",
		Deployment: "production",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("custom"); !ok || v != "production.toml" {
		t.Fatal(v, ok)
	}
}
