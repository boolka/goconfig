package source

import (
	"context"
	"io/fs"

	"github.com/boolka/goconfig/pkg/file"
)

type Originer interface {
	Get(context.Context, string) (any, bool)
}

type Source struct {
	Originer
	DirFs      fs.ReadDirFS
	Type       SourceType
	FilePath   string
	Hostname   string
	Deployment string
	Instance   string
}

func New(ctx context.Context, dirFs fs.ReadDirFS, fpath, hostname string) (*Source, error) {
	fileName := file.FileName(fpath)
	srcType, deployment, instance := ParseFilename(fileName, hostname)

	o := &Source{
		DirFs:      dirFs,
		Type:       srcType,
		Deployment: deployment,
		Instance:   instance,
		FilePath:   fpath,
		Hostname:   hostname,
	}

	return o, nil
}

func (s *Source) Get(ctx context.Context, path string) (any, bool) {
	return s.Originer.Get(ctx, path)
}
