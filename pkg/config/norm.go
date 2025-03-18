package config

import "math"

func norm(v any) any {
	switch v := v.(type) {
	case float64:
		if v == math.Trunc(v) {
			switch {
			case v > math.MaxInt && v <= math.MaxUint:
				return uint(v)
			case v >= math.MinInt && v <= math.MaxInt:
				return int(v)
			}
		}
	case float32:
		v64 := float64(v)
		if float64(v64) == math.Trunc(float64(v64)) {
			switch {
			case v > math.MaxInt && v <= math.MaxUint:
				return uint(v)
			case v >= math.MinInt && v <= math.MaxInt:
				return int(v)
			}
		}
	case uint64:
		switch {
		case v > math.MaxInt && v <= math.MaxUint:
			return uint(v)
		case v <= math.MaxInt:
			return int(v)
		}
	case uint32:
		switch {
		case v > math.MaxInt32 && v <= math.MaxUint32:
			return uint(v)
		case v <= math.MaxInt32:
			return int(v)
		}
	case uint16:
		return int(v)
	case uint8:
		return int(v)
	case int64:
		if v >= math.MinInt && v <= math.MaxInt {
			return int(v)
		}
	case int32:
		return int(v)
	case int16:
		return int(v)
	case int8:
		return int(v)
	}

	return v
}
