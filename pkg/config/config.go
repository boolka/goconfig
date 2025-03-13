package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/boolka/goconfig/pkg/entry"
)

// Config options:
//   - Directory: path to config files. May be multiple directories delimited by os specific path separator
//   - Instance: is concrete instance number in multi instance deployments
//   - Deployment: is concrete deployment. For example "production" or "development"
//   - Hostname: mean current machine hostname
type Options struct {
	Directory  string
	Instance   string
	Deployment string
	Hostname   string
}

type configEntry struct {
	entry.Entry
	source cfgSource
}

type Config struct {
	sources []configEntry
}

// Creates new config instance. Provide Options object to set
// config path and etc. This function call use recovery
// mechanism and will never panics.
func New(options Options) (cfg *Config, err error) {
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

	var sources = []configEntry{}

	for _, configDir := range strings.Split(configDirs, string(os.PathListSeparator)) {
		dirSources, err := loadDir(configDir, hostname, deployment, instance)

		if err != nil {
			return nil, err
		}

		sources = append(sources, dirSources...)
	}

	slices.SortFunc(sources, func(a, b configEntry) int {
		return int(b.source) - int(a.source)
	})

	cfg = &Config{
		sources: sources,
	}

	return cfg, nil
}

// Get method takes dot delimited configuration path and returns value if any.
// Second returned value states if it was found and follows comma ok idiom at all.
// If this method is used with nil config receiver it will return ErrUninitialized
// sentinel error in first return value.
func (c *Config) Get(path string) (any, bool) {
	if c == nil {
		return ErrUninitialized, false
	}

	var v any
	var ok = false

	for _, source := range c.sources {
		v, ok = source.Get(path)

		if ok {
			return v, ok
		}
	}

	return v, ok
}

func loadDir(configDir string, hostname, deployment, instance string) ([]configEntry, error) {
	var sources []configEntry

	dirEntries, err := os.ReadDir(configDir)

	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		var newEntry entry.Entry
		var err error

		fName := fileName(dirEntry.Name())
		source, fileDeployment, fileInstance := fileSource(fName, hostname)

		if (fileDeployment != "" && fileDeployment != deployment) || (fileInstance != "" && fileInstance != instance) {
			continue
		}

		f, err := os.Open(filepath.Join(configDir, dirEntry.Name()))

		if err != nil {
			return nil, err
		}

		defer f.Close()

		switch filepath.Ext(dirEntry.Name()) {
		case ".json":
			newEntry, err = entry.NewJson(f)
		case ".toml":
			newEntry, err = entry.NewToml(f)
		case ".yaml", ".yml":
			newEntry, err = entry.NewYaml(f)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		if fName == "env" {
			newEntry, err = entry.NewEnv(newEntry)
		}

		if err != nil {
			return nil, err
		}

		sources = append(sources, configEntry{
			Entry:  newEntry,
			source: source,
		})
	}

	return sources, nil
}
