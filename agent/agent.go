package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	s "strings"
	"sync"
	"time"

	"github.com/ilkeskin/icinga-g3000/lib"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
)

// getCPUUsage reads the change in the "user"-, "system"- and "idle"-values read from /proc/stats over 1 sec.
// The values are returned as a percentage of the total CPU usage.
// If an error occurs while reading those values from the os, an empty object is returned.
func getCPUUsage() lib.CPUUsage {
	before, err := cpu.Get()
	if err != nil {
		log.Fatal("Could not read CPU stats: " + err.Error())
		return lib.CPUUsage{}
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		log.Fatal("Could not read CPU stats: " + err.Error())
		return lib.CPUUsage{}
	}
	total := float64(after.Total - before.Total)
	user := float64(after.User-before.User) / total * 100
	sys := float64(after.System-before.System) / total * 100
	idle := float64(after.Idle-before.Idle) / total * 100

	return lib.CPUUsage{user, sys, idle}
}

// getMemUsage reads current memory consumption (used, cached, free, swap) of the os from /proc/meminfo.
// The values are returned as a percentage of the total available memory.
// If an error occurs while reading those values from the os, an empty object is returned.
func getMemUsage() lib.MemUsage {
	mem, err := memory.Get()
	if err != nil {
		log.Fatal("Could not read memory stats: " + err.Error())
		return lib.MemUsage{}
	}

	total := float64(mem.Total)
	used := float64(mem.Used) / total * 100
	cached := float64(mem.Cached) / total * 100
	free := float64(mem.Free) / total * 100
	//swapTotal := float64(mem.SwapTotal)
	//swapUsed := float64(mem.SwapUsed) / swapTotal * 100
	//swapFree := float64(mem.SwapFree) / swapTotal * 100

	return lib.MemUsage{used, cached, free} //, swapUsed, swapFree}
}

// getNetUsage determines the current RX- and TX-data rates of all availble NICs by sampling received
// and transmitted Bytes over the timespan of 1 sec. Data rates are return as Kbit per second.
// If an error occurs while reading those values from the os, an empty array of objects is returned.
func getNetUsage() []lib.NetUsage {
	before, err := network.Get()
	if err != nil {
		log.Fatal("Could not read network stats: " + err.Error())
		return []lib.NetUsage{}
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := network.Get()
	if err != nil {
		log.Fatal("Could not read network stats: " + err.Error())
		return []lib.NetUsage{}
	}

	var result []lib.NetUsage

	for i := 0; i < len(before); i++ {
		// Kbit/s = Bytes * (8 / 1000)
		rxKbps := float64(after[i].RxBytes-before[i].RxBytes) / 125
		txKbps := float64(after[i].TxBytes-before[i].TxBytes) / 125

		result = append(result, lib.NetUsage{before[i].Name, rxKbps, txKbps})
	}

	return result
}

// parseWGDump parses Wireguard peer information produced by the "wg show wg0 dump" command.
// Interface information is skipped (first line). In case of an error while command exuction
// an empty array of string-arrays is returned.
func parseWGDump() [][]string {

	var peers [][]string

	out, err := exec.Command("wg", "show", "wg0", "dump").Output()
	//out, err := exec.Command("cat", "wg-mock.txt").Output()
	if err != nil {
		log.Fatal("Could not read Wireguard config: " + err.Error())
		return peers
	}

	lines := s.Split(s.TrimSpace(string(out)), "\n")[1:]

	for i := 0; i < len(lines); i++ {
		peers = append(peers, s.Split(lines[i], "\t"))
	}

	return peers
}

// calcPeersRates calculates RX- and TX-date rates for every configured Wireguard peer,
// by sampling the change in received and transmitted Bytes over the timespan of 1 sec.
func calcPeersRates() []lib.PeerRate {

	before := parseWGDump()
	time.Sleep(time.Duration(1) * time.Second)
	after := parseWGDump()

	var result []lib.PeerRate
	for i := 0; i < len(before); i++ {
		rxBefore, err := strconv.ParseFloat(before[i][5], 64)
		rxAfter, err := strconv.ParseFloat(after[i][5], 64)
		txBefore, err := strconv.ParseFloat(before[i][6], 64)
		txAfter, err := strconv.ParseFloat(after[i][6], 64)
		if err != nil {
			log.Fatal("Could not parse RX/TX values for peers: " + err.Error())
		}
		// Kbit/s = Bytes * (8 / 1000)
		result = append(result, lib.PeerRate{(rxAfter - rxBefore) / 125, (txAfter - txBefore) / 125})
	}
	return result
}

// getWGPeers return all configured Wireguard peers as an array of Peer objects, each including
// its internal and external IP address, the epoch timestamp of its last succesful handshake with
// the gateway as well as its data rates.
func getWGPeers(peers [][]string, rates []lib.PeerRate) []lib.WGPeer {
	var result []lib.WGPeer
	for i := 0; i < len(peers); i++ {
		lastHS, err := strconv.ParseInt(peers[i][4], 10, 64)
		if err != nil {
			log.Fatal("Could not parse last-handshake value for peer " + strconv.Itoa(i) + ": " + err.Error())
		}

		// IntIPAddr ExtIPAddr LastHS PeerRate
		result = append(result, lib.WGPeer{peers[i][3], peers[i][2], lastHS, rates[i]})
	}

	return result
}

func main() {
	var wg sync.WaitGroup
	var cpuUsage lib.CPUUsage
	var memUsage lib.MemUsage
	var netUsage []lib.NetUsage
	var peerRates []lib.PeerRate
	wg.Add(4)
	go func() {
		defer wg.Done()
		cpuUsage = getCPUUsage()
	}()

	go func() {
		defer wg.Done()
		memUsage = getMemUsage()
	}()

	go func() {
		defer wg.Done()
		netUsage = getNetUsage()
	}()

	go func() {
		defer wg.Done()
		peerRates = calcPeersRates()
	}()
	wg.Wait()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Could not get hostname: " + err.Error())
	}

	skel := lib.JSONSkeleton{hostname, cpuUsage, memUsage, netUsage, getWGPeers(parseWGDump(), peerRates)}

	jsonRes, err := json.Marshal(skel)

	fmt.Println("HTTP/1.1 200 OK")
	fmt.Println("Content-Type: application/json; charset=utf-8")
	fmt.Println("Content-Length: " + strconv.Itoa(len(jsonRes)))
	fmt.Println("")
	fmt.Print(string(jsonRes))

}
