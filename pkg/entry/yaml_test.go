package entry_test

import (
	"context"
	"os"
	"slices"
	"sync"
	"testing"

	"github.com/boolka/goconfig/pkg/entry"
)

func readYamlFields(ctx context.Context, t *testing.T, yamlEntry *entry.YamlEntry) {
	if v, ok := yamlEntry.Get(ctx, "num"); !ok || v != 1 {
		t.Fatal(v, ok)
	}

	if v, ok := yamlEntry.Get(ctx, "e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := yamlEntry.Get(ctx, "bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	if v, ok := yamlEntry.Get(ctx, "nil"); !ok || v != nil {
		t.Fatal(v, ok)
	}

	var n []int

	if v, ok := yamlEntry.Get(ctx, "n_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			n = append(n, i.(int))
		}
	}

	if !slices.EqualFunc(n, []int{1, 2, 3}, func(f1, f2 int) bool {
		return f1 == f2
	}) {
		t.Fatal(n, []int{1, 2, 3})
	}

	var s []string

	if v, ok := yamlEntry.Get(ctx, "s_arr"); !ok {
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

	if v, ok := yamlEntry.Get(ctx, "str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}

	// nested
	if v, ok := yamlEntry.Get(ctx, "obj.num"); !ok || v != 1 {
		t.Fatal(v, ok)
	}

	if v, ok := yamlEntry.Get(ctx, "obj.e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := yamlEntry.Get(ctx, "obj.bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	if v, ok := yamlEntry.Get(ctx, "obj.nil"); !ok || v != nil {
		t.Fatal(v, ok)
	}

	var n1 []int

	if v, ok := yamlEntry.Get(ctx, "obj.n_arr"); !ok {
		t.Fatal(v, ok)
	} else {
		for _, i := range v.([]any) {
			n1 = append(n1, i.(int))
		}
	}

	if !slices.EqualFunc(n, []int{1, 2, 3}, func(f1, f2 int) bool {
		return f1 == f2
	}) {
		t.Fatal(n1, []int{1, 2, 3})
	}

	var s1 []string

	if v, ok := yamlEntry.Get(ctx, "obj.s_arr"); !ok {
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

	if v, ok := yamlEntry.Get(ctx, "obj.str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}
}

func TestYamlEntry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	f, err := os.Open("./testdata/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	yamlEntry, err := entry.NewYaml(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1000)

	for range 1000 {
		go func() {
			defer wg.Done()

			readYamlFields(ctx, t, yamlEntry)
		}()
	}

	wg.Wait()
}

func TestYamlContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	f, err := os.Open("./testdata/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	_, err = entry.NewYaml(ctx, f)
	if err == nil {
		t.Fatal(err)
	}
}
