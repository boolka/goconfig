package config_test

import (
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestJsonNumber(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/serializers/json",
	})

	if err != nil {
		t.Fatal(err)
	}

	zero, ok := cfg.Get("zero")
	zero = zero.(int)

	if !ok {
		t.Fatal(zero, ok)
	}

	max_int32, ok := cfg.Get("max_int32")
	max_int32 = max_int32.(int)

	if !ok {
		t.Fatal(max_int32, ok)
	}

	min_int32, ok := cfg.Get("min_int32")
	min_int32 = min_int32.(int)

	if !ok {
		t.Fatal(min_int32, ok)
	}

	max_uint32, ok := cfg.Get("max_uint32")
	max_uint32 = max_uint32.(int)

	if !ok {
		t.Fatal(max_uint32, ok)
	}

	max_int64, ok := cfg.Get("max_int64")
	max_int64 = max_int64.(int)

	if !ok {
		t.Fatal(max_int64, ok)
	}

	min_int64, ok := cfg.Get("min_int64")
	min_int64 = min_int64.(int)

	if !ok {
		t.Fatal(min_int64, ok)
	}

	max_uint64, ok := cfg.Get("max_uint64")
	max_uint64 = max_uint64.(uint)

	if !ok {
		t.Fatal(max_uint64, ok)
	}

	max, ok := cfg.Get("max")
	max = max.(float64)

	if !ok {
		t.Fatal(max, ok)
	}
}

func TestTomlNumber(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/serializers/toml",
	})

	if err != nil {
		t.Fatal(err)
	}

	zero, ok := cfg.Get("zero")
	zero = zero.(int)

	if !ok {
		t.Fatal(zero, ok)
	}

	max_int32, ok := cfg.Get("max_int32")
	max_int32 = max_int32.(int)

	if !ok {
		t.Fatal(max_int32, ok)
	}

	min_int32, ok := cfg.Get("min_int32")
	min_int32 = min_int32.(int)

	if !ok {
		t.Fatal(min_int32, ok)
	}

	max_uint32, ok := cfg.Get("max_uint32")
	max_uint32 = max_uint32.(int)

	if !ok {
		t.Fatal(max_uint32, ok)
	}

	max_int64, ok := cfg.Get("max_int64")
	max_int64 = max_int64.(int)

	if !ok {
		t.Fatal(max_int64, ok)
	}

	min_int64, ok := cfg.Get("min_int64")
	min_int64 = min_int64.(int)

	if !ok {
		t.Fatal(min_int64, ok)
	}

	max, ok := cfg.Get("max")
	max = max.(float64)

	if !ok {
		t.Fatal(max, ok)
	}
}

func TestYamlNumber(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/serializers/yaml",
	})

	if err != nil {
		t.Fatal(err)
	}

	zero, ok := cfg.Get("zero")
	zero = zero.(int)

	if !ok {
		t.Fatal(zero, ok)
	}

	max_int32, ok := cfg.Get("max_int32")
	max_int32 = max_int32.(int)

	if !ok {
		t.Fatal(max_int32, ok)
	}

	min_int32, ok := cfg.Get("min_int32")
	min_int32 = min_int32.(int)

	if !ok {
		t.Fatal(min_int32, ok)
	}

	max_uint32, ok := cfg.Get("max_uint32")
	max_uint32 = max_uint32.(int)

	if !ok {
		t.Fatal(max_uint32, ok)
	}

	max_int64, ok := cfg.Get("max_int64")
	max_int64 = max_int64.(int)

	if !ok {
		t.Fatal(max_int64, ok)
	}

	min_int64, ok := cfg.Get("min_int64")
	min_int64 = min_int64.(int)

	if !ok {
		t.Fatal(min_int64, ok)
	}

	max_uint64, ok := cfg.Get("max_uint64")
	max_uint64 = max_uint64.(uint)

	if !ok {
		t.Fatal(max_uint64, ok)
	}

	max, ok := cfg.Get("max")
	max = max.(float64)

	if !ok {
		t.Fatal(max, ok)
	}
}
