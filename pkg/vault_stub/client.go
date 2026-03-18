package vault_mock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type VaultClient struct {
	addr   string
	token  string
	client *http.Client
}

func NewVaultClient(addr, token string, client *http.Client) *VaultClient {
	return &VaultClient{
		addr:   addr,
		token:  token,
		client: client,
	}
}

func (v *VaultClient) WriteSecret(ctx context.Context, mount, path string, data map[string]any) error {
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

func (v *VaultClient) DeleteSecret(ctx context.Context, mount, path string) error {
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
