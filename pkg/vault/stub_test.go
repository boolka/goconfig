//go:build !goconfig_vault

package vault_test

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/vault"
)

func TestStubVault(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, err := vault.NewVaultSource(ctx, os.DirFS("testdata").(fs.ReadDirFS), "vault.toml", nil)
	if cfg == nil || err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "secret,goconfig_secret" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "secret,goconfig_secret,password2" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "broken_field"); !ok || v != "broken_value" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "broken_field1"); !ok || v != 1 {
		t.Fatal(v, ok)
	}
}
