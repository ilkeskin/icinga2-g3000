package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	s "strings"
	"sync"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/uptime"

	"github.com/ilkeskin/icinga-g3000/lib"
)

// getUptime reads uptime in secs from /proc/uptime and returns it as a time.Duration object.
// If an error occurs while reading those values from the os, an empty object is returned.
func getUptime() (time.Duration, error) {
	var result time.Duration

	result, err := uptime.Get()
	if err != nil {
		return result, err
	}
	return result, nil
}

// getCPUUsage reads the change in the "user"-, "system"- and "idle"-values read from /proc/stats over 1 sec.
// The values are returned as a percentage of the total CPU usage.
// If an error occurs while reading those values from the os, an empty object is returned.
func getCPUUsage() (lib.CPUUsage, error) {
	var result lib.CPUUsage

	before, err := cpu.Get()
	if err != nil {
		return result, err
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		return result, err
	}
	total := float64(after.Total - before.Total)
	user := float64(after.User-before.User) / total * 100
	sys := float64(after.System-before.System) / total * 100
	idle := float64(after.Idle-before.Idle) / total * 100

	result = lib.CPUUsage{User: user, System: sys, Idle: idle}

	return result, nil
}

// getMemUsage reads current memory consumption (used, cached, free, swap) of the os from /proc/meminfo.
// The values are returned as a percentage of the total available memory.
// If an error occurs while reading those values from the os, an empty object is returned.
func getMemUsage() (lib.MemUsage, error) {
	var result lib.MemUsage

	mem, err := memory.Get()
	if err != nil {
		return result, err
	}

	total := float64(mem.Total)
	used := float64(mem.Used) / total * 100
	cached := float64(mem.Cached) / total * 100
	free := float64(mem.Free) / total * 100
	//swapTotal := float64(mem.SwapTotal)
	//swapUsed := float64(mem.SwapUsed) / swapTotal * 100
	//swapFree := float64(mem.SwapFree) / swapTotal * 100
	result = lib.MemUsage{Used: used, Cached: cached, Free: free}

	return result, nil
}

// getNetUsage determines the current RX- and TX-data rates of all availble NICs by sampling received
// and transmitted Bytes over the timespan of 1 sec. Data rates are return as Kbit per second.
// If an error occurs while reading those values from the os, an empty array of objects is returned.
func getNetUsage() ([]lib.NetUsage, error) {
	var result []lib.NetUsage

	before, err := network.Get()
	if err != nil {
		return result, err
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := network.Get()
	if err != nil {
		return result, err
	}

	for i := 0; i < len(before); i++ {
		// Kbit/s = Bytes * (8 / 1000)
		rxKbps := float64(after[i].RxBytes-before[i].RxBytes) / 125
		txKbps := float64(after[i].TxBytes-before[i].TxBytes) / 125

		result = append(result, lib.NetUsage{Name: before[i].Name, Rx: rxKbps, Tx: txKbps})
	}

	return result, nil
}

var result []lib.NetUsage

// parseWGDump parses Wireguard peer information produced by the "wg show wg0 dump" command.
// Interface information is skipped (first line). In case of an error while command exuction
// an empty array of string-arrays is returned.
func parseWGDump() ([][]string, error) {

	var result [][]string

	out, err := exec.Command("wg", "show", "wg0", "dump").Output()
	//out, err := exec.Command("cat", "wg-mock.txt").Output()
	if err != nil {
		return result, fmt.Errorf("Executing \"wg show wg0 dump\" failed: %w", err)
	} else if len(out) == 0 {
		return result, errors.New("Executing \"wg show wg0 dump\" returned empty response")
	}

	lines := s.Split(s.TrimSpace(string(out)), "\n")[1:]
	for i := 0; i < len(lines); i++ {
		result = append(result, s.Split(lines[i], "\t"))
	}

	return result, nil
}

// calcPeersRates calculates RX- and TX-date rates for every configured Wireguard peer,
// by sampling the change in received and transmitted Bytes over the timespan of 1 sec.
func calcPeersRates() ([]lib.PeerRate, error) {
	var result []lib.PeerRate

	before, err := parseWGDump()
	if err != nil {
		return result, err
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := parseWGDump()
	if err != nil {
		return result, err
	}

	for i := 0; i < len(before); i++ {
		rxBefore, err := strconv.ParseFloat(before[i][5], 64)
		rxAfter, err := strconv.ParseFloat(after[i][5], 64)
		txBefore, err := strconv.ParseFloat(before[i][6], 64)
		txAfter, err := strconv.ParseFloat(after[i][6], 64)
		if err != nil {
			return result, err
		}
		// Kbit/s = Bytes * (8 / 1000)
		result = append(result, lib.PeerRate{Rx: (rxAfter - rxBefore) / 125, Tx: (txAfter - txBefore) / 125})
	}
	return result, nil
}

// getWGPeers return all configured Wireguard peers as an array of Peer objects, each including
// its internal and external IP address, the epoch timestamp of its last successful handshake with
// the gateway as well as its data rates.
func getWGPeers(peers [][]string, rates []lib.PeerRate) ([]lib.WGPeer, error) {
	var result []lib.WGPeer

	if len(peers) == 0 {
		return result, errors.New("Peer dump is empty")
	} else if len(rates) == 0 {
		return result, errors.New("Peer dump is empty")
	}

	for i := 0; i < len(peers); i++ {
		lastHS, err := strconv.ParseInt(peers[i][4], 10, 64)
		if err != nil {
			return result, fmt.Errorf("Parsing \"latest handshake\"-value for peer %s failed: %w", peers[i][3], err)
		}

		// IntIPAddr ExtIPAddr LastHS PeerRate
		result = append(result, lib.WGPeer{IntIPAddr: peers[i][3], ExtIPAddr: peers[i][2], LastHS: lastHS, PeerRate: rates[i]})
	}

	return result, nil
}

func sendError(err error) {
	var resp string

	resp = fmt.Sprintf("{\"Error\": \"%s\"}", err)
	fmt.Println("HTTP/1.1 500 Internal Server Error")
	fmt.Println("Content-Type: application/json; charset=utf-8")
	fmt.Println("Content-Length: " + strconv.Itoa(len(resp)))
	fmt.Println("")
	fmt.Print(resp)
}

func sendResult(skel lib.JSONSkeleton) {
	json, err := json.Marshal(skel)
	if err != nil {
		sendError(err)
	}

	fmt.Println("HTTP/1.1 200 OK")
	fmt.Println("Content-Type: application/json; charset=utf-8")
	fmt.Println("Content-Length: " + strconv.Itoa(len(json)))
	fmt.Println("")
	fmt.Print(string(json))
}

func main() {

	var wg sync.WaitGroup
	var uptime time.Duration
	var cpuUsage lib.CPUUsage
	var memUsage lib.MemUsage
	var netUsage []lib.NetUsage
	var peerRates []lib.PeerRate
	var err error

	wg.Add(5)
	go func() {
		defer wg.Done()
		uptime, err = getUptime()
		if err != nil {
		}
	}()

	go func() {
		defer wg.Done()
		cpuUsage, err = getCPUUsage()
		if err != nil {
		}
	}()

	go func() {
		defer wg.Done()
		memUsage, err = getMemUsage()
		if err != nil {
		}
	}()

	go func() {
		defer wg.Done()
		netUsage, err = getNetUsage()
		if err != nil {
		}
	}()

	go func() {
		defer wg.Done()
		peerRates, err = calcPeersRates()
		if err != nil {
		}
	}()
	wg.Wait()

	hostname, err := os.Hostname()
	if err != nil {
	}

	wgDump, err := parseWGDump()
	if err != nil {
	}

	peers, err := getWGPeers(wgDump, peerRates)
	if err != nil {
	}

	skel := lib.JSONSkeleton{
		Hostname:  hostname,
		Uptime:    uptime,
		CPU:       cpuUsage,
		Memory:    memUsage,
		Network:   netUsage,
		Wireguard: peers,
	}

	sendResult(skel)
}
