package config

import (
	"context"
	"strconv"

	vault "github.com/hashicorp/vault/api"
)

func loadVaultConfig(ctx context.Context, config *vault.Config, entries []configEntry) *vault.Config {
	if isEnable, ok := searchThroughSources(ctx, entries, "goconfig.vault.enable"); ok {
		if isEnable, ok := isEnable.(bool); ok && !isEnable {
			return nil
		}
	}

	if address, ok := searchThroughSources(ctx, entries, "goconfig.vault.address"); ok {
		if address, ok := address.(string); ok {
			config.Address = address
		}
	}

	if minRetryWait, ok := searchThroughSources(ctx, entries, "goconfig.vault.min_retry_wait"); ok {
		if minRetryWait, ok := minRetryWait.(string); ok {
			d, err := parseDuration(minRetryWait)

			if err == nil {
				config.MinRetryWait = d
			}
		}
	}

	if maxRetryWait, ok := searchThroughSources(ctx, entries, "goconfig.vault.max_retry_wait"); ok {
		if maxRetryWait, ok := maxRetryWait.(string); ok {
			d, err := parseDuration(maxRetryWait)

			if err == nil {
				config.MaxRetryWait = d
			}
		}
	}

	if maxRetries, ok := searchThroughSources(ctx, entries, "goconfig.vault.max_retries"); ok {
		if maxRetries, ok := maxRetries.(string); ok {
			maxRetries, err := strconv.ParseInt(maxRetries, 10, 64)

			if err == nil {
				config.MaxRetries = int(maxRetries)
			}
		}
	}

	if timeout, ok := searchThroughSources(ctx, entries, "goconfig.vault.timeout"); ok {
		if timeout, ok := timeout.(string); ok {
			d, err := parseDuration(timeout)

			if err == nil {
				config.Timeout = d
			}
		}
	}

	return config
}

type VaultConfigAuthType int

const (
	unknownVaultConfigAuthType VaultConfigAuthType = iota
	tokenVaultConfigAuthType
	userNameVaultConfigAuthType
	appRoleVaultConfigAuthType
)

func loadVaultAuth(ctx context.Context, entries []configEntry) (VaultConfigAuthType, []string) {
	var vaultMount string

	if token, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.token"); ok {
		if token, ok := token.(string); ok {
			return tokenVaultConfigAuthType, []string{token}
		}
	}

	if cfgMount, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.mount"); ok {
		if cfgMount, ok := cfgMount.(string); ok {
			vaultMount = cfgMount
		}
	}

	if _, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.roleid"); ok {
		var roleId, secretId string

		if cfgRoleId, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.roleid"); ok {
			if cfgRoleId, ok := cfgRoleId.(string); ok {
				roleId = cfgRoleId
			}
		}
		if cfgSecretId, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.secretid"); ok {
			if cfgSecretId, ok := cfgSecretId.(string); ok {
				secretId = cfgSecretId
			}
		}
		if vaultMount == "" {
			vaultMount = "/approle"
		}

		return appRoleVaultConfigAuthType, []string{roleId, secretId, vaultMount}
	}

	if _, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.username"); ok {
		var username, password string

		if cfgUsername, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.username"); ok {
			if cfgUsername, ok := cfgUsername.(string); ok {
				username = cfgUsername
			}
		}
		if cfgPassword, ok := searchThroughSources(ctx, entries, "goconfig.vault.auth.password"); ok {
			if cfgPassword, ok := cfgPassword.(string); ok {
				password = cfgPassword
			}
		}
		if vaultMount == "" {
			vaultMount = "/userpass"
		}

		return userNameVaultConfigAuthType, []string{username, password, vaultMount}
	}

	return unknownVaultConfigAuthType, nil
}
