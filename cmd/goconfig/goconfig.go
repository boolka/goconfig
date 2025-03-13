package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/boolka/goconfig/pkg/config"
)

const helpMsg = `Load application structured configuration. For more info follow https://github.com/boolka/goconfig.

--config (-c) sets configuration directory (default is ./config)
--deployment (-d) sets current deployment
--instance (-i) sets current instance
--hostname sets current hostname (by default will try to load os.Hostname() with the part after the first dot stripped off)
--get (-g) configuration path to lookup
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
			}
		}
	}

	cfg, err := config.New(config.Options{
		Directory:  configDirectory,
		Instance:   instance,
		Hostname:   hostname,
		Deployment: deployment,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	v, ok := cfg.Get(getPath)

	if !ok {
		fmt.Fprintln(os.Stderr, "\""+getPath+"\" key was not found")
	} else {
		fmt.Print(v)
	}
}
