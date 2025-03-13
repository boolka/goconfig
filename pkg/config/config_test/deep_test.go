package config_test

import (
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestDeepConfig(t *testing.T) {
	t.Setenv("TEST_FILE_ENV", "config_file_env_custom")

	cfg, err := config.New(config.Options{
		Directory: "testdata/deep",
	})

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get("deep.deep.custom"); !ok || v != "config_file_env_custom" {
		t.Fatal(v, ok)
	}
}
