package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/boolka/goconfig"
)

const helpMsg = `Load application structured configuration. For more info follow https://github.com/boolka/goconfig.

--config (-c) sets configuration directory (default is ./config)
--deployment (-d) sets current deployment
--instance (-i) sets current instance
--hostname sets current hostname (by default will try to load os.Hostname() with the part after the first dot stripped off)
--get (-g) configuration path to lookup
--verbose (-v) add debug and errors output
--help (-h) prints this message
`

func main() {
	var configDirectory, deployment, instance, hostname, getPath string
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

	cfg, err := goconfig.New(ctx, goconfig.Options{
		Directory:  configDirectory,
		Instance:   instance,
		Hostname:   hostname,
		Deployment: deployment,
		Logger:     logger,
	})
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
