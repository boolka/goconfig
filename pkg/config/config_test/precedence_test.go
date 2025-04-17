package config_test

import (
	"context"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

// No instance & deployment
func TestPrecedence(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/precedence",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "custom"); !ok || v != "default.toml" {
		t.Fatal(v, ok)
	}
}

func TestInstanceEnvPrecedence(t *testing.T) {
	t.Setenv("GO_INSTANCE", "1")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/precedence",
		Instance:  "2",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "custom"); !ok || v != "default-2.toml" {
		t.Fatal(v, ok)
	}
}

func TestDeploymentEnvPrecedence(t *testing.T) {
	t.Setenv("GO_DEPLOYMENT", "testing")

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory:  "testdata/precedence",
		Deployment: "production",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "custom"); !ok || v != "production.toml" {
		t.Fatal(v, ok)
	}
}
