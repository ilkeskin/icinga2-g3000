package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	s "strings"
	"sync"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
)

// CPUUsage holds CPU usage
type CPUUsage struct {
	User   float64
	System float64
	Idle   float64
}

// MemUsage holds memory usage
type MemUsage struct {
	RAMUsed   float64
	RAMCached float64
	RAMFree   float64
	SwapUsed  float64
	SwapFree  float64
}

// NetUsage holds network usage
type NetUsage struct {
	Name   string
	RxMbps float64
	TxMbps float64
}

// PeerRate holds Wireguard peers date rates
type PeerRate struct {
	RxMbps float64
	TxMbps float64
}

// WGPeer holds wireguard peer information
type WGPeer struct {
	IntIPAddr string
	ExtIPAddr string
	LastHS    int64
	PeerRate  PeerRate
}

func getCPUUsage() CPUUsage {
	before, err := cpu.Get()
	if err != nil {
		log.Fatal("Could not read CPU stats: " + err.Error())
		return CPUUsage{0.0, 0.0, 0.0}
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		log.Fatal("Could not read CPU stats: " + err.Error())
		return CPUUsage{0.0, 0.0, 0.0}
	}
	total := float64(after.Total - before.Total)
	user := float64(after.User-before.User) / total * 100
	sys := float64(after.System-before.System) / total * 100
	idle := float64(after.Idle-before.Idle) / total * 100

	return CPUUsage{user, sys, idle}
}

func getMemUsage() MemUsage {
	mem, err := memory.Get()
	if err != nil {
		log.Fatal("Could not read memory stats: " + err.Error())
		return MemUsage{0, 0, 0, 0, 0}
	}

	total := float64(mem.Total)
	used := float64(mem.Used) / total * 100
	cached := float64(mem.Cached) / total * 100
	free := float64(mem.Free) / total * 100
	swapTotal := float64(mem.SwapTotal)
	swapUsed := float64(mem.SwapUsed) / swapTotal * 100
	swapFree := float64(mem.SwapFree) / swapTotal * 100

	return MemUsage{used, cached, free, swapUsed, swapFree}
}

func getNetUsage() []NetUsage {
	before, err := network.Get()
	if err != nil {
		log.Fatal("Could not read network stats: " + err.Error())
		return []NetUsage{}
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := network.Get()
	if err != nil {
		log.Fatal("Could not read network stats: " + err.Error())
		return []NetUsage{}
	}

	var result []NetUsage

	for i := 0; i < len(before); i++ {
		// Kbit/s = Bytes * (8 / 1000) --> 1 Mbit/s = 1000 Kbit/s
		rxMbps := float64(after[i].RxBytes-before[i].RxBytes) / (125 * 1000)
		txMbps := float64(after[i].TxBytes-before[i].TxBytes) / (125 * 1000)

		result = append(result, NetUsage{before[i].Name, rxMbps, txMbps})
	}

	return result
}

func parseWGDump() [][]string {

	//out, err := exec.Command("wg", "show", "wg0", "dump").Output()
	out, err := exec.Command("cat", "wg-mock.txt").Output()
	if err != nil {
		log.Fatal("Could not read Wireguard config: " + err.Error())
	}

	lines := s.Split(s.TrimSpace(string(out)), "\n")[1:]

	var peers [][]string
	for i := 0; i < len(lines); i++ {
		peers = append(peers, s.Split(lines[i], "\t"))
	}

	return peers
}

func calcPeersRates() []PeerRate {

	before := parseWGDump()
	time.Sleep(time.Duration(1) * time.Second)
	after := parseWGDump()

	var result []PeerRate
	for i := 0; i < len(before); i++ {
		rxBefore, err := strconv.ParseFloat(before[i][5], 64)
		rxAfter, err := strconv.ParseFloat(after[i][5], 64)
		txBefore, err := strconv.ParseFloat(before[i][6], 64)
		txAfter, err := strconv.ParseFloat(after[i][6], 64)
		if err != nil {
			log.Fatal("Could not parse RX/TX values for peers: " + err.Error())
		}
		// Kbit/s = Bytes * (8 / 1000) --> 1 Mbit/s = 1000 Kbit/s
		result = append(result, PeerRate{(rxAfter - rxBefore) / (125 * 1000), (txAfter - txBefore) / (125 * 1000)})
	}
	return result
}

func getWGPeers(peers [][]string, rates []PeerRate) []WGPeer {
	var result []WGPeer
	for i := 0; i < len(peers); i++ {
		lastHS, err := strconv.ParseInt(peers[i][4], 10, 64)
		if err != nil {
			log.Fatal("Could not parse last-handshake value for peer " + strconv.Itoa(i) + ": " + err.Error())
		}

		// IntIPAddr ExtIPAddr LastHS PeerRate
		result = append(result, WGPeer{peers[i][3], peers[i][2], lastHS, rates[i]})
	}

	return result
}

func main() {
	var wg sync.WaitGroup
	var cpuUsage CPUUsage
	var memUsage MemUsage
	var netUsage []NetUsage
	var peerRates []PeerRate
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

	fmt.Println(cpuUsage)
	fmt.Println(memUsage)
	fmt.Println(netUsage)
	fmt.Println(getWGPeers(parseWGDump(), peerRates))
}
