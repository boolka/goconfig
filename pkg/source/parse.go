package source

import (
	"regexp"
	"strings"
)

var reDefInst = regexp.MustCompile(`^default-\d+$`)
var reLocInst = regexp.MustCompile(`^local-\d+$`)
var reLocDep = regexp.MustCompile(`^local-\w+$`)
var reLocDepInst = regexp.MustCompile(`^local-\w+-\d+$`)

var reInst = regexp.MustCompile(`^(\d+)$`)
var reWordInst = regexp.MustCompile(`^(.+)-(\d+)$`)

// returns source type, deployment, instance
func ParseFilename(filename, hostname string) (SourceType, string, string) {
	if filename == hostname {
		return HostSrc, "", ""
	}

	fileSplit := strings.Split(filename, "-")

	switch {
	case filename == "env":
		return EnvSrc, "", ""
	case filename == "vault":
		return VaultSrc, "", ""
	case filename == "default":
		return DefSrc, "", ""
	case reDefInst.MatchString(filename):
		return DefInstSrc, "", fileSplit[1]
	case filename == "local":
		return LocSrc, "", ""
	case reLocInst.MatchString(filename):
		return LocInstSrc, "", fileSplit[1]
	case reLocDep.MatchString(filename):
		return LocDepSrc, fileSplit[1], ""
	case reLocDepInst.MatchString(filename):
		return LocDepInstSrc, fileSplit[1], fileSplit[2]
	case strings.Contains(filename, hostname) && hostname != "":
		filename = strings.TrimPrefix(filename, hostname+"-")

		switch {
		case reInst.MatchString(filename):
			r := reInst.FindStringSubmatch(filename)

			return HostInstSrc, "", r[1]
		case reWordInst.MatchString(filename):
			r := reWordInst.FindStringSubmatch(filename)

			return HostDepInstSrc, r[1], r[2]
		default:
			return HostDepSrc, filename, ""
		}
	case reWordInst.MatchString(filename):
		r := reWordInst.FindStringSubmatch(filename)

		return DepInstSrc, r[1], r[2]
	}

	return DepSrc, filename, ""
}
