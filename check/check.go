package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// CPUUsage holds CPU usage
type CPUUsage struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

// MemUsage holds memory usage
type MemUsage struct {
	RAMUsed   float64 `json:"used"`
	RAMCached float64 `json:"cached"`
	RAMFree   float64 `json:"free"`
	//SwapUsed  float64 `json:"swap-used"`
	//SwapFree  float64 `json:"swap-free"`
}

// NetUsage holds network usage
type NetUsage struct {
	Name   string  `json:"device"`
	RxKbps float64 `json:"rx"`
	TxKbps float64 `json:"tx"`
}

// PeerRate holds Wireguard peers date rates
type PeerRate struct {
	RxKbps float64 `json:"rx"`
	TxKbps float64 `json:"tx"`
}

// WGPeer holds wireguard peer information
type WGPeer struct {
	IntIPAddr string   `json:"internal-ip"`
	ExtIPAddr string   `json:"external-ip"`
	LastHS    int64    `json:"latest-handshake"`
	PeerRate  PeerRate `json:"data-rates"`
}

// JSONSkeleton defines the structure of the JSON response
type JSONSkeleton struct {
	Hostname  string     `json:"hostname"`
	CPU       CPUUsage   `json:"cpu"`
	Memory    MemUsage   `json:"memory"`
	Network   []NetUsage `json:"network"`
	Wireguard []WGPeer   `json:"wireguard"`
}

func getJSON(host string, port int) JSONSkeleton {
	resp, err := http.Get("http://" + host + ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var skel JSONSkeleton
	json.Unmarshal(body, &skel)

	return skel
}

func getValueByKey(key string) string {

	return ""
}

func main() {

	fmt.Println(getJSON("192.168.25.9", 5665))

}
