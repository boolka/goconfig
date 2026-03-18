//go:build goconfig_vault

package vault

import "strings"

const trimChars = "\t\r\n\x20"

func parsePath(cfgPath string) (mount string, secret string, key string, err error) {
	sepPath := strings.Split(cfgPath, ",")

	switch len(sepPath) {
	case 2:
		return strings.Trim(sepPath[0], trimChars), strings.Trim(sepPath[1], trimChars), "", nil
	case 3:
		return strings.Trim(sepPath[0], trimChars), strings.Trim(sepPath[1], trimChars), strings.Trim(sepPath[2], trimChars), nil
	}

	return "", "", "", ErrInvalidPath
}
