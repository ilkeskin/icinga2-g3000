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

// CLIArguments holds arguments passed from the cli
type CLIArguments struct {
	Hostname  *string
	Port      *int
	Warning   *float64
	Critical  *float64
	Timeout   *int
	NetDevice *string
	Peer      *int64
	Verbose   bool
}

func setRequired(hostname string, port int, timeout int) CLIArguments {
	var args CLIArguments

	args.Hostname = &hostname
	args.Port = &port
	args.Timeout = &timeout
	args.Verbose = false

	return args
}

func (args *CLIArguments) setWarning(warning float64) {
	args.Warning = &warning
}

func (args *CLIArguments) setCritical(critical float64) {
	args.Critical = &critical
}

func (args *CLIArguments) setNetDevice(netdevice string) {
	args.NetDevice = &netdevice
}

func (args *CLIArguments) setPeer(peer int64) {
	args.Peer = &peer
}

func (args *CLIArguments) setVerbose() {
	args.Verbose = true
}

func checkRequiredFlags(args *CLIArguments) bool {
	if args.Hostname == nil || *args.Hostname == "" {
		fmt.Println("No hostname or IP address was set")
		return false
	}

	if args.Port == nil || *args.Port < 1023 || *args.Port > 65535 {
		fmt.Println("No port was set or port number out of range")
		return false
	}

	if args.Timeout == nil || *args.Timeout == 0 || *args.Timeout > 120 {
		fmt.Println("No timeout was set or timeout is 0 or greater than 2 mins")
		return false
	}

	return true
}

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
				Aliases:     []string{"up", "u"},
				Usage:       "get device uptime (in s)",
				Description: "retrieves the device uptime since last (re)boot",
				Action: func(c *cli.Context) error {
					cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

					if c.IsSet("warning") {
						cliArgs.setWarning(c.Float64("warning"))
					}

					if c.IsSet("critical") {
						cliArgs.setCritical(c.Float64("critical"))
					}

					if !checkRequiredFlags(&cliArgs) {
						os.Exit(exitUnknown)
					}

					CheckUptime(cliArgs)

					os.Exit(GlobalReturnCode)
					return nil
				},
			},
			&cli.Command{
				Name:        "cpu",
				Aliases:     []string{"c"},
				Usage:       "get CPU usage (in %)",
				Description: "retrieves the current CPU usage as a percentage of total cpu time split between user, system and idle",
				Action: func(c *cli.Context) error {
					cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

					if c.IsSet("warning") {
						cliArgs.setWarning(c.Float64("warning"))
					}

					if c.IsSet("critical") {
						cliArgs.setCritical(c.Float64("critical"))
					}

					if !checkRequiredFlags(&cliArgs) {
						os.Exit(exitUnknown)
					}

					CheckCPU(cliArgs)

					os.Exit(GlobalReturnCode)
					return nil
				},
			},
			&cli.Command{
				Name:        "memory",
				Aliases:     []string{"mem", "m"},
				Usage:       "get memory usage (in %)",
				Description: "retrieves the current Memory usage as a percentage of total memory available split between used, cached and free RAM",
				Action: func(c *cli.Context) error {
					cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

					if c.IsSet("warning") {
						cliArgs.setWarning(c.Float64("warning"))
					}

					if c.IsSet("critical") {
						cliArgs.setCritical(c.Float64("critical"))
					}

					if !checkRequiredFlags(&cliArgs) {
						os.Exit(exitUnknown)
					}

					CheckMemory(cliArgs)

					os.Exit(GlobalReturnCode)
					return nil
				},
			},
			&cli.Command{
				Name:        "network",
				Aliases:     []string{"net", "n"},
				Usage:       "get network usage (in kbps) of a NIC",
				Description: "retrieves the current network usage split into kbps up- and downstream",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "device",
						Aliases:     []string{"d"},
						Value:       "eth0",
						DefaultText: "eth0",
						Usage:       "Specifies the device that should be queried",
					},
				},
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:        "upstream",
						Aliases:     []string{"up", "u"},
						Usage:       "get NIC upstream (in kbps)",
						Description: "retrieves current upstream in kbps for a given network device",
						Action: func(c *cli.Context) error {
							cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

							if c.IsSet("device") {
								cliArgs.setNetDevice(c.String("device"))
							} else {
								cli.ShowCommandHelp(c, "network")
								os.Exit(exitUnknown)
							}

							if c.IsSet("warning") {
								cliArgs.setWarning(c.Float64("warning"))
							}

							if c.IsSet("critical") {
								cliArgs.setCritical(c.Float64("critical"))
							}

							if !checkRequiredFlags(&cliArgs) {
								os.Exit(exitUnknown)
							}

							CheckUpstream(cliArgs)

							os.Exit(GlobalReturnCode)
							return nil
						},
					},
					&cli.Command{
						Name:        "downstream",
						Aliases:     []string{"down", "d"},
						Usage:       "get NIC donwstream (in kbps)",
						Description: "retrieves current downstream in kbps for a given network device",
						Action: func(c *cli.Context) error {
							cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

							if c.IsSet("device") {
								cliArgs.setNetDevice(c.String("device"))
							} else {
								cli.ShowCommandHelp(c, "network")
								os.Exit(exitUnknown)
							}

							if c.IsSet("warning") {
								cliArgs.setWarning(c.Float64("warning"))
							}

							if c.IsSet("critical") {
								cliArgs.setCritical(c.Float64("critical"))
							}

							if !checkRequiredFlags(&cliArgs) {
								os.Exit(exitUnknown)
							}

							CheckDownstream(cliArgs)

							os.Exit(GlobalReturnCode)
							return nil
						},
					},
				},
			},
			&cli.Command{
				Name:        "wireguard",
				Aliases:     []string{"wg", "w"},
				Usage:       "get WireGuard related information",
				Description: "retrieves up- and downstream speeds and time since the last handshake for a specified WireGuard peer",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "peer",
						Aliases:     []string{"P"},
						Value:       1,
						DefaultText: "1",
						Usage:       "Specifies the WireGuard peer which should be queried. Peers are identified by the last octet of their IP address",
					},
				},
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:        "handshake",
						Aliases:     []string{"hs"},
						Usage:       "get seconds since last handshake (in s)",
						Description: "retrieves seconds since last handshake with gateway",
						Action: func(c *cli.Context) error {
							cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

							if c.IsSet("peer") {
								cliArgs.setPeer(c.Int64("peer"))
							} else {
								os.Exit(exitUnknown)
							}

							if c.IsSet("warning") {
								cliArgs.setWarning(c.Float64("warning"))
							}

							if c.IsSet("critical") {
								cliArgs.setCritical(c.Float64("critical"))
							}

							if !checkRequiredFlags(&cliArgs) {
								os.Exit(exitUnknown)
							}

							CheckPeerHandshake(cliArgs)

							os.Exit(GlobalReturnCode)
							return nil
						},
					},
					&cli.Command{
						Name:        "downstream",
						Aliases:     []string{"down", "d"},
						Usage:       "get current downstream (in kbps)",
						Description: "retrieves the current downstream of the selected WireGuard peer in kbit per second",
						Action: func(c *cli.Context) error {
							cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

							if c.IsSet("peer") {
								cliArgs.setPeer(c.Int64("peer"))
							} else {
								os.Exit(exitUnknown)
							}

							if c.IsSet("warning") {
								cliArgs.setWarning(c.Float64("warning"))
							}

							if c.IsSet("critical") {
								cliArgs.setCritical(c.Float64("critical"))
							}

							if !checkRequiredFlags(&cliArgs) {
								os.Exit(exitUnknown)
							}

							CheckPeerDownstream(cliArgs)

							os.Exit(GlobalReturnCode)
							return nil
						},
					},
					&cli.Command{
						Name:        "upstream",
						Aliases:     []string{"up", "u"},
						Usage:       "get current upstream (in kbps)",
						Description: "retrieves the current uprstream of the selected WireGuard peer in kbit per second",
						Action: func(c *cli.Context) error {
							cliArgs := setRequired(c.String("hostname"), c.Int("port"), c.Int("timeout"))

							if c.IsSet("peer") {
								cliArgs.setPeer(c.Int64("peer"))
							} else {
								os.Exit(exitUnknown)
							}

							if c.IsSet("warning") {
								cliArgs.setWarning(c.Float64("warning"))
							}

							if c.IsSet("critical") {
								cliArgs.setCritical(c.Float64("critical"))
							}

							if !checkRequiredFlags(&cliArgs) {
								os.Exit(exitUnknown)
							}

							CheckPeerUpstream(cliArgs)

							os.Exit(GlobalReturnCode)
							return nil
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "hostname",
				Aliases:     []string{"H"},
				Value:       "192.168.25.10",
				DefaultText: "192.168.25.10",
				Usage:       "Specifies the hostname or IP-address of the G3000 gateway",
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       5665,
				DefaultText: "5665",
				Usage:       "Specifies the port to query",
			},
			&cli.IntFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Value:       90,
				DefaultText: "90",
				Usage:       "Specifies the timeout for requests",
			},
			&cli.Float64Flag{
				Name:    "warning",
				Aliases: []string{"w"},
				Usage:   "Specifies the warning threshold",
			},
			&cli.Float64Flag{
				Name:    "critical",
				Aliases: []string{"c"},
				Usage:   "Specifies the critical threshold",
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				DefaultText: "false",
				Usage:       "Prints additional information for debugging",
			},
		},
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"V"},
		Usage: "Prints the plugin version",
	}

	app.Run(os.Args)
}
