package config

import "regexp"

type cfgSource int

const (
	_ cfgSource = iota
	def
	defInst
	dep
	depInst
	host
	hostInst
	hostDep
	hostDepInst
	loc
	locInst
	locDep
	locDepInst
	env
)

var reDefInst = regexp.MustCompile(`^default-\d+$`)
var reLocInst = regexp.MustCompile(`^local-\d+$`)
var reLocDep = regexp.MustCompile(`^local-\w+$`)
var reLocDepInst = regexp.MustCompile(`^local-\w+-\d+$`)

var reInst = regexp.MustCompile(`^(\d+)$`)
var reWord = regexp.MustCompile(`^(.+)$`)
var reWordInst = regexp.MustCompile(`^(.+)-(\d+)$`)
