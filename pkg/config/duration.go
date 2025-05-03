package config

import (
	"errors"
	"strconv"
	"time"
)

func parseDuration(s string) (time.Duration, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return time.Duration(n) * time.Second, nil
	}

	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}

	return time.Duration(0), errors.New("invalid duration")
}
