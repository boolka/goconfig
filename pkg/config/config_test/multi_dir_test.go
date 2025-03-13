package config_test

import (
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestMultiConfigDir(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "testdata/deep" + string(os.PathListSeparator) + "testdata/config",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("deep.deep.custom"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get("default"); !ok || v != "default.json" {
		t.Fatal(v, ok)
	}
}
