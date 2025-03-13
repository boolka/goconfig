package entry_test

import (
	"os"
	"slices"
	"testing"

	"github.com/boolka/goconfig/pkg/entry"
)

func TestJsonEntry(t *testing.T) {
	t.Parallel()

	f, err := os.Open("./testdata/config.json")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	l, err := entry.NewJson(f)

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := l.Get("num"); !ok || v != 1. {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("nil"); !ok || v != nil {
		t.Fatal(v, ok)
	}

	var n []float64

	if v, ok := l.Get("n_arr"); !ok {
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

	if v, ok := l.Get("s_arr"); !ok {
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

	if v, ok := l.Get("str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}

	// nested
	if v, ok := l.Get("obj.num"); !ok || v != 1. {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("obj.e_num"); !ok || v != float64(1e2) {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("obj.bool"); !ok || v != true {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("obj.nil"); !ok || v != nil {
		t.Fatal(v, ok)
	}

	var n1 []float64

	if v, ok := l.Get("obj.n_arr"); !ok {
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

	if v, ok := l.Get("obj.s_arr"); !ok {
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

	if v, ok := l.Get("obj.str"); !ok || v != "\"custom string\"\n" {
		t.Fatal(v, ok)
	}
}

func TestMissedField(t *testing.T) {
	t.Parallel()

	f, err := os.Open("./testdata/config.json")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	l, err := entry.NewJson(f)

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := l.Get("some"); ok {
		t.Fatal(v, ok)
	}

	if v, ok := l.Get("some.nested.missed.value"); ok {
		t.Fatal(v, ok)
	}
}

func TestDeepInner(t *testing.T) {
	t.Parallel()

	f, err := os.Open("./testdata/config.json")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	l, err := entry.NewJson(f)

	if err != nil {
		t.Fatal(err)
	}

	if v, ok := l.Get("obj.obj.obj.inner"); !ok || v != 1. {
		t.Fatal(v, ok)
	}
}
