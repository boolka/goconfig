//go:build vault

package config_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
	vaultMock "github.com/boolka/goconfig/pkg/vault"
)

func prepareVaultSecret(ctx context.Context, t *testing.T, token string) {
	client := vaultMock.NewClient("http://127.0.0.1:8200", token, http.DefaultClient)

	err := client.WriteSecret(ctx, "secret", "goconfig_secret_1", map[string]any{
		"password1": "abc123",
		"password2": "correct horse battery staple",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = client.DeleteSecret(ctx, "secret", "goconfig_secret_1")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestVaultSingleFile(t *testing.T) {
	ctx := context.Background()
	prepareVaultSecret(ctx, t, "root")

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/config",
		VaultAuth: vaultMock.NewTokenAuth("root"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultTokenAuth(t *testing.T) {
	ctx := context.Background()
	prepareVaultSecret(ctx, t, "root")

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/token",
		VaultAuth: vaultMock.NewTokenAuth("root"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultTokenAuthFromFile(t *testing.T) {
	ctx := context.Background()
	prepareVaultSecret(ctx, t, "root")

	t.Setenv("VAULT_TOKEN", "root")

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/token",
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultGetCertainSource(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/config",
		VaultAuth: vaultMock.NewTokenAuth("root"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "password1_default" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "password2_default" {
		t.Fatal(v, ok)
	}

	if v, ok := cfg.Get(ctx, "password1", "vault"); ok {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2", "vault"); ok {
		t.Fatal(v, ok)
	}
}

func TestVaultMustGetCertainSource(t *testing.T) {
	ctx := context.Background()

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("must panic")
		}
	}()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/config",
		VaultAuth: vaultMock.NewTokenAuth("root"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "password1_default" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "password2_default" {
		t.Fatal(v, ok)
	}

	cfg.MustGet(ctx, "password1", "vault")
}
