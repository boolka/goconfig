package config

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/boolka/goconfig/pkg/entry"
	tokenRoleAuth "github.com/boolka/goconfig/pkg/vault"
	vault "github.com/hashicorp/vault/api"
	appRoleAuth "github.com/hashicorp/vault/api/auth/approle"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

// Config options:
//
//   - Directory: path to config files. May be multiple directories delimited by os specific path separator
//
//   - FileSystem: an interface that provides access to some point in the file system. The directory path
//     will be considered relative if it is non nil. It can be useful when the configuration is embedded
//
//   - Instance: is concrete instance number in multi instance deployments
//
//   - Deployment: is concrete deployment. For example "production" or "development"
//
//   - Hostname: mean current machine hostname
//
//   - Logger: produce output to supplied logger. Module will be silent if nil was received.
//
//   - VaultClient: vault client if you don't want to create a new one
//
//   - VaultAuth: is AuthMethod interface from [vault] api module that provides Login method
//
// [vault]: https://github.com/hashicorp/vault
type Options struct {
	Directory   string
	FileSystem  fs.ReadDirFS
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

	var directory = options.Directory
	var instance = options.Instance
	var deployment = options.Deployment
	var hostname = options.Hostname
	var fsys = options.FileSystem
	var logger = options.Logger

	if logger != nil {
		logger = logger.With("module", "github.com/boolka/goconfig")

		ctx = ContextWithLogger(ctx, logger)
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

	if directory == "" {
		directory = os.Getenv("GO_CONFIG_PATH")

		if directory == "" && fsys == nil {
			directory = "./config"
		}
	}

	if logger != nil && logger.Enabled(ctx, slog.LevelDebug) {
		logger.DebugContext(ctx, fmt.Sprintf(`directory: "%s"`, directory))
		logger.DebugContext(ctx, fmt.Sprintf(`fsys: "%t"`, fsys != nil))
		logger.DebugContext(ctx, fmt.Sprintf(`hostname: "%s"`, hostname))
		logger.DebugContext(ctx, fmt.Sprintf(`deployment: "%s"`, deployment))
		logger.DebugContext(ctx, fmt.Sprintf(`instance: "%s"`, instance))
	}

	var sources = []configEntry{}

	for _, dir := range strings.Split(directory, string(os.PathListSeparator)) {
		if dir == "" {
			continue
		}

		if fsys, ok := os.DirFS(dir).(fs.ReadDirFS); ok {
			dirSources, err := loadDir(ctx, fsys, ".", hostname, deployment, instance)
			if err != nil {
				return nil, err
			}

			sources = append(sources, dirSources...)
		}
	}

	if fsys != nil {
		dirSources, err := loadDir(ctx, fsys, directory, hostname, deployment, instance)
		if err != nil {
			return nil, err
		}

		sources = append(sources, dirSources...)
	}

	slices.SortFunc(sources, func(a, b configEntry) int {
		return int(b.source) - int(a.source)
	})

	for i, source := range sources {
		var replaceEntry entry.Entry
		var src cfgSource
		var err error

		switch fileName(source.file) {
		case "env":
			src = envSrc
			replaceEntry, err = entry.NewEnv(source)
		case "vault":
			src = vaultSrc
			var vaultClient *vault.Client = options.VaultClient

			if vaultClient == nil {
				// load vault config from source files
				vaultCfg := loadVaultConfig(ctx, vault.DefaultConfig(), sources)

				// vault disabled
				if vaultCfg == nil {
					// remove correspond file from sources
					sources = slices.Delete(sources, i, i+1)
					continue
				}

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
					vaultAuth = tokenRoleAuth.NewTokenAuth(creds[0])
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

type sourceValue struct {
	v  any
	ok bool
}

// Get method takes dot delimited configuration path and returns value if any.
// The last parameter specifies which files to allow for searching both with or without extension.
// If omitted, all files will be search through.
// The sequence of transmitted files does not change the original order for searching.
// Second returned value states if it was found and follows comma ok idiom at all.
func (c *Config) Get(ctx context.Context, path string, files ...string) (any, bool) {
	if c.logger != nil && c.logger.Enabled(ctx, slog.LevelDebug) {
		ctx = ContextWithLogger(ctx, c.logger)
		c.logger.DebugContext(ctx, fmt.Sprintf(`trying to get "%s" field`, path))
	}

	value := make(chan sourceValue)

	go func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				value <- sourceValue{
					v:  panicErr,
					ok: false,
				}
			}
		}()

		var cfgSources []configEntry = c.sources

		if len(files) != 0 {
			cfgSources = filterSources(ctx, c.sources, files...)
		}

		v, ok := searchThroughSources(ctx, cfgSources, path)

		value <- sourceValue{
			v:  v,
			ok: ok,
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err(), false
	case res := <-value:
		return res.v, res.ok
	}
}

// MustGet method is the same as Get except that it panics if the path does not exist
func (c *Config) MustGet(ctx context.Context, path string, files ...string) any {
	if c.logger != nil && c.logger.Enabled(ctx, slog.LevelDebug) {
		ctx = ContextWithLogger(ctx, c.logger)
		c.logger.DebugContext(ctx, fmt.Sprintf(`trying to get "%s" field`, path))
	}

	value := make(chan sourceValue)

	go func() {
		defer close(value)

		var cfgSources []configEntry = c.sources

		if len(files) != 0 {
			cfgSources = filterSources(ctx, c.sources, files...)
		}

		v, ok := searchThroughSources(ctx, cfgSources, path)

		value <- sourceValue{
			v:  v,
			ok: ok,
		}
	}()

	select {
	case <-ctx.Done():
		panic(ctx.Err())
	case res := <-value:
		if !res.ok {
			panic("path " + path + " not found")
		}

		return res.v
	}
}

// Return created or directly passed vault client
func (c *Config) GetVaultClient() *vault.Client {
	for _, source := range c.sources {
		if fileName(source.file) == "vault" {
			if e, ok := source.Entry.(*entry.VaultEntry); ok {
				return e.Client()
			}
		}
	}

	return nil
}

func filterSources(ctx context.Context, sources []configEntry, files ...string) []configEntry {
	var cfgSources []configEntry

	for _, source := range sources {
		if slices.ContainsFunc(files, func(file string) bool {
			return source.file == file || fileName(source.file) == file
		}) {
			cfgSources = append(cfgSources, source)
		}
	}

	return cfgSources
}

func searchThroughSources(ctx context.Context, sources []configEntry, path string) (any, bool) {
	var v any
	var ok bool

	for _, source := range sources {
		v, ok = source.Get(ctx, path)

		if ok {
			break
		} else {
			if logger, ok := LoggerFromContext(ctx); ok && logger.Enabled(ctx, slog.LevelInfo) {
				switch v := v.(type) {
				case error:
					logger.InfoContext(ctx, v.Error())
				}
			}
		}
	}

	return norm(v), ok
}
