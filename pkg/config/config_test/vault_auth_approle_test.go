//go:build vault

package config_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
	vaultMock "github.com/boolka/goconfig/pkg/vault"
	appRoleAuth "github.com/hashicorp/vault/api/auth/approle"
)

func prepareVaultAppRole(ctx context.Context, t *testing.T, token string) (string, string) {
	t.Helper()
	client := vaultMock.NewClient("http://127.0.0.1:8200", token, http.DefaultClient)

	err := client.EnableAppRoleMethod(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = client.DisableAppRoleMethod(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	policyName := "goconfig_approle_policy_name"

	err = client.CreateSecretPolicy(ctx, policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = client.DeleteSecretPolicy(ctx, policyName)
		if err != nil {
			t.Fatal(err)
		}
	})

	roleId, secretId, err := client.CreateAppRole(ctx, "goconfig_approle", policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = client.DeleteAppRole(ctx, "goconfig_approle")
		if err != nil {
			t.Fatal(err)
		}
	})

	return roleId, secretId
}

func TestVaultAppRole(t *testing.T) {
	ctx := context.Background()

	roleId, secretId := prepareVaultAppRole(ctx, t, "root")
	prepareVaultSecret(ctx, t, "root")

	auth, err := appRoleAuth.NewAppRoleAuth(
		roleId,
		&appRoleAuth.SecretID{FromString: secretId},
	)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/token",
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

func TestVaultAppRoleFromFile(t *testing.T) {
	ctx := context.Background()

	roleId, secretId := prepareVaultAppRole(ctx, t, "root")
	prepareVaultSecret(ctx, t, "root")

	f, err := os.Create("testdata/vault/approle/default-1.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	_, err = f.Write([]byte(fmt.Sprintf(`[goconfig.vault.auth]
roleid="%s"
secretid="%s"`, roleId, secretId)))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/approle",
		Instance:  "1",
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

func TestVaultAppRoleFromEnvFile(t *testing.T) {
	ctx := context.Background()

	roleId, secretId := prepareVaultAppRole(ctx, t, "root")
	prepareVaultSecret(ctx, t, "root")

	f, err := os.Create("testdata/vault/approle/env.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	_, err = f.Write([]byte(`[goconfig.vault.auth]
roleid="VAULT_ROLEID"
secretid="VAULT_SECRETID"`))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	t.Setenv("VAULT_ROLEID", roleId)
	t.Setenv("VAULT_SECRETID", secretId)

	cfg, err := config.New(ctx, config.Options{
		Directory: "testdata/vault/approle",
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
