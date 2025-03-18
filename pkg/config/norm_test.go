package config

import (
	"math"
	"testing"
)

func TestFloat64Normalization(t *testing.T) {
	t.Parallel()

	f64_uint := norm(float64(math.MaxUint64))
	f64_uint = f64_uint.(uint)

	f64_int := norm(float64(math.MaxInt))
	f64_int = f64_int.(int)

	f64_int = norm(float64(0))
	f64_int = f64_int.(int)

	f64_int = norm(float64(math.MinInt))
	f64_int = f64_int.(int)

	f64_f64 := norm(float64(math.MaxUint64 * 2))
	f64_f64 = f64_f64.(float64)

	f64_f64 = norm(-float64(math.MaxUint64 * 2))
	f64_f64 = f64_f64.(float64)
}

func TestFloat32Normalization(t *testing.T) {
	t.Parallel()

	f32_uint := norm(float32(math.MaxUint64))
	f32_uint = f32_uint.(uint)

	f32_int := norm(float32(math.MaxInt))
	f32_int = f32_int.(int)

	f32_int = norm(float32(0))
	f32_int = f32_int.(int)

	f32_int = norm(float32(math.MinInt))
	f32_int = f32_int.(int)

	f32_f32 := norm(float32(math.MaxUint64 * 2))
	f32_f32 = f32_f32.(float32)

	f32_f32 = norm(-float32(math.MaxUint64 * 2))
	f32_f32 = f32_f32.(float32)
}

func TestUint64Normalization(t *testing.T) {
	t.Parallel()

	uint64_uint := norm(uint64(math.MaxUint64))
	uint64_uint = uint64_uint.(uint)

	uint64_int := norm(uint64(math.MaxInt64))
	uint64_int = uint64_int.(int)

	uint64_int = norm(uint64(0))
	uint64_int = uint64_int.(int)
}

func TestUint32Normalization(t *testing.T) {
	t.Parallel()

	uint32_uint := norm(uint32(math.MaxUint32))
	uint32_uint = uint32_uint.(uint)

	uint32_int := norm(uint32(math.MaxInt32))
	uint32_int = uint32_int.(int)

	uint32_int = norm(uint32(0))
	uint32_int = uint32_int.(int)
}

func TestUint16Normalization(t *testing.T) {
	t.Parallel()

	uint16_uint := norm(uint16(math.MaxUint16))
	uint16_uint = uint16_uint.(int)

	uint16_int := norm(uint16(math.MaxInt16))
	uint16_int = uint16_int.(int)

	uint16_int = norm(uint16(0))
	uint16_int = uint16_int.(int)
}

func TestUint8Normalization(t *testing.T) {
	t.Parallel()

	uint8_uint := norm(uint8(math.MaxUint8))
	uint8_uint = uint8_uint.(int)

	uint8_int := norm(uint8(math.MaxInt8))
	uint8_int = uint8_int.(int)

	uint8_int = norm(uint8(0))
	uint8_int = uint8_int.(int)
}

func TestInt64Normalization(t *testing.T) {
	t.Parallel()

	int64_int := norm(int64(math.MaxInt64))
	int64_int = int64_int.(int)

	int64_int = norm(int64(0))
	int64_int = int64_int.(int)

	int64_int = norm(int64(math.MinInt64))
	int64_int = int64_int.(int)
}

func TestInt32Normalization(t *testing.T) {
	t.Parallel()

	int32_int := norm(int32(math.MaxInt32))
	int32_int = int32_int.(int)

	int32_int = norm(int32(0))
	int32_int = int32_int.(int)

	int32_int = norm(int32(math.MinInt32))
	int32_int = int32_int.(int)
}

func TestInt16Normalization(t *testing.T) {
	t.Parallel()

	int16_int := norm(int16(math.MaxInt16))
	int16_int = int16_int.(int)

	int16_int = norm(int16(0))
	int16_int = int16_int.(int)

	int16_int = norm(int16(math.MinInt16))
	int16_int = int16_int.(int)
}

func TestInt8Normalization(t *testing.T) {
	t.Parallel()

	int8_int := norm(int8(math.MaxInt8))
	int8_int = int8_int.(int)

	int8_int = norm(int8(0))
	int8_int = int8_int.(int)

	int8_int = norm(int8(math.MinInt8))
	int8_int = int8_int.(int)
}
