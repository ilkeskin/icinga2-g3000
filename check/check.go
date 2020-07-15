package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const version = "0.0.1"

const (
	exitOk       = 0
	exitWarning  = 1
	exitCritical = 2
	exitUnknown  = 3
)

// GlobalReturnCode holds last issued exit code
var GlobalReturnCode = exitUnknown

func printVersion() {
	fmt.Println("check_g3000 v" + version)
	GlobalReturnCode = exitOk
}

func main() {
	app := &cli.App{
		Name:    "check_g3000",
		Usage:   "Check plugin to monitor a TDT G3000 gateway",
		Version: version,
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "uptime",
				Aliases:     []string{"u"},
				Usage:       "get device uptime (in s)",
				Description: "retrieves the device uptime since last (re)boot",
			},
			&cli.Command{
				Name:        "cpu",
				Aliases:     []string{"c"},
				Usage:       "get CPU usage (in %)",
				Description: "retrieves the current CPU usage as a percentage of total cpu time split between user, system and idle",
			},
			&cli.Command{
				Name:        "mem",
				Aliases:     []string{"m"},
				Usage:       "get memory usage (in %)",
				Description: "retrieves the current Memory usage as a percentage of total memory available split between used, cached and free RAM",
			},
			&cli.Command{
				Name:        "net",
				Aliases:     []string{"n"},
				Usage:       "get network usage (in kbps) of a NIC",
				Description: "",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "hostname",
				Aliases:     []string{"H"},
				Value:       "192.168.25.9",
				DefaultText: "192.168.25.9",
				Usage:       "Specifies the hostname or IP-address to query.",
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       5665,
				DefaultText: "5665",
				Usage:       "Specifies the port to query.",
			},
			&cli.StringFlag{
				Name:        "method",
				Aliases:     []string{"m"},
				Value:       "cpu-usage",
				DefaultText: "cpu-usage",
				Usage:       "Specifies the check method that should be used.",
			},
			&cli.IntFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Value:       90,
				DefaultText: "90",
				Usage:       "Specifies the timeout for requests.",
			},
			&cli.Float64Flag{
				Name:    "warning",
				Aliases: []string{"w"},
				Usage:   "Specifies the warning threshold.",
			},
			&cli.Float64Flag{
				Name:    "critical",
				Aliases: []string{"c"},
				Usage:   "Specifies the critical threshold.",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Prints additional information for debugging.",
			},
		},
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"V"},
		Usage: "Prints the plugin version.",
	}

	app.Run(os.Args)
}
