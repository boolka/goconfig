package config_test

import (
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestMultiInstance(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/multi_instance",
		Instance:  "2",
	})

	if err != nil {
		t.Fatal(err)
	}

	v, ok := cfg.Get("field")

	if !ok || v != "2" {
		t.Fatal(v, ok)
	}
}
