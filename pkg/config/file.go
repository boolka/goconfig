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
		return host, "", ""
	}

	fileSplit := strings.Split(filename, "-")

	switch {
	case filename == "default":
		return def, "", ""
	case reDefInst.MatchString(filename):
		return defInst, "", fileSplit[1]
	case filename == "local":
		return loc, "", ""
	case reLocInst.MatchString(filename):
		return locInst, "", fileSplit[1]
	case reLocDep.MatchString(filename):
		return locDep, fileSplit[1], ""
	case reLocDepInst.MatchString(filename):
		return locDepInst, fileSplit[1], fileSplit[2]
	case filename == "env":
		return env, "", ""
	case strings.Contains(filename, hostname) && hostname != "":
		filename = strings.TrimPrefix(filename, hostname+"-")

		switch {
		case reInst.MatchString(filename):
			r := reInst.FindStringSubmatch(filename)

			return hostInst, "", r[1]
		case reWordInst.MatchString(filename):
			r := reWordInst.FindStringSubmatch(filename)

			return hostDepInst, r[1], r[2]
		default:
			return hostDep, filename, ""
		}
	case reWordInst.MatchString(filename):
		r := reWordInst.FindStringSubmatch(filename)

		return depInst, r[1], r[2]
	}

	return dep, filename, ""
}
