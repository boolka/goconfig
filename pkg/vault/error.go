//go:build goconfig_vault

package vault

import "errors"

var ErrInvalidPath = errors.New("invalid vault path")
