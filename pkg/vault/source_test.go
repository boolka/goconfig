//go:build goconfig_vault

package vault_test

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/vault"
	vaultStub "github.com/boolka/goconfig/pkg/vault_stub"
	vaultApi "github.com/hashicorp/vault/api"
)

func prepareSecret(ctx context.Context, t *testing.T, client *vaultStub.VaultClient) {
	err := client.WriteSecret(ctx, "secret", "goconfig_secret", map[string]any{
		"password1": "abc123",
		"password2": "correct horse battery staple",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = client.DeleteSecret(ctx, "secret", "goconfig_secret")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestVaultSource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vaultServer := vaultStub.NewVaultServer("root")
	t.Cleanup(vaultServer.Close)

	c := vaultStub.NewVaultClient(vaultServer.URL, "root", vaultServer.Client())

	prepareSecret(ctx, t, c)

	vaultCfg := vaultApi.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vaultApi.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("root")

	cfg, err := vault.NewVaultSource(ctx, os.DirFS("testdata").(fs.ReadDirFS), "vault.toml", client)
	if cfg == nil || err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := cfg.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultUnauthClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, err := vaultApi.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	c, err := vault.NewVaultSource(ctx, os.DirFS("testdata").(fs.ReadDirFS), "vault.toml", client)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := c.Get(ctx, "password1"); ok {
		t.Fatal("unauth pass")
	}
}

func TestVaultBrokenPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	v, err := vault.NewVaultSource(ctx, os.DirFS("testdata").(fs.ReadDirFS), "vault.toml", nil)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := v.Get(ctx, "broken_field"); ok || v != vault.ErrInvalidPath {
		t.Fatal(v, ok)
	}

	if v, ok := v.Get(ctx, "broken_field1"); ok || v != vault.ErrInvalidPath {
		t.Fatal(v, ok)
	}
}
