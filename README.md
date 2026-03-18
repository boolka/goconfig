# Configure your go applications with *goconfig*

## The purpose

The main purpose of *goconfig* is to provide hierarchical configurations to go applications. It lets you define a set of default parameters, and extend them for different deployment environments (development, qa, staging, production, etc.) or external sources (vault).

For the library user, the idea is to treat the configuration directory as a "black box" and work with it as a whole, rather than using individual files.

## Install

```bash
go get github.com/boolka/goconfig@latest
```

Use `github.com/boolka/goconfig/pkg/config` package:
- `New` function to create new config instance
- `Options` struct to pass config options

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
	"context"
	"fmt"

	"github.com/boolka/goconfig"
)

func main() {
	ctx := context.Background()
	cfg, _ := goconfig.New(ctx, goconfig.Options{})

	syslogDomain := cfg.MustGet(ctx, "syslog-ng.domain")
	syslogPort, _ := cfg.Get(ctx, "syslog-ng.port")

	fmt.Printf("%s:%d", syslogDomain, syslogPort) // "syslog-ng:601"
}
```

Run it:

```bash
go run ./cfg.go
```

Suppose now we want to load environment variable. Create another configuration file in `./config` directory. Name it `env.toml` and copy the contents:

```toml
[syslog-ng]
domain = 'SYSLOG_NG_DOMAIN'
```

In that case if "SYSLOG_NG_DOMAIN" environment contains non empty value, then it will be loaded when we call `cfg.Get(ctx, "syslog-ng.domain")`. Try it:

```bash
SYSLOG_NG_DOMAIN=custom-syslog-ng-domain go run ./cfg.go
```

## Deeper

### API

- `New(context.Context, config.Options) (*config.Config, error)`. Creates new config instance. Provide `config.Options` object to set config path and etc. If configuration directory is empty the `ErrEmptyDir` sentinel error will be returned.
- `(*config.Config) Get(context.Context, path string, files ...string) (any, bool)`. Get method takes dot-delimited configuration path and returns a value if any. The last parameter specifies which files to search, with or without extension. If omitted, all files will be search through. The sequence of passed files does not change the search order. Second returned value states if it was found and follows comma ok idiom.
- `(*config.Config) MustGet(context.Context, path string, files ...string) any`. MustGet method is the same as Get except that it panics if the path does not exist.

#### Config options

Options specify configuration directory, instance, deployment, hostname and vault settings. All options are optional and can be omitted.

```go
goconfig.Options{
	Directory:         "/path/to/directory",       // trying to load from ./config by default
	FileSystem: 	   fs.ReadDirFS,               // useful in case of configuration embed
	Instance:          "1",
	Deployment:        "production",
	Hostname:          "localhost",                // os.Hostname() by default
	Logger:            *slog.Logger,               // goconfig will remain silent when nil is received
	VaultClient:       *vault.Client,              // vault client instance
}
```

##### Directory

Directory can be set by `Directory` option explicitly or implicitly via `GO_CONFIG_PATH` environment variable and must contain `.json`, `.yaml` (`.yml`) or `.toml` configuration files. All other files will be ignored. You can provide multiple directories delimited by `os.PathListSeparator`. Think of it as if you were putting all files together into one directory.

##### FileSystem

You can specify `fs.ReadDirFS` interface to restrict file system access. Can be useful to embed configuration. If specified then `Directory` option is treated like path relative to `FileSystem`. Can be omitted.

##### Deployment

Deployment can be set by `Deployment` option explicitly or implicitly via `GO_DEPLOYMENT` environment variable. For example "testing", "development" or "production" is common used deployment types. There is no default value. So you need to provide it somehow. If it is omitted then all deployment configuration files will be ignored.

##### Instance

For support multi instance configuration use `Instance` option. Can also be implicitly accepted via `GO_INSTANCE` environment variable. Instance value can only be a number. Meaning "default-1.toml" is valid instance file configuration, but "custom-instance.toml" is not.

Multi instance configuration common usage is for get specific options for horizontal scaled multi pod environments. Suppose we have workers set "worker-1", "worker-2" ... "worker-n". By the multi instance configurations you can provide specific options for every single worker.

##### Hostname

If you do not provide `Hostname` explicitly then `os.Hostname()` is called and the part after the first dot stripped off. For example suppose the `MacBook-Pro-5.local` is hostname. Then `Hostname` will borrow `MacBook-Pro-5`. It may be identical to the `hostname -s` call.

`Hostname` must not contain dots in general. This is important for searching a specific value in the configuration, because the dot is a field separator. Choose to provide it explicitly when in doubt.

##### Logger

Produce output to supplied logger. Module will be silent if nil was received. Can be helpful for state some source errors. For example if vault was unavailable then logger will receive message describes whats going on.

##### VaultClient

To use vault abilities you must use the `goconfig_vault` build tag and pass the vault client through the `VaultClient` option. Unauthorized client will lead to the runtime errors. For more details look at [Vault](####Vault) section below.

### Configuration files

Application configuration is stored in `.json`, `.yaml` (`.yml`) or `.toml` files. Other files will be ignored. Special case is the `env.EXT` ([Environment](####Environment)) and `vault.EXT` ([Vault](####Vault)) files.

#### Configuration files and field lookup order

When looking up a value using the `Get` or `MustGet` method of a configuration, the sources(files) in the configuration directory(ies) are searched in the following order (from highest to lowest):

- vault.EXT
- env.EXT
- local-{deployment}-{instance}.EXT
- local-{deployment}.EXT
- local-{instance}.EXT
- local.EXT
- {hostname}-{deployment}-{instance}.EXT
- {hostname}-{deployment}.EXT
- {hostname}-{instance}.EXT
- {hostname}.EXT
- {deployment}-{instance}.EXT
- {deployment}.EXT
- default-{instance}.EXT
- default.EXT

Where:

- EXT can be `.yaml` (`.yml`), `.json` or `.toml`
- {instance} is an optional instance name string for multi-instance deployments
- {hostname} is your hostname (don't use dots)
- {deployment} is the deployment name
- `env.EXT` and `vault.EXT` has special meanings and will be explained below

If you don't specify deployment, instance or hostname then the corresponding files will be ignored. All files with unknown filename signature will be treated as {deployment}.EXT and will be ignored if the deployment option is not provided. Dot prefixed files will be ignored.

#### Local files

`local.EXT`, `local-{deployment}-{instance}.EXT`, `local-{deployment}.EXT`, `local-{instance}.EXT` files are intended for use locally and to not be tracked in your version control system. Use it to overlap some definitions in local development for example.

#### Environment

The `env.EXT` file contains environment variable names that will be mapped to your configuration structure. For example, the file `env.toml`:

```toml
[server]
port = "SERVER_PORT"
```

defines that we will try to load `SERVER_PORT` environment variable value into port configuration field. That means the expression

```go
cfg.MustGet(ctx, "server.port").(string) == os.Getenv("SERVER_PORT")
```

will match. Loading environment variables is dynamic. *goconfig* will not save values while the configuration module is initializing. That means that if the runtime changes environment variable while the application is running then this value will be loaded.

Environment file may be any supported file extension - `.json`, `.yaml` (`.yml`) or `.toml`.

#### Vault

Vault file will have special meaning only if you specify `goconfig_vault` build tag during compilation. Otherwise it would be treated as plain text file and field values loaded as they are.

If you create `vault.EXT` file then fields from that file will be looked up from the vault server. When *goconfig* instance is created then vault domain will be checked up.

Loading vault variables is dynamic just like the environment variables. The field has special syntax `mount_path,secret_path,secret_key` where the vault secret mount path, secret path, and secret key separated by comma. Suppose we have a vault server with a configured postgresql secret and running at `http://localhost:8200`. Create directory `config` and place `vault.toml` file in to it with contents:

```toml
[postgresql]
username = "secret,postgresql,username"
password = "secret,postgresql,password"
```

Then create go file to upload keys:

```go
package main

import (
	"context"
	"fmt"

	"github.com/boolka/goconfig"
	vaultApi "github.com/hashicorp/vault/api"
)

func main() {
	ctx := context.Background()

	vaultCfg := vaultApi.DefaultConfig()
	vaultCfg.Address = "http://localhost:8200"

	client, err := vaultApi.NewClient(vaultCfg)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken("vault_token")

	cfg, _ := goconfig.New(ctx, goconfig.Options{
		VaultClient: client,
	})

	username := cfg.MustGet(ctx, "postgresql.username")
	password := cfg.MustGet(ctx, "postgresql.password")

	fmt.Printf("username: %s, password: %s", username, password)
}
```

If you create correct vault secret with keys `username` and `password` and provide access then these values was printed.

To check that certain field available direct from vault source you need to specify files argument:

```go
if password, ok := cfg.Get(ctx, "postgresql.password", "vault"); ok {
	// available from vault source
}
```

Managing vault auth methods, policies and secrets is out of scope.

### Config embedding

Suppose we want to embed configuration. We have `./config` directory containing our files. Embed the files first:

```go
import "embed"

//go:embed config/*
var configDir embed.FS
```

Specify `FileSystem` option and corresponding `Directory`:

```go
cfg, err := goconfig.New(ctx, config.Options{
	FileSystem: &configDir,
	Directory:  "config",
})
```

You must specify `Directory` option explicitly when `FileSystem` is given even if the configuration files reside in the default `config` directory.

## Cli usage

*goconfig* can be used to load config fields in the terminal.

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

Execute `goconfig --help` for more info.

## Under the hood

*goconfig* loads configuration files at instance creation time. If the configuration has changed, you will need to create another instance to see them.

The following libraries are used to load concrete configuration files:

- json: internal go library `encoding/json`
- toml: `github.com/pelletier/go-toml/v2`
- yaml: `gopkg.in/yaml.v3`

*goconfig* handles a variety of number types that specific serialization libraries provide, by normalizing them. Type depends on number magnitude range:

- x < MinInt or x > MaxUint: `float64`
- x > MaxInt and x <= MaxUint: `uint`
- x >= MinInt and x <= MaxInt: `int`

## Links

- Inspired with great library for node.js ecosystem [node-config](https://github.com/node-config/node-config)
- [Vault](https://github.com/hashicorp/vault) api
- [toml](https://github.com/pelletier/go-toml) format
- [yaml](https://gopkg.in/yaml.v3) format
