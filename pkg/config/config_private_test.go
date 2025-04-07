package config

import (
	"context"
	"io/fs"
	"os"
	"slices"
	"testing"
)

func TestConfigSourcesPrecedence(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := New(ctx, Options{
		Directory:  "config_test/testdata/config",
		Deployment: "testing",
		Hostname:   "host-name",
		Instance:   "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !slices.IsSortedFunc(cfg.sources, func(a, b configEntry) int {
		return int(b.source) - int(a.source)
	}) {
		t.Fatal("incorrect order of sources")
	}

	if len(cfg.sources) != 14 {
		for i, source := range cfg.sources {
			t.Error(i, source.file)
		}
		t.Fatal("expected", 14, "found", len(cfg.sources))
	}
}

func TestConfigFileSystem(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := New(ctx, Options{
		FileSystem: os.DirFS("config_test/testdata/config").(fs.ReadDirFS),
		Directory:  ".",
		Deployment: "testing",
		Hostname:   "host-name",
		Instance:   "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !slices.IsSortedFunc(cfg.sources, func(a, b configEntry) int {
		return int(b.source) - int(a.source)
	}) {
		t.Fatal("incorrect order of sources")
	}

	if len(cfg.sources) != 14 {
		for i, source := range cfg.sources {
			t.Error(i, source.file)
		}
		t.Fatal("expected", 14, "found", len(cfg.sources))
	}
}
