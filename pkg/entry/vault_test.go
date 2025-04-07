package entry

import (
	"context"
	"os"
	"testing"

	vaultMock "github.com/boolka/goconfig/pkg/vault"
	vault "github.com/hashicorp/vault/api"
	appRoleAuth "github.com/hashicorp/vault/api/auth/approle"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

func prepareSecret(ctx context.Context, t *testing.T, client *vaultMock.Client) {
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

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})
}

func TestVaultPredefinedClient(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(vaultServer.Close)

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	prepareSecret(ctx, t, u)

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("root")

	// ignore auth if client has token
	cfg, err := NewVault(ctx, tomlEntry, client, vaultMock.NewTokenAuth("broken"))
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

func TestVaultAuthMiss(t *testing.T) {
	ctx := context.Background()

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	client, err := vault.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = NewVault(ctx, tomlEntry, client, nil); err != ErrVaultUnauthorized {
		t.Fatal(err)
	}
}

func TestVaultTokenAuth(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(func() {
		vaultServer.Close()
	})

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	prepareSecret(ctx, t, u)

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	v, err := NewVault(ctx, tomlEntry, client, vaultMock.NewTokenAuth("root"))
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := v.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := v.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultBrokenTokenAuth(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(func() {
		vaultServer.Close()
	})

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	prepareSecret(ctx, t, u)

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)

	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("broken")

	v, err := NewVault(ctx, tomlEntry, client, nil)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := v.Get(ctx, "password1"); ok {
		t.Fatal(v, ok)
	}
}

func TestVaultUserPassAuth(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(func() {
		vaultServer.Close()
	})

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

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

	policyName := "goconfig_userpass_policy_name"

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

	err = u.CreateUserPass(ctx, "goconfig_userpass_login", "goconfig_userpass_password", policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DeleteUserPass(ctx, "goconfig_userpass_login")
		if err != nil {
			t.Fatal(err)
		}
	})

	prepareSecret(ctx, t, u)

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	vaultUserPassAuth, err := userPassAuth.NewUserpassAuth("goconfig_userpass_login",
		&userPassAuth.Password{FromString: "goconfig_userpass_password"},
		userPassAuth.WithMountPath("/userpass"),
	)
	if err != nil {
		t.Fatal(err)
	}

	v, err := NewVault(ctx, tomlEntry, client, vaultUserPassAuth)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := v.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := v.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultUserPassAuthDenied(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(func() {
		vaultServer.Close()
	})

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

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

	policyName := "goconfig_userpass_policy_name"

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

	err = u.CreateUserPass(ctx, "goconfig_userpass_login", "goconfig_userpass_password", policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DeleteUserPass(ctx, "goconfig_userpass_login")
		if err != nil {
			t.Fatal(err)
		}
	})

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	vaultUserPassAuth, err := userPassAuth.NewUserpassAuth("wrong_login",
		&userPassAuth.Password{FromString: "goconfig_userpass_password"},
		userPassAuth.WithMountPath("/userpass"),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewVault(ctx, tomlEntry, client, vaultUserPassAuth)
	if err == nil {
		t.Fatal(err)
	}
}

func TestVaultAppRoleAuth(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(func() {
		vaultServer.Close()
	})

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	err := u.EnableAppRoleMethod(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DisableAppRoleMethod(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	policyName := "goconfig_approle_policy_name"

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

	roleId, secretId, err := u.CreateAppRole(ctx, "goconfig_approle", policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DeleteAppRole(ctx, "goconfig_approle")
		if err != nil {
			t.Fatal(err)
		}
	})

	prepareSecret(ctx, t, u)

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	appRoleVaultAuth, err := appRoleAuth.NewAppRoleAuth(roleId,
		&appRoleAuth.SecretID{FromString: secretId},
		appRoleAuth.WithMountPath("/approle"),
	)
	if err != nil {
		t.Fatal(err)
	}

	v, err := NewVault(ctx, tomlEntry, client, appRoleVaultAuth)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := v.Get(ctx, "password1"); !ok || v != "abc123" {
		t.Fatal(v, ok)
	}
	if v, ok := v.Get(ctx, "userpass.password2"); !ok || v != "correct horse battery staple" {
		t.Fatal(v, ok)
	}
}

func TestVaultAppRoleAuthDenied(t *testing.T) {
	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(func() {
		vaultServer.Close()
	})

	u := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	err := u.EnableAppRoleMethod(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DisableAppRoleMethod(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	policyName := "goconfig_approle_policy_name"

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

	roleId, _, err := u.CreateAppRole(ctx, "goconfig_approle", policyName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = u.DeleteAppRole(ctx, "goconfig_approle")
		if err != nil {
			t.Fatal(err)
		}
	})

	prepareSecret(ctx, t, u)

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultServer.URL

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}

	appRoleVaultAuth, err := appRoleAuth.NewAppRoleAuth(roleId,
		&appRoleAuth.SecretID{FromString: "broken_secretid"},
		appRoleAuth.WithMountPath("/approle"),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewVault(ctx, tomlEntry, client, appRoleVaultAuth)
	if err == nil {
		t.Fatal(err)
	}
}

func TestVaultBrokenPath(t *testing.T) {
	ctx := context.Background()

	f, err := os.Open("./testdata/vault.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	tomlEntry, err := NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	client, err := vault.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("root")

	v, err := NewVault(ctx, tomlEntry, client, nil)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := v.Get(ctx, "broken_field"); ok || v != ErrVaultInvalidPath {
		t.Fatal(v, ok)
	}
}
