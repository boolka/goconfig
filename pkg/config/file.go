package config

import (
	"path/filepath"
	"strings"
)

func fileName(path string) string {
	_, fileExt := filepath.Split(path)
	return strings.TrimSuffix(fileExt, filepath.Ext(fileExt))
}

func fileSource(filename, hostname string) (cfgSource, string, string) {
	if filename == hostname {
		return hostSrc, "", ""
	}

	fileSplit := strings.Split(filename, "-")

	switch {
	case filename == "default":
		return defSrc, "", ""
	case reDefInst.MatchString(filename):
		return defInstSrc, "", fileSplit[1]
	case filename == "local":
		return locSrc, "", ""
	case reLocInst.MatchString(filename):
		return locInstSrc, "", fileSplit[1]
	case reLocDep.MatchString(filename):
		return locDepSrc, fileSplit[1], ""
	case reLocDepInst.MatchString(filename):
		return locDepInstSrc, fileSplit[1], fileSplit[2]
	case filename == "env":
		return envSrc, "", ""
	case filename == "vault":
		return vaultSrc, "", ""
	case strings.Contains(filename, hostname) && hostname != "":
		filename = strings.TrimPrefix(filename, hostname+"-")

		switch {
		case reInst.MatchString(filename):
			r := reInst.FindStringSubmatch(filename)

			return hostInstSrc, "", r[1]
		case reWordInst.MatchString(filename):
			r := reWordInst.FindStringSubmatch(filename)

			return hostDepInstSrc, r[1], r[2]
		default:
			return hostDepSrc, filename, ""
		}
	case reWordInst.MatchString(filename):
		r := reWordInst.FindStringSubmatch(filename)

		return depInstSrc, r[1], r[2]
	}

	return depSrc, filename, ""
}
