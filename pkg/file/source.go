package file

import (
	"context"
	"io/fs"

	"github.com/boolka/goconfig/pkg/datamap"
)

type FileSource struct {
	data map[string]any
}

func NewPlainFileSource(ctx context.Context, dirFs fs.ReadDirFS, fpath string) (*FileSource, error) {
	data, err := datamap.NewDataMapFromFile(ctx, dirFs, fpath)
	if err != nil {
		return nil, err
	}

	return &FileSource{
		data: data,
	}, nil
}

func (f *FileSource) Get(_ context.Context, path string) (any, bool) {
	return datamap.GetByPath(f.data, path)
}
