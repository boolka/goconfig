package config_test

import (
	"fmt"
	"testing"

	"github.com/boolka/goconfig/pkg/config"
)

func TestJsonNumber(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/serializers/json",
	})

	if err != nil {
		t.Fatal(err)
	}

	syslogDomain, ok := cfg.Get("syslog-ng.domain")

	if !ok {
		t.Fatal(syslogDomain, ok)
	}

	syslogPort, ok := cfg.Get("syslog-ng.port")

	if !ok {
		t.Fatal(syslogPort, ok)
	}

	if fmt.Sprintf("%s:%.0f", syslogDomain, syslogPort) != "syslog-ng:601" {
		t.Fatal(syslogDomain, syslogPort)
	}
}

func TestTomlNumber(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/serializers/toml",
	})

	if err != nil {
		t.Fatal(err)
	}

	syslogDomain, ok := cfg.Get("syslog-ng.domain")

	if !ok {
		t.Fatal(syslogDomain, ok)
	}

	syslogPort, ok := cfg.Get("syslog-ng.port")

	if !ok {
		t.Fatal(syslogPort, ok)
	}

	if fmt.Sprintf("%s:%d", syslogDomain, syslogPort) != "syslog-ng:601" {
		t.Fatal(syslogDomain, syslogPort)
	}
}

func TestYamlNumber(t *testing.T) {
	t.Parallel()

	cfg, err := config.New(config.Options{
		Directory: "./testdata/serializers/yaml",
	})

	if err != nil {
		t.Fatal(err)
	}

	syslogDomain, ok := cfg.Get("syslog-ng.domain")

	if !ok {
		t.Fatal(syslogDomain, ok)
	}

	syslogPort, ok := cfg.Get("syslog-ng.port")

	if !ok {
		t.Fatal(syslogPort, ok)
	}

	if fmt.Sprintf("%s:%d", syslogDomain, syslogPort) != "syslog-ng:601" {
		t.Fatal(syslogDomain, syslogPort)
	}
}
