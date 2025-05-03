package config_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/boolka/goconfig/pkg/config"
	"github.com/boolka/goconfig/pkg/entry"
)

func TestVaultConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/config",
	})
	if err != nil {
		t.Fatal(err)
	}

	client := cfg.GetVaultClient()
	if client == nil {
		t.Fatal("client is nil")
	}

	config := client.CloneConfig()

	if config.Address != "http://127.0.0.1:8200" {
		t.Fatal(config.Address)
	}

	if config.MinRetryWait != time.Duration(time.Second*3) {
		t.Fatal(config.MinRetryWait)
	}

	if config.MaxRetryWait != time.Duration(time.Second*5) {
		t.Fatal(config.MaxRetryWait)
	}

	if config.Timeout != time.Duration(time.Second*30) {
		t.Fatal(config.Timeout)
	}

	if client.Token() != "root" {
		t.Fatal(client.CloneToken())
	}
}

func TestUnavailableServer(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBuffer([]byte{})

	ctx := context.Background()

	_, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/unauthorized",
		Logger: slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	})
	if err != entry.ErrVaultUnauthorized {
		t.Fatal(err)
	}
}

func TestBrokenPath(t *testing.T) {
	t.Setenv("TEST_FILE_VAULT", "config_file_vault_custom")

	buf := bytes.NewBuffer([]byte{})

	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/config",
		Logger: slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "vault"); !ok || v != "config_file_vault_custom" {
		t.Fatal(v, ok)
	}

	if !strings.Contains(buf.String(), "invalid vault path") {
		t.Fatal("valid path")
	}
}
