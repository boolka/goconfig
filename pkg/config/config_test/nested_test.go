package config_test

import (
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestNested(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "testdata/nested",
		Instance:  "1",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("nested.custom-1"); !ok || v != "value-1" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("nested.custom-2"); !ok || v != "value-2" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("nested.custom-3"); !ok || v != "value-3" {
		t.Fatal(v, ok)
	}
}
