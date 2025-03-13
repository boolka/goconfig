package config

import (
	"testing"
)

func TestFileName(t *testing.T) {
	t.Parallel()

	if fileName("hostname-development-1.json") != "hostname-development-1" {
		t.Fatal("hostname-development-1.json")
	}

	if fileName("hostname-development-1.toml.json") != "hostname-development-1.toml" {
		t.Fatal("hostname-development-1.toml.json")
	}
}

func TestFileParse(t *testing.T) {
	t.Parallel()

	if s, _, _ := fileSource("default", ""); s != def {
		t.Fatal("default")
	}

	if s, _, i := fileSource("default-1", ""); s != defInst || i != "1" {
		t.Fatal("default-1")
	}

	if s, d, _ := fileSource("production", ""); s != dep || d != "production" {
		t.Fatal("production")
	}

	if s, _, i := fileSource("development-1", ""); s != depInst || i != "1" {
		t.Fatal("development-1")
	}

	if s, _, _ := fileSource("hostname", "hostname"); s != host {
		t.Fatal("hostname")
	}

	if s, _, i := fileSource("hostname-1", "hostname"); s != hostInst || i != "1" {
		t.Fatal("hostname-1")
	}

	if s, d, _ := fileSource("hostname-production", "hostname"); s != hostDep || d != "production" {
		t.Fatal("hostname-production")
	}

	if s, d, i := fileSource("hostname-development-1", "hostname"); s != hostDepInst || d != "development" || i != "1" {
		t.Fatal("hostname-development-1")
	}

	if s, _, _ := fileSource("local", ""); s != loc {
		t.Fatal("local")
	}

	if s, _, i := fileSource("local-1", ""); s != locInst || i != "1" {
		t.Fatal("local-1")
	}

	if s, d, _ := fileSource("local-production", ""); s != locDep || d != "production" {
		t.Fatal("local-production")
	}

	if s, d, i := fileSource("local-production-1", ""); s != locDepInst || d != "production" || i != "1" {
		t.Fatal("local-production-1")
	}

	if s, _, _ := fileSource("env", ""); s != env {
		t.Fatal("env")
	}

	if s, d, i := fileSource("unexpected", ""); s != dep && d == "" && i == "" {
		t.Fatal("unexpected")
	}

	if s, d, i := fileSource("unexpected-source", ""); s != dep && d == "" && i == "" {
		t.Fatal("unexpected-source")
	}

	if s, d, i := fileSource("unexpected-source-1", ""); s != dep && d == "" && i == "" {
		t.Fatal("unexpected-source-1")
	}

	if s, d, i := fileSource("unexpected-source-1-json", ""); s != dep && d == "" && i == "" {
		t.Fatal("unexpected-source-1-json")
	}
}
