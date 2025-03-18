package config

import "errors"

var ErrUninitialized = errors.New("config is uninitialized")
var ErrContextDone = errors.New("context is done")
