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

	if s, _, _ := fileSource("default", ""); s != defSrc {
		t.Fatal("default")
	}

	if s, _, i := fileSource("default-1", ""); s != defInstSrc || i != "1" {
		t.Fatal("default-1")
	}

	if s, d, _ := fileSource("production", ""); s != depSrc || d != "production" {
		t.Fatal("production")
	}

	if s, _, i := fileSource("development-1", ""); s != depInstSrc || i != "1" {
		t.Fatal("development-1")
	}

	if s, _, _ := fileSource("hostname", "hostname"); s != hostSrc {
		t.Fatal("hostname")
	}

	if s, _, i := fileSource("hostname-1", "hostname"); s != hostInstSrc || i != "1" {
		t.Fatal("hostname-1")
	}

	if s, d, _ := fileSource("hostname-production", "hostname"); s != hostDepSrc || d != "production" {
		t.Fatal("hostname-production")
	}

	if s, d, i := fileSource("hostname-development-1", "hostname"); s != hostDepInstSrc || d != "development" || i != "1" {
		t.Fatal("hostname-development-1")
	}

	if s, _, _ := fileSource("local", ""); s != locSrc {
		t.Fatal("local")
	}

	if s, _, i := fileSource("local-1", ""); s != locInstSrc || i != "1" {
		t.Fatal("local-1")
	}

	if s, d, _ := fileSource("local-production", ""); s != locDepSrc || d != "production" {
		t.Fatal("local-production")
	}

	if s, d, i := fileSource("local-production-1", ""); s != locDepInstSrc || d != "production" || i != "1" {
		t.Fatal("local-production-1")
	}

	if s, _, _ := fileSource("env", ""); s != envSrc {
		t.Fatal("env")
	}

	if s, d, i := fileSource("unexpected", ""); s != depSrc && d == "" && i == "" {
		t.Fatal("unexpected")
	}

	if s, d, i := fileSource("unexpected-source", ""); s != depSrc && d == "" && i == "" {
		t.Fatal("unexpected-source")
	}

	if s, d, i := fileSource("unexpected-source-1", ""); s != depSrc && d == "" && i == "" {
		t.Fatal("unexpected-source-1")
	}

	if s, d, i := fileSource("unexpected-source-1-json", ""); s != depSrc && d == "" && i == "" {
		t.Fatal("unexpected-source-1-json")
	}
}
