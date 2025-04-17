package entry_test

import (
	"context"
	"os"
	"slices"
	"testing"

	"github.com/boolka/goconfig/pkg/entry"
)

func TestTomlEntry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	f, err := os.Open("./testdata/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	l, err := entry.NewToml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := l.Get(ctx, "empty_field"); ok {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "num"); !ok || v != int64(1) {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "f_num"); !ok || v != 1. {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	var n []int64

	if v, ok := l.Get(ctx, "n_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			n = append(n, i.(int64))
		}
	}

	if !slices.EqualFunc(n, []int64{1, 2, 3}, func(f1, f2 int64) bool {
		return f1 == f2
	}) {
		t.Fatal(n, []int64{1, 2, 3})
	}

	var s []string

	if v, ok := l.Get(ctx, "s_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			s = append(s, i.(string))
		}
	}

	if !slices.EqualFunc(s, []string{"1", "2", "3"}, func(s1, s2 string) bool {
		return s1 == s2
	}) {
		t.Fatal(s, []string{"1", "2", "3"})
	}

	if v, ok := l.Get(ctx, "str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}

	// nested
	if v, ok := l.Get(ctx, "obj.num"); !ok || v != int64(1) {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "f_num"); !ok || v != 1. {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "obj.e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get(ctx, "obj.bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	var n1 []int64

	if v, ok := l.Get(ctx, "obj.n_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			n1 = append(n1, i.(int64))
		}
	}

	if !slices.EqualFunc(n, []int64{1, 2, 3}, func(f1, f2 int64) bool {
		return f1 == f2
	}) {
		t.Fatal(n1, []int64{1, 2, 3})
	}

	var s1 []string

	if v, ok := l.Get(ctx, "obj.s_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			s1 = append(s1, i.(string))
		}
	}

	if !slices.EqualFunc(s1, []string{"1", "2", "3"}, func(s1, s2 string) bool {
		return s1 == s2
	}) {
		t.Fatal(s1, []string{"1", "2", "3"})
	}

	if v, ok := l.Get(ctx, "obj.str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}
}

func TestTomlContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	f, err := os.Open("./testdata/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	_, err = entry.NewToml(ctx, f)
	if err == nil {
		t.Fatal(err)
	}
}
