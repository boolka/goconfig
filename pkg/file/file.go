package file

import (
	"path/filepath"
	"strings"
)

func FileName(fpath string) string {
	return strings.TrimSuffix(filepath.Base(fpath), filepath.Ext(fpath))
}
