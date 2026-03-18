package config

import (
	"context"
	"fmt"

	goconfigLogger "github.com/boolka/goconfig/pkg/logger"
)

type srcValue struct {
	v  any
	ok bool
}

// Get method takes dot delimited configuration path and returns value if any.
// The last parameter specifies which files to allow for searching both with or without extension.
// If omitted, all files will be search through.
// The sequence of transmitted files does not change the original order for searching.
// Second returned value states if it was found and follows comma ok idiom at all.
func (c *Config) Get(ctx context.Context, path string, files ...string) (any, bool) {
	if c.logger != nil {
		ctx = goconfigLogger.ContextWithLogger(ctx, c.logger)
		c.logger.DebugContext(ctx, fmt.Sprintf("get %s field with files %v", path, files))
	}

	done := make(chan srcValue, 1)

	go func() {
		if ctx.Err() != nil {
			return
		}

		defer func() {
			if panicErr := recover(); panicErr != nil {
				done <- srcValue{
					v:  panicErr,
					ok: false,
				}
			}
		}()

		v, ok := searchSources(ctx, c.sources, path, files...)

		done <- srcValue{
			v:  v,
			ok: ok,
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err(), false
	case res := <-done:
		return res.v, res.ok
	}
}

// MustGet method is the same as Get except that it panics if the path does not exist
func (c *Config) MustGet(ctx context.Context, path string, files ...string) any {
	if c.logger != nil {
		ctx = goconfigLogger.ContextWithLogger(ctx, c.logger)
		c.logger.DebugContext(ctx, fmt.Sprintf("must get %s field with files: %v", path, files))
	}

	done := make(chan srcValue, 1)

	go func() {
		if ctx.Err() != nil {
			return
		}

		v, ok := searchSources(ctx, c.sources, path, files...)

		done <- srcValue{
			v:  v,
			ok: ok,
		}
	}()

	select {
	case <-ctx.Done():
		panic(ctx.Err())
	case res := <-done:
		if !res.ok {
			panic("path " + path + " not found")
		}

		return res.v
	}
}
