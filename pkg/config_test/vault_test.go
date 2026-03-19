//go:build goconfig_vault

package config_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
	"github.com/boolka/goconfig/pkg/vault"
	vaultStub "github.com/boolka/goconfig/pkg/vault_stub"
	vaultApi "github.com/hashicorp/vault/api"
)

const vaultToken = "root"

func prepareSecret(ctx context.Context, t *testing.T, addr string) {
	vaultStubClient := vaultStub.NewVaultClient(addr, vaultToken, http.DefaultClient)

	err := vaultStubClient.WriteSecret(ctx, "secret", "goconfig_secret", map[string]any{
		"password1": "abc123",
		"password2": "correct horse battery staple",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = vaultStubClient.DeleteSecret(ctx, "secret", "goconfig_secret")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestVaultBrokenPath(t *testing.T) {
	ctx := context.Background()

	client, err := vaultApi.NewClient(vaultApi.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.New(ctx, config.Options{
		Directory:   "testdata/vault",
		VaultClient: client,
	})
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := cfg.Get(ctx, "broken_field"); ok || v != vault.ErrInvalidPath {
		t.Fatal(v, ok)
	}
}

func TestVault(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultStub.NewVaultServer(vaultToken)

	prepareSecret(ctx, t, vaultServer.URL)

	vaultCfg := vaultApi.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vaultApi.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(vaultToken)

	cfg, err := config.New(ctx, config.Options{
		Directory:   "testdata/vault",
		VaultClient: client,
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
