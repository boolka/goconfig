package config

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/boolka/goconfig/pkg/env"
	"github.com/boolka/goconfig/pkg/file"
	goconfigLogger "github.com/boolka/goconfig/pkg/logger"
	"github.com/boolka/goconfig/pkg/source"
	vault "github.com/boolka/goconfig/pkg/vault"
)

// Config options:
//
//   - Directory: path to config files
//
//   - DirFS: an interface that provides access to some point in the file system. The directory path
//     will be considered relative if it is non nil. It can be useful when the configuration is embedded
//
//   - Instance: is concrete instance number in multi instance deployments
//
//   - Deployment: is concrete deployment. For example "production" or "development"
//
//   - Hostname: mean current machine hostname
//
//   - Logger: produce output to supplied logger. Config will be silent if nil was received.
//
//   - VaultClient: [vault] client. To use vault api set "vault" build tag. Otherwise it would act as plain file.
//
// [vault]: https://github.com/hashicorp/vault
type Options struct {
	Directory   string
	DirFS       fs.ReadDirFS
	Instance    string
	Deployment  string
	Hostname    string
	Logger      *slog.Logger
	VaultClient any
}

type Config struct {
	logger  *slog.Logger
	sources []*source.Source
}

// Creates new config instance. Provide Options object to set
// config path and etc. If configuration directory is empty the ErrEmptyDir
// sentinel error will be returned.
func New(ctx context.Context, options Options) (cfg *Config, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			cfg = nil

			if e, ok := panicErr.(error); ok {
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
	var logger = options.Logger

	if logger != nil {
		logger = logger.With("module", "github.com/boolka/goconfig")

		ctx = goconfigLogger.ContextWithLogger(ctx, logger)
	}

	if hostname == "" {
		if h, err := os.Hostname(); err != nil {
			hostname = ""
		} else {
			// cut hostname in case of dot specified
			hostname = strings.Split(h, ".")[0]
		}
	}

	if deployment == "" {
		deployment = os.Getenv("GO_DEPLOYMENT")
	}

	if instance == "" {
		instance = os.Getenv("GO_INSTANCE")
	}

	if directory == "" {
		directory = os.Getenv("GO_CONFIG_PATH")

		if directory == "" {
			directory = "config"
		}
	}

	if logger != nil {
		logger.DebugContext(ctx, fmt.Sprintf("directory: %s, fsys: %t, hostname: %s, deployment: %s, instance: %s", directory, options.DirFS != nil, hostname, deployment, instance))
	}

	var dirFs []fs.ReadDirFS

	if options.DirFS == nil {
		for _, dir := range strings.Split(directory, string(os.PathListSeparator)) {
			dir = strings.TrimSpace(dir)

			if fs, ok := os.DirFS(dir).(fs.ReadDirFS); ok {
				dirFs = append(dirFs, fs)
			} else {
				return nil, err
			}
		}

		directory = "."
	} else {
		dirFs = append(dirFs, options.DirFS)
	}

	var sources []*source.Source

	for _, fs := range dirFs {
		dirSources, err := loadDir(ctx, fs, directory, hostname)
		if err != nil {
			return nil, err
		}

		sources = append(sources, dirSources...)
	}

	sortSources(sources)
	sources = filterSources(sources, hostname, deployment, instance)

	if len(sources) == 0 {
		return nil, ErrEmptyDir
	}

	for i, src := range sources {
		var org source.Originer

		switch src.Type {
		case source.EnvSrc:
			org, err = env.NewEnvSource(ctx, src.DirFs, src.FilePath)
		case source.VaultSrc:
			org, err = vault.NewVaultSource(ctx, src.DirFs, src.FilePath, options.VaultClient)
		default:
			org, err = file.NewPlainFileSource(ctx, src.DirFs, src.FilePath)
		}

		if err != nil {
			return nil, err
		}

		src.Originer = org

		if logger != nil {
			logger.DebugContext(ctx, fmt.Sprintf(`%d loaded source, file: "%s", source type: %s`, i, src.FilePath, src.Type))
		}
	}

	return &Config{
		logger:  logger,
		sources: sources,
	}, nil
}
