package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	addr   string
	token  string
	client *http.Client
}

func NewClient(addr, token string, httpClient *http.Client) *Client {
	return &Client{
		addr:   addr,
		token:  token,
		client: httpClient,
	}
}

func (v *Client) EnableUserPassMethod(ctx context.Context) error {
	b, err := json.Marshal(map[string]string{
		"path": "userpass",
		"type": "userpass",
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/sys/auth/userpass", v.addr)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 500 {
		return fmt.Errorf("vaultEnableUserPassMethod respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) DisableUserPassMethod(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/sys/auth/userpass", v.addr)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 500 {
		return fmt.Errorf("VaultDisableUserPassMethod respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) EnableAppRoleMethod(ctx context.Context) error {
	b, err := json.Marshal(map[string]string{
		"path": "approle",
		"type": "approle",
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/sys/auth/approle", v.addr)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 500 {
		return fmt.Errorf("VaultEnableAppRoleMethod respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) DisableAppRoleMethod(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/sys/auth/approle", v.addr)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 500 {
		return fmt.Errorf("VaultDisableAppRoleMethod respond with status: %s", res.Status)
	}

	return nil
}

const secretPolicy = `path "secret/*" {
	capabilities = ["create", "read", "update", "patch", "delete", "list"]
}`

func (v *Client) CreateSecretPolicy(ctx context.Context, policy string) error {
	b, err := json.Marshal(map[string]string{
		"name":   policy,
		"policy": secretPolicy,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/sys/policy/%s", v.addr, policy)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("vaultCreateSecretPolicy respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) DeleteSecretPolicy(ctx context.Context, policy string) error {
	url := fmt.Sprintf("%s/v1/sys/policy/%s", v.addr, policy)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("VaultDeleteSecretPolicy respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) CreateUserPass(ctx context.Context, username, password, policy string) error {
	b, err := json.Marshal(map[string]string{
		"username":       username,
		"password":       password,
		"token_policies": policy,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/auth/userpass/users/%s", v.addr, username)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("VaultCreateUserPass respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) DeleteUserPass(ctx context.Context, username string) error {
	url := fmt.Sprintf("%s/v1/auth/userpass/users/%s", v.addr, username)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return fmt.Errorf("VaultDeleteUserPass respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) CreateAppRole(ctx context.Context, roleName, policy string) (string, string, error) {
	b, err := json.Marshal(map[string]any{
		"role_name":      roleName,
		"token_policies": policy,
	})
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("%s/v1/auth/approle/role/%s", v.addr, roleName)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := v.client.Do(req)
	if err != nil {
		return "", "", err
	}

	if res.StatusCode >= 300 {
		return "", "", fmt.Errorf("vaultCreateAppRole respond with status: %s", res.Status)
	}

	b, err = json.Marshal(map[string]any{
		"role_name": roleName,
	})
	if err != nil {
		return "", "", err
	}

	url = fmt.Sprintf("%s/v1/auth/approle/role/%s/secret-id", v.addr, roleName)

	req, err = http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Content-Type", "application/json")

	res, err = v.client.Do(req)
	if err != nil {
		return "", "", err
	}

	if res.StatusCode >= 300 {
		return "", "", fmt.Errorf("vaultCreateAppRole respond with status: %s", res.Status)
	}

	defer res.Body.Close()

	r := map[string]any{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return "", "", err
	}

	secretId := r["data"].(map[string]any)["secret_id"]

	url = fmt.Sprintf("%s/v1/auth/approle/role/%s/role-id", v.addr, roleName)

	req, err = http.NewRequestWithContext(ctx,
		http.MethodGet,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Set("Accept", "application/json")

	res, err = v.client.Do(req)
	if err != nil {
		return "", "", err
	}

	if res.StatusCode >= 300 {
		return "", "", fmt.Errorf("vaultCreateAppRole respond with status: %s", res.Status)
	}

	defer res.Body.Close()

	r = map[string]any{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return "", "", err
	}

	roleId := r["data"].(map[string]any)["role_id"]

	return roleId.(string), secretId.(string), nil
}

func (v *Client) DeleteAppRole(ctx context.Context, roleName string) error {
	url := fmt.Sprintf("%s/v1/auth/approle/role/%s", v.addr, roleName)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("vaultDeleteAppRole respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) WriteSecret(ctx context.Context, mount, path string, data map[string]any) error {
	b, err := json.Marshal(map[string]any{
		"secret-mount-path": mount,
		"path":              path,
		"data":              data,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v1/%s/data/%s", v.addr, mount, path)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)
	req.Header.Add("Content-Type", "application-json")

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("vaultWriteSecret respond with status: %s", res.Status)
	}

	return nil
}

func (v *Client) DeleteSecret(ctx context.Context, mount, path string) error {
	url := fmt.Sprintf("%s/v1/%s/data/%s", v.addr, mount, path)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("VaultDeleteSecret respond with status: %s", res.Status)
	}

	return nil
}
