package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/boolka/goconfig/pkg/entry"
	vault "github.com/hashicorp/vault/api"
	appRoleAuth "github.com/hashicorp/vault/api/auth/approle"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

// Config options:
//
//   - Directory: path to config files. May be multiple directories delimited by os specific path separator
//
//   - Instance: is concrete instance number in multi instance deployments
//
//   - Deployment: is concrete deployment. For example "production" or "development"
//
//   - Hostname: mean current machine hostname
//
//   - Logger: produce debug info and errors to provided logger. Module will be silent if nil was received
//
//   - VaultClient: vault client if you don't want to create a new one
//
//   - VaultAuth: is AuthMethod interface from [vault] api module that provides Login method
//
// [vault]: https://github.com/hashicorp/vault
type Options struct {
	Directory   string
	Instance    string
	Deployment  string
	Hostname    string
	Logger      *slog.Logger
	VaultClient *vault.Client
	VaultAuth   vault.AuthMethod
}

type configEntry struct {
	entry.Entry
	source cfgSource
	file   string
}

type Config struct {
	sources []configEntry
	logger  *slog.Logger
}

// Creates new config instance. Provide Options object to set
// config path and etc. If configuration directory is empty the ErrEmptyDir
// sentinel error will be returned.
func New(ctx context.Context, options Options) (cfg *Config, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			cfg = nil

			e, ok := panicErr.(error)

			if ok {
				err = fmt.Errorf("module \"github.com/boolka/goconfig\" panics with error: %w", e)
			} else {
				err = fmt.Errorf("module \"github.com/boolka/goconfig\" panics with value: %v", panicErr)
			}
		}
	}()

	var configDirs = options.Directory
	var instance = options.Instance
	var deployment = options.Deployment
	var hostname = options.Hostname
	var logger = options.Logger

	if logger != nil {
		logger = logger.With("module", "github.com/boolka/goconfig")

		ctx = ContextWithLogger(ctx, logger)
	}

	if logger != nil && logger.Enabled(ctx, slog.LevelDebug) {
		logger.DebugContext(ctx, "new config instantiating")
	}

	if hostname == "" {
		h, err := os.Hostname()

		if err != nil {
			hostname = ""
		} else {
			hostname = strings.Split(h, ".")[0]
		}

		hostname = h
	}

	if deployment == "" {
		deployment = os.Getenv("GO_DEPLOYMENT")
	}

	if instance == "" {
		instance = os.Getenv("GO_INSTANCE")
	}

	if configDirs == "" {
		configDirs = os.Getenv("GO_CONFIG_PATH")

		if configDirs == "" {
			configDirs = "./config"
		}
	}

	if logger != nil && logger.Enabled(ctx, slog.LevelDebug) {
		logger.DebugContext(ctx, fmt.Sprintf(`goconfig hostname: "%s"`, hostname))
		logger.DebugContext(ctx, fmt.Sprintf(`goconfig deployment: "%s"`, deployment))
		logger.DebugContext(ctx, fmt.Sprintf(`goconfig instance: "%s"`, instance))
		logger.DebugContext(ctx, fmt.Sprintf(`goconfig configDirs: "%s"`, configDirs))
	}

	var sources = []configEntry{}

	for _, configDir := range strings.Split(configDirs, string(os.PathListSeparator)) {
		dirSources, err := loadDir(ctx, configDir, hostname, deployment, instance)

		if err != nil {
			return nil, err
		}

		sources = append(sources, dirSources...)
	}

	for i, source := range sources {
		var replaceEntry entry.Entry
		var src cfgSource
		var err error

		switch source.file {
		case "env":
			src = envSrc
			replaceEntry, err = entry.NewEnv(source)
		case "vault":
			src = vaultSrc
			var vaultClient *vault.Client = options.VaultClient

			if vaultClient == nil {
				// load vault config from source files
				vaultCfg := loadVaultConfig(ctx, vault.DefaultConfig(), sources)

				// set goconfig logger
				if options.Logger != nil {
					vaultCfg.Logger = options.Logger
				}

				vaultClient, err = vault.NewClient(vaultCfg)
				if err != nil {
					return nil, fmt.Errorf("can't create vault client: %w\n", err)
				}
			}

			var vaultAuth = options.VaultAuth

			if vaultAuth == nil {
				// trying to load auth from source files
				authType, creds := loadVaultAuth(ctx, sources)

				switch authType {
				case tokenVaultConfigAuthType:
					vaultClient.SetToken(creds[0])
				case appRoleVaultConfigAuthType:
					roleId := creds[0]
					secretId := creds[1]
					mount := creds[2]

					vaultAuth, err = appRoleAuth.NewAppRoleAuth(roleId,
						&appRoleAuth.SecretID{FromString: secretId},
						appRoleAuth.WithMountPath(mount),
					)
					if err != nil {
						return nil, fmt.Errorf("invalid approle auth credentials: %w", err)
					}
				case userNameVaultConfigAuthType:
					username := creds[0]
					password := creds[1]
					mount := creds[2]

					vaultAuth, err = userPassAuth.NewUserpassAuth(username,
						&userPassAuth.Password{FromString: password},
						userPassAuth.WithMountPath(mount),
					)
					if err != nil {
						return nil, fmt.Errorf("invalid username auth credentials: %w", err)
					}
				case unknownVaultConfigAuthType:
					return nil, entry.ErrVaultUnauthorized
				}
			}

			replaceEntry, err = entry.NewVault(ctx, source, vaultClient, vaultAuth)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		sources[i] = configEntry{
			Entry:  replaceEntry,
			source: src,
			file:   source.file,
		}
	}

	slices.SortFunc(sources, func(a, b configEntry) int {
		return int(b.source) - int(a.source)
	})

	if logger != nil && logger.Enabled(ctx, slog.LevelDebug) {
		for i, source := range sources {
			logger.DebugContext(ctx, fmt.Sprintf(`%d loaded source, file: "%s", cfgSource: %d`, i, source.file, source.source))
		}
	}

	if len(sources) == 0 {
		return nil, ErrEmptyDir
	}

	cfg = &Config{
		sources: sources,
		logger:  logger,
	}

	return cfg, nil
}

// Get method takes dot delimited configuration path and returns value if any.
// Second returned value states if it was found and follows comma ok idiom at all.
func (c *Config) Get(ctx context.Context, path string) (any, bool) {
	if c.logger != nil && c.logger.Enabled(ctx, slog.LevelDebug) {
		ctx = ContextWithLogger(ctx, c.logger)
		c.logger.DebugContext(ctx, fmt.Sprintf(`trying to get "%s" field`, path))
	}

	var v any
	var ok = false

	done := make(chan struct{})

	go func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				v = panicErr
				ok = false
			}

			close(done)
		}()

		v, ok = get(ctx, c.sources, path)
	}()

	select {
	case <-ctx.Done():
		v, ok = ctx.Err(), false
	case <-done:
	}

	return v, ok
}

// MustGet method is the same as Get except that it panics if the path does not exist
func (c *Config) MustGet(ctx context.Context, path string) any {
	if c.logger != nil && c.logger.Enabled(ctx, slog.LevelDebug) {
		ctx = ContextWithLogger(ctx, c.logger)
		c.logger.DebugContext(ctx, fmt.Sprintf(`trying to get "%s" field`, path))
	}

	var v any
	var ok = false

	done := make(chan struct{})

	go func() {
		defer close(done)

		v, ok = get(ctx, c.sources, path)
	}()

	select {
	case <-ctx.Done():
		panic(ctx.Err())
	case <-done:
	}

	if !ok {
		panic("path " + path + " not found")
	}

	return v
}

// Return created or directly passed vault client
func (c *Config) GetVaultClient() *vault.Client {
	for _, source := range c.sources {
		if source.file == "vault" {
			if e, ok := source.Entry.(*entry.VaultEntry); ok {
				return e.Client()
			}
		}
	}

	return nil
}

func get(ctx context.Context, sources []configEntry, path string) (any, bool) {
	var v any
	var ok = false

	for _, source := range sources {
		v, ok = source.Get(ctx, path)

		if ok {
			break
		} else {
			if logger, ok := LoggerFromContext(ctx); ok && logger.Enabled(ctx, slog.LevelDebug) {
				switch v := v.(type) {
				case error:
					logger.Debug(v.Error())
				}
			}
		}
	}

	return norm(v), ok
}
