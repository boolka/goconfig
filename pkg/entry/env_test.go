package entry_test

import (
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/entry"
)

func TestEnv(t *testing.T) {
	t.Setenv("CUSTOM_ENV", "variable1234")
	t.Setenv("CUSTOM_ENV_1", "variable4321")

	f, err := os.Open("./testdata/env.toml")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	tomlEntry, err := entry.NewToml(f)

	if err != nil {
		t.Fatal(err)
	}

	envEntry, err := entry.NewEnv(tomlEntry)

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := envEntry.Get("custom"); !ok || v != "variable1234" {
		t.Fatal(v, ok)
	}

	if v, ok := envEntry.Get("ojb.custom"); !ok || v != "variable4321" {
		t.Fatal(v, ok)
	}
}
