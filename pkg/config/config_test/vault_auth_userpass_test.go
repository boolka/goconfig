//go:build vault

package config_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
	vaultMock "github.com/boolka/goconfig/pkg/vault"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

func prepareVaultUserName(ctx context.Context, t *testing.T, token string) {
	u := vaultMock.NewClient("http://127.0.0.1:8200", token, http.DefaultClient)

	err := u.EnableUserPassMethod(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DisableUserPassMethod(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	policyName := "goconfig_vault_policy"

	err = u.CreateSecretPolicy(ctx, policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DeleteSecretPolicy(ctx, policyName)
		if err != nil {
			t.Fatal(err)
		}
	})

	err = u.CreateUserPass(ctx, "goconfig_vault_login", "goconfig_vault_password", policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DeleteUserPass(ctx, "goconfig_vault_login")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestVaultUserPassAuth(t *testing.T) {
	ctx := context.Background()

	prepareVaultUserName(ctx, t, "root")
	prepareVaultSecret(ctx, t, "root")

	auth, err := userPassAuth.NewUserpassAuth("goconfig_vault_login", &userPassAuth.Password{
		FromString: "goconfig_vault_password",
	})
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/userpass",
		VaultAuth: auth,
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

func TestVaultUserPassAuthFromFile(t *testing.T) {
	ctx := context.Background()

	prepareVaultUserName(ctx, t, "root")
	prepareVaultSecret(ctx, t, "root")

	t.Setenv("VAULT_PASSWORD", "goconfig_vault_password")

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/userpass",
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
