package config

import (
	"context"
	"io/fs"
	"os"
	"slices"
	"testing"

	"github.com/boolka/goconfig/pkg/source"
)

func TestConfigSources(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("sort", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "host-name")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)

		if !slices.IsSortedFunc(sources, func(a, b *source.Source) int {
			return int(b.Type) - int(a.Type)
		}) {
			t.Fatal("incorrect order of sources")
		}

		if len(sources) != 14 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 14, "found", len(sources))
		}
	})

	t.Run("load by hostname", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "host-name")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "host-name", "", "")

		if len(sources) != 5 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 5, "found", len(sources))
		}
	})

	t.Run("load by deployment", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "", "testing", "")

		if len(sources) != 6 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 6, "found", len(sources))
		}
	})

	t.Run("load by instance", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "", "", "1")

		if len(sources) != 6 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 6, "found", len(sources))
		}
	})

	t.Run("load by hostname & instance", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "host-name")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "host-name", "", "1")

		if len(sources) != 8 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 8, "found", len(sources))
		}
	})

	t.Run("load by hostname & deployment", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "host-name")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "host-name", "testing", "")

		if len(sources) != 8 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 8, "found", len(sources))
		}
	})

	t.Run("load by deployment & instance", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "", "testing", "1")

		if len(sources) != 10 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 10, "found", len(sources))
		}
	})

	t.Run("load by hostname & deployment & instance", func(t *testing.T) {
		sources, err := loadDir(ctx, os.DirFS("testdata").(fs.ReadDirFS), "config", "host-name")
		if err != nil {
			t.Fatal(err)
		}

		sortSources(sources)
		sources = filterSources(sources, "host-name", "testing", "1")

		if len(sources) != 14 {
			for i, src := range sources {
				t.Error(i, src.FilePath)
			}
			t.Fatal("expected", 14, "found", len(sources))
		}
	})
}
