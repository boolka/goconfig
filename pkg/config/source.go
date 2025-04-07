package config

import "regexp"

type cfgSource int

const (
	_ cfgSource = iota
	defSrc
	defInstSrc
	depSrc
	depInstSrc
	hostSrc
	hostInstSrc
	hostDepSrc
	hostDepInstSrc
	locSrc
	locInstSrc
	locDepSrc
	locDepInstSrc
	envSrc
	vaultSrc
)

var reDefInst = regexp.MustCompile(`^default-\d+$`)
var reLocInst = regexp.MustCompile(`^local-\d+$`)
var reLocDep = regexp.MustCompile(`^local-\w+$`)
var reLocDepInst = regexp.MustCompile(`^local-\w+-\d+$`)

var reInst = regexp.MustCompile(`^(\d+)$`)
var reWord = regexp.MustCompile(`^(.+)$`)
var reWordInst = regexp.MustCompile(`^(.+)-(\d+)$`)
