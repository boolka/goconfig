# Configure your go applications with *goconfig*

## The purpose

The *goconfig* main purpose is to provide hierarchical configurations to go applications. It lets you define a set of default parameters, and extend them for different deployment environments (development, qa, staging, production, etc.).

## Quick start

Create a `./config` directory and place configuration files in it, for example a `default.toml` file with the following contents:

```toml
[syslog-ng]
domain = 'syslog-ng'
port = 601
```

Then create go file `cfg.go` to load appropriate configuration:

```go
package main

import (
	"fmt"

	"github.com/boolka/goconfig/pkg/config"
)

func main() {
	cfg, _ := config.New(config.Options{})

	syslogDomain, _ := cfg.Get("syslog-ng.domain")
	syslogPort, _ := cfg.Get("syslog-ng.port")

	syslog := fmt.Sprintf("%s:%d", syslogDomain, syslogPort)
	fmt.Println(syslog) // "syslog-ng:601"
}
```

Run it:

```bash
go run ./cfg.go
```

Suppose now we want to load environment variable. Create another one configuration file in `./config` directory. Name it `env.toml` and copy the contents:

```toml
[syslog-ng]
domain = 'SYSLOG_NG_DOMAIN'
```

In that case if "SYSLOG_NG_DOMAIN" environment contains non empty value then it will be loaded when we call `cfg.Get("syslog-ng.domain")`. Try it:

```bash
SYSLOG_NG_DOMAIN=custom-syslog-ng-domain go run ./cfg.go
```

## Deeper

### Api

- `New(config.Options) (*config.Config, error)`. Creates new config instance. Provide `config.Options` object to set config path and etc. This function call use recovery mechanism and will never panics.
- `(*config.Config) Get(path string) (any, bool)`. Get method takes dot delimited configuration path and returns value if any. Second returned value states if it was found and follows comma ok idiom at all. If this method is used with `nil` config receiver it will return `config.ErrUninitialized` sentinel error in first return value.

#### Config options

Options specify configuration directory, instance, deployment and hostname. All options can be omitted.

```go
config.Options{
    Directory:  "/path/to/directory", // trying to load from ./config by default
    Instance:   "1",
    Deployment: "production",
    Hostname:   "localhost",          // os.Hostname() by default
}
```

##### Directory

Directory can be set by `Directory` option explicitly or implicitly via `GO_CONFIG_PATH` environment variable and must contains `.json`, `.yaml` (`.yml`) or `.toml` configuration files. All other files will be ignored. You can provide multiple directories delimited by `os.PathListSeparator`. Think of it as if you were putting all files together into one directory.

##### Deployment

Deployment can be set by `Deployment` option explicitly or implicitly via `GO_DEPLOYMENT` environment variable. For example "testing", "development" or "production" is common used deployment types. There is no default value. So you need to provide it somehow. If it omitted then all deployment configuration files will be ignored.

##### Multi instance

For support multi instance configuration use `Instance` option. Can also be implicitly accepted via `GO_INSTANCE` environment variable. Instance value can only be a number. Mean "default-1.toml" is valid instance file configuration, but "custom-instance.toml" is not and will be skipped. If you do not specify instance at all then every instance suffix file will be ignored.

Multi instance configuration common usage is for get specific options for horizontal scaled multi pod environments. Suppose we have workers set "worker-1", "worker-2" ... "worker-n". By the multi instance configurations you can provide specific options for every single pod.

##### Hostname

If you do not provide `Hostname` explicitly then `os.Hostname()` is called with the part after the first dot stripped off. For example suppose the `MacBook-Pro-5.local` is hostname. Then `Hostname` will borrow `MacBook-Pro-5`. It may be identical with `hostname -s` call.

`Hostname` must not contain dots in general. This is important for searching a specific value in the configuration, as long as the dot is a field separator. Choose to provide it explicitly when in doubt.

### Configuration files

Application configuration persists in any of `.json`, `.yaml` (`.yml`) or `.toml` files. Other files will be ignored. Special case is `env.EXT` file. The `env.EXT` file contains environment variable names that will be loaded. For example file `env.toml`:

```toml
[server]
port = "SERVER_PORT"
```

defines that we will try to load `SERVER_PORT` environment variable value to port configuration field. That means the expression

```go
cfg.Get("server.port") == os.Getenv("SERVER_PORT")
```

will match. Loading environment variables is dynamic. *goconfig* will not save values while configuration module initializing. That means that if the runtime spoof environment variable while application running than this value will be loaded.

Environment file may be any supported file extension - `.json`, `.yaml` (`.yml`) or `.toml`.

#### File Load Order

Files in the config directory are loaded in the following order (from lowest to highest):

- default.EXT
- default-{instance}.EXT
- {deployment}.EXT
- {deployment}-{instance}.EXT
- {hostname}.EXT
- {hostname}-{instance}.EXT
- {hostname}-{deployment}.EXT
- {hostname}-{deployment}-{instance}.EXT
- local.EXT
- local-{instance}.EXT
- local-{deployment}.EXT
- local-{deployment}-{instance}.EXT
- env.EXT

Where:
- EXT can be `.yaml` (`.yml`), `.json` or `.toml`
- {instance} is an optional instance name string for multi-instance deployments
- {hostname} is your hostname. Don't use dots
- {deployment} is the deployment name

`default` has the lowest priority and `env` has the highest priority when searching for a specific configuration value.

`local` files are intended to use locally and to not be tracked in your version control system.

## Cli usage

*goconfig* can be used to load config fields in terminal.

First of all install it:

```bash
go install github.com/boolka/goconfig/cmd/goconfig@latest
```

Then for example create directory `./config` and place `default.toml` file in to it with contents:

```toml
delay = 1
```

Execute this example command:

```bash
goconfig --get delay | xargs sleep
```

Other options duplicate `config.Options`. Execute `goconfig --help` for more info.

## Under the hood

Following libraries used to load concrete configuration files:

- json: internal go library `encoding/json`
- toml: `github.com/pelletier/go-toml/v2`
- yaml: `gopkg.in/yaml.v3`

*goconfig* takes care over variety of number types specific serialization library provides by normalizing them. Type depends on number magnitude range:

- x < MinInt or x > MaxUint: `float64`
- x > MaxInt and x <= MaxUint: `uint`
- x >= MinInt and x <= MaxInt: `int`

## Links

- Inspired with great library for node.js ecosystem [node-config](https://github.com/node-config/node-config)
