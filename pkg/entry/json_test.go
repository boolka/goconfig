package entry_test

import (
	"context"
	"os"
	"slices"
	"sync"
	"testing"

	"github.com/boolka/goconfig/pkg/entry"
)

func readJsonFields(ctx context.Context, t *testing.T, jsonEntry *entry.JsonEntry) {
	if v, ok := jsonEntry.Get(ctx, "num"); !ok || v != 1. {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "nil"); !ok || v != nil {
		t.Fatal(v, ok)
	}

	var n []float64

	if v, ok := jsonEntry.Get(ctx, "n_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			n = append(n, i.(float64))
		}
	}

	if !slices.EqualFunc(n, []float64{1, 2, 3}, func(f1, f2 float64) bool {
		return f1 == f2
	}) {
		t.Fatal(n, []float64{1, 2, 3})
	}

	var s []string

	if v, ok := jsonEntry.Get(ctx, "s_arr"); !ok {
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

	if v, ok := jsonEntry.Get(ctx, "str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}

	// nested
	if v, ok := jsonEntry.Get(ctx, "obj.num"); !ok || v != 1. {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "obj.e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "obj.bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "obj.nil"); !ok || v != nil {
		t.Fatal(v, ok)
	}

	var n1 []float64

	if v, ok := jsonEntry.Get(ctx, "obj.n_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			n1 = append(n1, i.(float64))
		}
	}

	if !slices.EqualFunc(n, []float64{1, 2, 3}, func(f1, f2 float64) bool {
		return f1 == f2
	}) {
		t.Fatal(n1, []float64{1, 2, 3})
	}

	var s1 []string

	if v, ok := jsonEntry.Get(ctx, "obj.s_arr"); !ok {
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

	if v, ok := jsonEntry.Get(ctx, "obj.str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}
}

func TestJsonEntry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	f, err := os.Open("./testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	jsonEntry, err := entry.NewJson(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1000)

	for range 1000 {
		go func() {
			defer wg.Done()

			readJsonFields(ctx, t, jsonEntry)
		}()
	}

	wg.Wait()
}

func TestJsonMissedField(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	f, err := os.Open("./testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	jsonEntry, err := entry.NewJson(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := jsonEntry.Get(ctx, "some"); ok {
		t.Fatal(v, ok)
	}

	if v, ok := jsonEntry.Get(ctx, "some.nested.missed.value"); ok {
		t.Fatal(v, ok)
	}
}

func TestJsonDeepInner(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	f, err := os.Open("./testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	jsonEntry, err := entry.NewJson(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := jsonEntry.Get(ctx, "obj.obj.obj.inner"); !ok || v != 1. {
		t.Fatal(v, ok)
	}
}

func TestJsonContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	f, err := os.Open("./testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	_, err = entry.NewJson(ctx, f)
	if err == nil {
		t.Fatal(err)
	}
}
