package config_test

import (
	"context"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestMultiInstance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "./testdata/multi_instance",
		Instance:  "2",
	})

	if err != nil {
		t.Fatal(err)
	}

	v, ok := cfg.Get(ctx, "field")

	if !ok || v != "2" {
		t.Fatal(v, ok)
	}
}
