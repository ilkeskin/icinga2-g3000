package main

import (
	"fmt"
	"time"

	"github.com/ilkeskin/icinga-g3000/lib"
	"github.com/mitchellh/mapstructure"
	ms "github.com/mitchellh/mapstructure"
)

// CheckUptime checks device uptime
func CheckUptime(args CLIArguments) {
	var uptime lib.Uptime
	res, err := lib.QueryData(*args.Hostname, *args.Port, "/uptime", *args.Timeout)
	err = ms.Decode(res, &uptime)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}
	GlobalReturnCode = exitOk

	fmt.Printf("OK - 'uptime'=%ds\n", int(uptime.Uptime.Seconds()))
}

// CheckCPU checks current CPU usage
func CheckCPU(args CLIArguments) {
	var cpu lib.CPUUsage
	res, err := lib.QueryData(*args.Hostname, *args.Port, "/cpu", *args.Timeout)
	err = ms.Decode(res, &cpu)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	var totalUsed = cpu.System + cpu.User
	output, err := lib.ParseCPUUsage(cpu)
	GlobalReturnCode = exitOk

	if args.Warning != nil && totalUsed > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && totalUsed > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK - " + output + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get CPU usage\n")
	}
}

// CheckMemory checks current memory (RAM) usage
func CheckMemory(args CLIArguments) {
	var mem lib.MemUsage
	res, err := lib.QueryData(*args.Hostname, *args.Port, "/memory", *args.Timeout)
	err = ms.Decode(res, &mem)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	var totalUsed = mem.Cached + mem.Used
	output, err := lib.ParseMemUsage(mem)
	GlobalReturnCode = exitOk

	if args.Warning != nil && totalUsed > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && totalUsed > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK - " + output + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get memory usage\n")
	}
}

// CheckUpstream checks current upstream of a selected network device
func CheckUpstream(args CLIArguments) {
	var netArr []lib.NetUsage

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/network", *args.Timeout)

	config := &ms.DecoderConfig{
		TagName: "json",
		Result:  &netArr,
	}
	decoder, err := mapstructure.NewDecoder(config)
	decoder.Decode(res)

	var upstream float64
	for i := range netArr {
		if netArr[i].Name == *args.NetDevice {
			upstream = netArr[i].Rx
		}
	}

	output, err := lib.ParseNetUsage(netArr, *args.NetDevice)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}
	GlobalReturnCode = exitOk

	if args.Warning != nil && upstream > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && upstream > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK - " + output[0] + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output[0] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output[0] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get upstream of device " + *args.NetDevice + "\n")
	}
}

// CheckDownstream checks current donwstream of a selected network device
func CheckDownstream(args CLIArguments) {
	var netArr []lib.NetUsage

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/network", *args.Timeout)

	config := &ms.DecoderConfig{
		TagName: "json",
		Result:  &netArr,
	}
	decoder, err := mapstructure.NewDecoder(config)
	decoder.Decode(res)

	var upstream float64
	for i := range netArr {
		if netArr[i].Name == *args.NetDevice {
			upstream = netArr[i].Rx
		}
	}

	output, err := lib.ParseNetUsage(netArr, *args.NetDevice)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	GlobalReturnCode = exitOk

	if args.Warning != nil && upstream > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && upstream > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK - " + output[1] + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output[1] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output[1] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get upstream of device " + *args.NetDevice + "\n")
	}
}

// CheckPeerHandshake checks secs since last handshake for a given WireGuard peer
func CheckPeerHandshake(args CLIArguments) {
	var peerArr []lib.WGPeer

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/wireguard", *args.Timeout)

	config := &ms.DecoderConfig{
		TagName: "json",
		Result:  &peerArr,
	}
	decoder, err := mapstructure.NewDecoder(config)
	decoder.Decode(res)

	peer, err := lib.GetPeerByIndex(peerArr, *args.Peer)
	output, err := lib.ParsePeer(peerArr, *args.Peer)
	if err != nil {
		fmt.Printf("UNKNOWN - %s\n", err)
		return
	}

	//fmt.Printf("Epoch is %d\nSecs since last handshake are %d\n", time.Now().Unix(), peer.LastHS)
	secSinceHS := float64(time.Now().Unix() - peer.LastHS)
	GlobalReturnCode = exitOk

	if args.Warning != nil && secSinceHS > *args.Warning {
		GlobalReturnCode = exitWarning
	}

	if args.Critical != nil && secSinceHS > *args.Critical {
		GlobalReturnCode = exitCritical
	}

	switch GlobalReturnCode {
	case exitOk:
		fmt.Print("OK - " + output[0] + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output[0] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output[0] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get secs since last handshake\n")
	}
}

// CheckPeerUpstream checks current upstream for a given WireGuard peer
func CheckPeerUpstream(args CLIArguments) {
	var peerArr []lib.WGPeer

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/wireguard", *args.Timeout)

	config := &ms.DecoderConfig{
		TagName: "json",
		Result:  &peerArr,
	}
	decoder, err := mapstructure.NewDecoder(config)
	decoder.Decode(res)

	peer, err := lib.GetPeerByIndex(peerArr, *args.Peer)
	output, err := lib.ParsePeer(peerArr, *args.Peer)
	GlobalReturnCode = exitOk

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
		fmt.Print("OK - " + output[1] + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output[1] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output[1] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get peer upstream\n")
	}
}

// CheckPeerDownstream checks current downstream for a given WireGuard peer
func CheckPeerDownstream(args CLIArguments) {
	var peerArr []lib.WGPeer

	res, err := lib.QueryData(*args.Hostname, *args.Port, "/wireguard", *args.Timeout)

	config := &ms.DecoderConfig{
		TagName: "json",
		Result:  &peerArr,
	}
	decoder, err := mapstructure.NewDecoder(config)
	decoder.Decode(res)

	peer, err := lib.GetPeerByIndex(peerArr, *args.Peer)
	output, err := lib.ParsePeer(peerArr, *args.Peer)
	GlobalReturnCode = exitOk

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
		fmt.Print("OK - " + output[2] + "\n")
	case exitWarning:
		fmt.Print("WARNING - " + output[2] + "\n")
	case exitCritical:
		fmt.Print("CRITICAL - " + output[2] + "\n")
	default:
		GlobalReturnCode = exitUnknown
		fmt.Print("UNKNOWN - Could not get peer downstream\n")
	}
}
