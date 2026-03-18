package config

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/boolka/goconfig/pkg/datamap"
	"github.com/boolka/goconfig/pkg/source"
)

func loadDir(ctx context.Context, dirFs fs.ReadDirFS, directory string, hostname string) ([]*source.Source, error) {
	var sources []*source.Source

	dirEntries, err := fs.ReadDir(dirFs, directory)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		fName := dirEntry.Name()
		if strings.HasPrefix(fName, ".") {
			continue
		}
		fPath := filepath.Join(directory, fName)

		src, err := source.New(ctx, dirFs, fPath, hostname)
		if err != nil {
			if err == datamap.ErrUnknownFileSource {
				continue
			}

			return nil, err
		}

		sources = append(sources, src)
	}

	return sources, nil
}
