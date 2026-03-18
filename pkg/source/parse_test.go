package source_test

import (
	"testing"

	"github.com/boolka/goconfig/pkg/source"
)

func TestParseFilenameToSourceType(t *testing.T) {
	t.Parallel()

	if s, _, _ := source.ParseFilename("default", ""); s != source.DefSrc {
		t.Fatal("default")
	}

	if s, _, i := source.ParseFilename("default-1", ""); s != source.DefInstSrc || i != "1" {
		t.Fatal("default-1")
	}

	if s, d, _ := source.ParseFilename("production", ""); s != source.DepSrc || d != "production" {
		t.Fatal("production")
	}

	if s, _, i := source.ParseFilename("development-1", ""); s != source.DepInstSrc || i != "1" {
		t.Fatal("development-1")
	}

	if s, _, _ := source.ParseFilename("hostname", "hostname"); s != source.HostSrc {
		t.Fatal("hostname")
	}

	if s, _, i := source.ParseFilename("hostname-1", "hostname"); s != source.HostInstSrc || i != "1" {
		t.Fatal("hostname-1")
	}

	if s, d, _ := source.ParseFilename("hostname-production", "hostname"); s != source.HostDepSrc || d != "production" {
		t.Fatal("hostname-production")
	}

	if s, d, i := source.ParseFilename("hostname-development-1", "hostname"); s != source.HostDepInstSrc || d != "development" || i != "1" {
		t.Fatal("hostname-development-1")
	}

	if s, _, _ := source.ParseFilename("local", ""); s != source.LocSrc {
		t.Fatal("local")
	}

	if s, _, i := source.ParseFilename("local-1", ""); s != source.LocInstSrc || i != "1" {
		t.Fatal("local-1")
	}

	if s, d, _ := source.ParseFilename("local-production", ""); s != source.LocDepSrc || d != "production" {
		t.Fatal("local-production")
	}

	if s, d, i := source.ParseFilename("local-production-1", ""); s != source.LocDepInstSrc || d != "production" || i != "1" {
		t.Fatal("local-production-1")
	}

	if s, d, i := source.ParseFilename("unexpected", ""); s != source.DepSrc && d == "" && i == "" {
		t.Fatal("unexpected")
	}

	if s, d, i := source.ParseFilename("unexpected-source", ""); s != source.DepSrc && d == "" && i == "" {
		t.Fatal("unexpected-source")
	}

	if s, d, i := source.ParseFilename("unexpected-source-1", ""); s != source.DepSrc && d == "" && i == "" {
		t.Fatal("unexpected-source-1")
	}

	if s, d, i := source.ParseFilename("unexpected-source-1-json", ""); s != source.DepSrc && d == "" && i == "" {
		t.Fatal("unexpected-source-1-json")
	}
}
