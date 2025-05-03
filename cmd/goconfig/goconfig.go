package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/boolka/goconfig"
	vaultMock "github.com/boolka/goconfig/pkg/vault"
)

const helpMsg = `Load application structured configuration. For more info follow https://github.com/boolka/goconfig.

--config (-c) sets configuration directory (default is ./config)
--deployment (-d) sets current deployment
--instance (-i) sets current instance
--hostname sets current hostname (by default will try to load os.Hostname() with the part after the first dot stripped off)
--get (-g) configuration path to lookup
--token set vault token
--verbose (-v) add debug and errors output
--help (-h) prints this message
`

func main() {
	var configDirectory = ""
	var deployment = ""
	var instance = ""
	var hostname = ""
	var getPath = ""
	var configDirectoryArg = false
	var deploymentArg = false
	var instanceArg = false
	var hostnameArg = false
	var getArg = false
	var vaultTokenArg = false
	var verbose = false
	var vaultToken = ""

	for _, arg := range os.Args {
		switch {
		case configDirectoryArg:
			configDirectory = arg
			configDirectoryArg = false
			continue
		case deploymentArg:
			deployment = arg
			deploymentArg = false
		case instanceArg:
			instance = arg
			instanceArg = false
		case hostnameArg:
			hostname = arg
			hostnameArg = false
		case getArg:
			getPath = arg
			getArg = false
			continue
		case vaultTokenArg:
			vaultToken = arg
			vaultTokenArg = false
			continue
		}

		switch arg {
		case "--config", "-c":
			configDirectoryArg = true
			continue
		case "--deployment", "-d":
			deploymentArg = true
			continue
		case "--instance", "-i":
			instanceArg = true
			continue
		case "--hostname":
			hostnameArg = true
			continue
		case "--get", "-g":
			getArg = true
			continue
		case "--verbose", "-v":
			verbose = true
		case "--token":
			vaultTokenArg = true
		case "--help", "-h":
			fmt.Print(helpMsg)
			return
		default:
			if !strings.Contains(arg, "=") {
				continue
			}

			param := strings.Split(arg, "=")[1]

			switch {
			case strings.Contains(arg, "--config"):
				configDirectory = param
			case strings.Contains(arg, "--deployment"):
				deployment = param
			case strings.Contains(arg, "--instance"):
				instance = param
			case strings.Contains(arg, "--hostname"):
				hostname = param
			case strings.Contains(arg, "--get"):
				getPath = param
			case strings.Contains(arg, "--token"):
				vaultToken = param
			}
		}
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

	if vaultToken != "" {
		opts.VaultAuth = vaultMock.NewTokenAuth(vaultToken)
	}

	cfg, err := goconfig.New(ctx, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	v, ok := cfg.Get(ctx, getPath)
	if !ok {
		fmt.Fprintln(os.Stderr, "\""+getPath+"\" key not found")
	} else {
		fmt.Print(v)
	}
}
