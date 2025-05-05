package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/boolka/goconfig"
	vaultMock "github.com/boolka/goconfig/pkg/vault"
	appRoleAuth "github.com/hashicorp/vault/api/auth/approle"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

const helpMsg = `Load application structured configuration. For more info follow https://github.com/boolka/goconfig.

--config (-c) sets configuration directory (default is ./config)
--deployment (-d) sets current deployment
--instance (-i) sets current instance
--hostname sets current hostname (by default will try to load os.Hostname() with the part after the first dot stripped off)
--get (-g) configuration path to lookup
--token vault token. For token auth
--username vault userpass auth username
--password vault userpass auth password
--roleid vault approle auth roleid
--secretid vault approle auth secretid
--verbose (-v) add debug and errors output
--help (-h) prints this message
`

func main() {
	var configDirectory, deployment, instance, hostname, getPath string
	var vaultToken, vaultUsername, vaultPassword, vaultRoleId, vaultSecretId string
	var verbose, help bool

	flag.StringVar(&configDirectory, "config", "", "provide optional configuration files directory")
	flag.StringVar(&configDirectory, "c", "", "provide optional configuration files directory")

	flag.StringVar(&deployment, "deployment", "", "provide optional configuration deployment")
	flag.StringVar(&deployment, "d", "", "provide optional configuration deployment")

	flag.StringVar(&instance, "instance", "", "provide optional configuration instance")
	flag.StringVar(&instance, "i", "", "provide optional configuration instance")

	flag.StringVar(&hostname, "hostname", "", "provide optional configuration hostname")

	flag.StringVar(&getPath, "get", "", "path to the configuration field")
	flag.StringVar(&getPath, "g", "", "path to the configuration field")

	flag.StringVar(&vaultToken, "token", "", "the vault token")
	flag.StringVar(&vaultUsername, "username", "", "the vault username (for userpass auth)")
	flag.StringVar(&vaultPassword, "password", "", "the vault password (for userpass auth)")
	flag.StringVar(&vaultRoleId, "roleid", "", "the vault roleid (for approle auth)")
	flag.StringVar(&vaultSecretId, "secretid", "", "the vault secretid (for approle auth)")

	flag.BoolVar(&verbose, "verbose", false, "provide optional configuration verbose option")
	flag.BoolVar(&verbose, "v", false, "provide optional configuration verbose option")

	flag.BoolVar(&help, "help", false, "help message")
	flag.BoolVar(&help, "h", false, "help message")

	flag.Parse()

	if help {
		fmt.Print(helpMsg)
		return
	}

	ctx := context.Background()

	var logger *slog.Logger

	if verbose {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}))
	}

	opts := goconfig.Options{
		Directory:  configDirectory,
		Instance:   instance,
		Hostname:   hostname,
		Deployment: deployment,
		Logger:     logger,
	}

	switch {
	case vaultToken != "":
		opts.VaultAuth = vaultMock.NewTokenAuth(vaultToken)
	case vaultUsername != "" && vaultPassword != "":
		auth, err := userPassAuth.NewUserpassAuth(vaultUsername, &userPassAuth.Password{
			FromString: vaultPassword,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		opts.VaultAuth = auth
	case vaultRoleId != "" && vaultSecretId != "":
		auth, err := appRoleAuth.NewAppRoleAuth(
			vaultRoleId,
			&appRoleAuth.SecretID{FromString: vaultSecretId},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		opts.VaultAuth = auth
	}

	cfg, err := goconfig.New(ctx, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	v, ok := cfg.Get(ctx, getPath)
	if !ok {
		fmt.Fprintln(os.Stderr, "\""+getPath+"\" key not found")
	} else {
		fmt.Print(v)
	}
}
