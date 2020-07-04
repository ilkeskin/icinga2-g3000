package main

import (
	"fmt"
	"log"
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
		// Kbit/s = Bytes * (8 / 1000) --> 1 Mbit/s = 1048.58 Kbit/s
		rxMbps := float64(after[i].RxBytes-before[i].RxBytes) / (125 * 1048.58)
		txMbps := float64(after[i].TxBytes-before[i].TxBytes) / (125 * 1048.58)

		result = append(result, NetUsage{before[i].Name, rxMbps, txMbps})
	}

	return result
}

func main() {
	fmt.Println(getCPUUsage())
	fmt.Println(getMemUsage())
	fmt.Println(getNetUsage())

}
