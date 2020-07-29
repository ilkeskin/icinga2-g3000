package main

import (
	"fmt"
	"time"

	"github.com/ilkeskin/icinga-g3000/lib"
)

// CheckUptime checks device uptime
func CheckUptime(args CLIArguments) {
	res, err := lib.QueryData(*args.Hostname, *args.Port, "/uptime")
	uptime := res.(lib.Uptime)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	fmt.Printf("OK - %s\n", uptime)
	GlobalReturnCode = exitOk
}

// CheckCPU checks current CPU usage
func CheckCPU(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/cpu")
	cpu := res.(lib.CPUUsage)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	var totalUsed = cpu.System + cpu.User
	output, err := lib.ParseCPUUsage(cpu)

	if args.Warning != nil && totalUsed > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && totalUsed > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get CPU usage\n")
	}
}

// CheckMemory checks current memory (RAM) usage
func CheckMemory(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/memory")
	mem := res.(lib.MemUsage)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	var totalUsed = mem.Cached + mem.Used
	output, err := lib.ParseMemUsage(mem)

	if args.Warning != nil && totalUsed > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && totalUsed > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get memory usage\n")
	}
}

// CheckUpstream checks current upstream of a selected network device
func CheckUpstream(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/network")
	net := res.([]lib.NetUsage)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	var upstream = -1.0
	for i := range net {
		if net[i].Name == *args.NetDevice {
			upstream = net[i].Tx
		}
	}

	output, err := lib.ParseNetUsage(net, *args.NetDevice)

	if args.Warning != nil && upstream > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && upstream > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output[0] + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output[0] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output[0] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get memory usage\n")
	}
}

// CheckDownstream checks current donwstream of a selected network device
func CheckDownstream(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/network")
	netArr := res.([]lib.NetUsage)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	var downstream = -1.0
	for i := range netArr {
		if netArr[i].Name == *args.NetDevice {
			downstream = netArr[i].Rx
		}
	}

	output, err := lib.ParseNetUsage(netArr, *args.NetDevice)

	if args.Warning != nil && downstream > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && downstream > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output[0] + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output[0] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output[0] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get network usage\n")
	}
}

// CheckPeerHandshake checks secs since last handshake for a given WireGuard peer
func CheckPeerHandshake(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/wireguard")
	peerArr := res.([]lib.WGPeer)
	peer, err := lib.GetPeerByIndex(peerArr, *args.Peer)
	output, err := lib.ParsePeer(peerArr, *args.Peer)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	secSinceHS := float64(time.Now().Unix() - peer.LastHS)

	if args.Warning != nil && secSinceHS > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && secSinceHS > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output[0] + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output[0] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output[0] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get secs since laast handshake\n")
	}
}

// CheckPeerUpstream checks current upstream for a given WireGuard peer
func CheckPeerUpstream(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/wireguard")
	peerArr := res.([]lib.WGPeer)
	peer, err := lib.GetPeerByIndex(peerArr, *args.Peer)
	output, err := lib.ParsePeer(peerArr, *args.Peer)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	if args.Warning != nil && peer.PeerRate.Tx > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && peer.PeerRate.Tx > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output[1] + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output[1] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output[1] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get peer upstream\n")
	}
}

// CheckPeerDownstream checks current downstream for a given WireGuard peer
func CheckPeerDownstream(args CLIArguments) {

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/wireguard")
	peerArr := res.([]lib.WGPeer)
	peer, err := lib.GetPeerByIndex(peerArr, *args.Peer)
	output, err := lib.ParsePeer(peerArr, *args.Peer)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	if args.Warning != nil && peer.PeerRate.Rx > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && peer.PeerRate.Rx > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK" + output[2] + "\n")
	case exitWarning:
		fmt.Print("WARNING" + output[2] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL" + output[2] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get peer downstream\n")
	}
}
