package config_test

import (
	"context"
	"errors"
	"io/fs"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestUnexpectedConfigDirectory(t *testing.T) {
	t.Parallel()

	_, err := config.New(context.Background(), config.Options{
		Directory: "unexpected/directory/path",
	})

	var pathError *fs.PathError
	if !errors.As(err, &pathError) {
		t.Fatal(err)
	}
}

func TestSkipUnsupportedFile(t *testing.T) {
	t.Parallel()

	_, err := config.New(context.Background(), config.Options{
		Directory: "./testdata/unsupported",
	})
	if err != config.ErrEmptyDir {
		t.Fatal(err)
	}
}

func TestBrokenJson(t *testing.T) {
	t.Parallel()

	_, err := config.New(context.Background(), config.Options{
		Directory: "./testdata/unexpected/json",
	})

	if err == nil {
		t.Fatal("unexpected json")
	}
}

func TestBrokenToml(t *testing.T) {
	t.Parallel()

	_, err := config.New(context.Background(), config.Options{
		Directory: "./testdata/unexpected/toml",
	})

	if err == nil {
		t.Fatal("unexpected toml")
	}
}

func TestBrokenYaml(t *testing.T) {
	t.Parallel()

	_, err := config.New(context.Background(), config.Options{
		Directory: "./testdata/unexpected/yaml",
	})

	if err == nil {
		t.Fatal("unexpected yaml")
	}
}
