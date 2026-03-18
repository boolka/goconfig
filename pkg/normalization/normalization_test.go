package normalization_test

import (
	"math"
	"testing"

	"github.com/boolka/goconfig/pkg/normalization"
)

func TestFloat64Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(float64(math.MaxUint)).(uint); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float64(math.MaxInt)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float64(0)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float64(math.MinInt)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float64(math.MaxUint * 2)).(float64); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(-float64(math.MaxUint * 2)).(float64); !ok {
		t.Fatal(v)
	}
}

func TestFloat32Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(float32(math.MaxUint64)).(uint); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float32(math.MaxInt)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float32(0)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float32(math.MinInt)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(float32(math.MaxUint64 * 2)).(float32); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(-float32(math.MaxUint64 * 2)).(float32); !ok {
		t.Fatal(v)
	}
}

func TestUint64Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(uint64(math.MaxUint64)).(uint); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint64(math.MaxInt64)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint64(0)).(int); !ok {
		t.Fatal(v)
	}
}

func TestUint32Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(uint32(math.MaxUint32)).(uint); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint32(math.MaxInt32)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint32(0)).(int); !ok {
		t.Fatal(v)
	}
}

func TestUint16Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(uint16(math.MaxUint16)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint16(math.MaxInt16)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint16(0)).(int); !ok {
		t.Fatal(v)
	}
}

func TestUint8Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(uint8(math.MaxUint8)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint8(math.MaxInt8)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(uint8(0)).(int); !ok {
		t.Fatal(v)
	}
}

func TestInt64Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(int64(math.MaxInt64)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int64(0)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int64(math.MinInt64)).(int); !ok {
		t.Fatal(v)
	}
}

func TestInt32Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(int32(math.MaxInt32)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int32(0)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int32(math.MinInt32)).(int); !ok {
		t.Fatal(v)
	}
}

func TestInt16Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(int16(math.MaxInt16)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int16(0)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int16(math.MinInt16)).(int); !ok {
		t.Fatal(v)
	}
}

func TestInt8Normalization(t *testing.T) {
	t.Parallel()

	if v, ok := normalization.Number(int8(math.MaxInt8)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int8(0)).(int); !ok {
		t.Fatal(v)
	}

	if v, ok := normalization.Number(int8(math.MinInt8)).(int); !ok {
		t.Fatal(v)
	}
}
