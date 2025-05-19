# Configure your go applications with *goconfig*

## The purpose

The *goconfig* main purpose is to provide hierarchical configurations to go applications. It lets you define a set of default parameters, and extend them for different deployment environments (development, qa, staging, production, etc.) or external sources (vault).

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

	syslogDomain, _ := cfg.Get(ctx, "syslog-ng.domain")
	syslogPort, _ := cfg.Get(ctx, "syslog-ng.port")

	fmt.Printf("%s:%d", syslogDomain, syslogPort) // "syslog-ng:601"
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

In that case if "SYSLOG_NG_DOMAIN" environment contains non empty value then it will be loaded when we call `cfg.Get(ctx, "syslog-ng.domain")`. Try it:

```bash
SYSLOG_NG_DOMAIN=custom-syslog-ng-domain go run ./cfg.go
```

## Deeper

### Api

- `New(context.Context, config.Options) (*config.Config, error)`. Creates new config instance. Provide `config.Options` object to set config path and etc. If configuration directory is empty the `ErrEmptyDir` sentinel error will be returned.
- `(*config.Config) Get(context.Context, path string, files ...string) (any, bool)`. Get method takes dot delimited configuration path and returns value if any. The last parameter specifies which files to allow for searching both with or without extension. If omitted, all files will be search through. The sequence of transmitted files does not change the original order for searching. Second returned value states if it was found and follows comma ok idiom at all.
- `(*config.Config) MustGet(context.Context, path string, files ...string) any`. MustGet method is the same as Get except that it panics if the path does not exist.
- `(*config.Config) GetVaultClient() *vault.Client`. Returns vault client created and configured or provided directly by option.

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
	VaultClient:       *vault.Client,              // vault client if you don't want to create a new one
	VaultAuth:         vault.AuthMethod,           // is the AuthMethod interface from `github.com/hashicorp/vault/api` module that provides Login method
}
```

##### Directory

Directory can be set by `Directory` option explicitly or implicitly via `GO_CONFIG_PATH` environment variable and must contains `.json`, `.yaml` (`.yml`) or `.toml` configuration files. All other files will be ignored. You can provide multiple directories delimited by `os.PathListSeparator`. Think of it as if you were putting all files together into one directory.

##### FileSystem

You can specify `fs.ReadDirFS` interface to restrict file system access. Can be useful to embed configuration. If specified then `Directory` option treated like path relative to `FileSystem`. Can be omitted.

##### Deployment

Deployment can be set by `Deployment` option explicitly or implicitly via `GO_DEPLOYMENT` environment variable. For example "testing", "development" or "production" is common used deployment types. There is no default value. So you need to provide it somehow. If it omitted then all deployment configuration files will be ignored.

##### Instance

For support multi instance configuration use `Instance` option. Can also be implicitly accepted via `GO_INSTANCE` environment variable. Instance value can only be a number. Mean "default-1.toml" is valid instance file configuration, but "custom-instance.toml" is not and will be skipped. If you do not specify instance at all then every instance suffix file will be ignored.

Multi instance configuration common usage is for get specific options for horizontal scaled multi pod environments. Suppose we have workers set "worker-1", "worker-2" ... "worker-n". By the multi instance configurations you can provide specific options for every single worker.

##### Hostname

If you do not provide `Hostname` explicitly then `os.Hostname()` is called with the part after the first dot stripped off. For example suppose the `MacBook-Pro-5.local` is hostname. Then `Hostname` will borrow `MacBook-Pro-5`. It may be identical with `hostname -s` call.

`Hostname` must not contain dots in general. This is important for searching a specific value in the configuration, as long as the dot is a field separator. Choose to provide it explicitly when in doubt.

##### Logger

Produce output to supplied logger. Module will be silent if nil was received. Can be helpful for state some source errors. For example if vault was unavailable then logger will receive message describes whats going on.

##### VaultClient

Pass your own vault client through the `VaultClient` option if you don't want to create a new one. It may already have the token or will be authorized. Client can be configured by `goconfig.vault` path and authorized by `goconfig.vault.auth` path declared in any of configuration files. See below in [Vault](####Vault) section for more details.

##### VaultAuth

`VaultAuth` is AuthMethod interface from `github.com/hashicorp/vault/api` module that provides Login method. You can use one of vault auth methods or create your own. Visit on [vault github page](https://github.com/hashicorp/vault) and lookup for `api/auth` subpath contains various vault auth methods. For example `github.com/hashicorp/vault/api/auth/aws` module provides vault aws auth.

### Configuration files

Application configuration persists in any of `.json`, `.yaml` (`.yml`) or `.toml` files. Other files will be ignored. Special case is the `env.EXT` ([Environment](####Environment)) and `vault.EXT` ([Vault](####Vault)) files.

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

If you don't specify deployment, instance or hostname then the correspond files will be ignored. All files with unknown filename signature will be treated as {deployment}.EXT and will be ignored if no exact deployment option.

#### Local files

`local.EXT`, `local-{deployment}-{instance}.EXT`, `local-{deployment}.EXT`, `local-{instance}.EXT`, `local.EXT` files are intended to use locally and to not be tracked in your version control system. Use it to overlap some definitions in local development for example.

#### Environment

The `env.EXT` file contains environment variable names that will be mapped to your configuration structure. For example file `env.toml`:

```toml
[server]
port = "SERVER_PORT"
```

defines that we will try to load `SERVER_PORT` environment variable value to port configuration field. That means the expression

```go
cfg.MustGet(ctx, "server.port").(string) == os.Getenv("SERVER_PORT")
```

will match. Loading environment variables is dynamic. *goconfig* will not save values while configuration module initializing. That means that if the runtime spoof environment variable while application running than this value will be loaded.

Environment file may be any supported file extension - `.json`, `.yaml` (`.yml`) or `.toml`.

#### Vault

If you create `vault.EXT` file then fields from that file will proceed to lookup from vault server. When *goconfig* instance was created then vault domain will be checkup. If you want to disable vault lookup on some configurations then use `goconfig.vault.enable` system option. All vault system options will be explained below.

Loading vault variables is dynamic too both the same with environment. The field has special syntax `mount_path,secret_path,secret_key` where are listed vault secret mount path, secret path, and secret key separated by comma. Suppose we have vault server with configured postgresql secret and started on `http://localhost:8200`. Lets load for example `username` and `password` keys to our config. And suppose for our simple example that vault setup has the `root` token to access. All again: mount - secret, path - postgresql, keys - username, password. Create directory `config` and place `vault.toml` file in to it with contents:

```toml
[postgresql]
username = "secret,postgresql,username"
password = "secret,postgresql,password"

[goconfig.vault]
address = "http://localhost:8200"

[goconfig.vault.auth]
token = "root"
```

Then create go file to upload keys:

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

	username, _ := cfg.Get(ctx, "postgresql.username")
	password, _ := cfg.Get(ctx, "postgresql.password")

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

##### Vault Configuration

So, before we can use vault keys we must configure vault client. There are two ways:

- Provide already configured vault client though `VaultClient` option. In this case client may be already authorized. Otherwise specify `VaultAuth` or provide config files with `[goconfig.vault.auth]` section.
- Through [goconfig.vault] section from configuration files. It does't important in what source (`default.EXT`, `local.EXT` and etc) it would be written. Specify address and optionally retry/timeout params. For example create `default.toml` file:

```toml
[goconfig.vault]
address = "http://127.0.0.1:8200"
enable = false
min_retry_wait = "3"
max_retry_wait = "5"
max_retries = "10"
timeout = "30"
```

You can specify optional `enable` param to turn off vault on specific configuration environment.

Params `min_retry_wait`, `max_retry_wait`, `max_retries` and `timeout` must be specified as strings. Their values treated as seconds when specified without postfix. Otherwise the value will be supplied to `time.ParseDuration` function.

##### Vault Authorization

Vault authorization can be achieved by `VaultAuth` option or by specified fields under path `goconfig.vault.auth.*` in configuration files. `VaultAuth` option have precedence over file configurations. As explained above [VaultAuth](#####VaultAuth) is just AuthMethod interface defined in `github.com/hashicorp/vault/api` module.

Vault file auth configuration includes:

- token authorization. `goconfig.vault.auth.token` field must be supplied, for example:

```toml
[goconfig.vault.auth]
token = "root"
```

- userpass authorization. `goconfig.vault.auth.username` and `goconfig.vault.auth.password` fields must be supplied, for example:

```toml
[goconfig.vault.auth]
username = "vault_username"
password = "qwerty123456"
```

- approle authorization: `goconfig.vault.auth.roleid` and `goconfig.vault.auth.secretid` fields must be supplied, for example:

```toml
[goconfig.vault.auth]
roleid = "db02de05-fa39-4855-059b-67221c5c2f63"
secretid = "6a174c20-f6de-a53c-74d2-6018fcceff64"
```

It is not acceptable in many cases to provide credentials directly to file. So you can provide essential info by environment variables loaded via `env.toml`:

```toml
[goconfig.vault.auth]
token = "CUSTOM_VAULT_TOKEN"
```

And so on for other cases.

### Config embedding

Suppose we want to embed configuration. We have `./config` directory contains our files. Embed the files first:

```go
import "embed"

//go:embed config/*
var configDir embed.FS
```

Specify `FileSystem` option and correspond `Directory`:

```go
cfg, err := goconfig.New(ctx, config.Options{
	FileSystem: &configDir,
	Directory:  "config",
})
```

You must specify `Directory` option explicitly when `FileSystem` is given even if configuration files persists in treated as default `config` directory.

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

Execute `goconfig --help` for more info.

## Under the hood

*goconfig* loads configuration files at instance creation time. If the configuration has changed, you will need to create another instance to see them.

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
- [Vault](https://github.com/hashicorp/vault) api
- [toml](https://github.com/pelletier/go-toml) format
- [yaml](https://gopkg.in/yaml.v3) format
