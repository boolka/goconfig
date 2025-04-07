package config

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/boolka/goconfig/pkg/entry"
)

func loadDir(ctx context.Context, fsys fs.ReadDirFS, directory string, hostname, deployment, instance string) ([]configEntry, error) {
	var sources []configEntry

	dirEntries, err := fs.ReadDir(fsys, directory)
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

		if fName != "env" && fName != "vault" && ((fileDeployment != "" && fileDeployment != deployment) || (fileInstance != "" && fileInstance != instance)) {
			continue
		}

		f, err := fsys.Open(filepath.Join(directory, dirEntry.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		switch filepath.Ext(dirEntry.Name()) {
		case ".json":
			newEntry, err = entry.NewJson(ctx, f)
		case ".toml":
			newEntry, err = entry.NewToml(ctx, f)
		case ".yaml", ".yml":
			newEntry, err = entry.NewYaml(ctx, f)
		default:
			continue
		}

		if err != nil {
			return nil, err
		}

		sources = append(sources, configEntry{
			Entry:  newEntry,
			source: source,
			file:   fName,
		})
	}

	return sources, nil
}
