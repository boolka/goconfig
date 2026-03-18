package config

import (
	"context"
	"slices"

	"github.com/boolka/goconfig/pkg/file"
	goconfigLogger "github.com/boolka/goconfig/pkg/logger"
	"github.com/boolka/goconfig/pkg/source"
)

func sortSources(sources []*source.Source) {
	slices.SortFunc(sources, func(a, b *source.Source) int {
		return int(b.Type) - int(a.Type)
	})
}

// retain only relevant to current environment sources
func filterSources(sources []*source.Source, hostname, deployment, instance string) []*source.Source {
	return slices.DeleteFunc(sources, func(o *source.Source) bool {
		if o.Type == source.EnvSrc || o.Type == source.VaultSrc || o.Type == source.DefSrc || o.Type == source.LocSrc {
			return false
		}

		if (o.Hostname == "" || o.Hostname == hostname) &&
			(o.Deployment == "" || o.Deployment == deployment) &&
			(o.Instance == "" || o.Instance == instance) {
			return false
		}

		return true
	})
}

func searchSources(ctx context.Context, sources []*source.Source, path string, files ...string) (any, bool) {
	var v any
	var ok bool

	for _, src := range sources {
		if len(files) > 0 && !slices.ContainsFunc(files, func(fp string) bool {
			return fp == src.FilePath || file.FileName(fp) == file.FileName(src.FilePath)
		}) {
			continue
		}

		v, ok = src.Get(ctx, path)
		if ok {
			break
		}

		if logger, ok := goconfigLogger.LoggerFromContext(ctx); ok {
			switch v := v.(type) {
			case error:
				logger.InfoContext(ctx, v.Error())
			}
		}
	}

	return v, ok
}
