package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	s "strconv"
	"strings"

	"github.com/fatih/structs"
)

// CPUUsage holds CPU usage
type CPUUsage struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

// MemUsage holds memory usage
type MemUsage struct {
	Used   float64 `json:"used"`
	Cached float64 `json:"cached"`
	Free   float64 `json:"free"`
	//SwapUsed  float64 `json:"swap-used"`
	//SwapFree  float64 `json:"swap-free"`
}

// NetUsage holds network usage
type NetUsage struct {
	Name string  `json:"device"`
	Rx   float64 `json:"rx"`
	Tx   float64 `json:"tx"`
}

// PeerRate holds Wireguard peers date rates
type PeerRate struct {
	Rx float64 `json:"rx"`
	Tx float64 `json:"tx"`
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

func getJSON(host string, port int, path string) JSONSkeleton {
	resp, err := http.Get("http://" + host + ":" + s.Itoa(port) + path)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var skel JSONSkeleton
	json.Unmarshal(body, &skel)

	return skel
}

func getCPUUsage(data JSONSkeleton) string {
	m := structs.Map(data.CPU)
	var result string
	for k, v := range m {
		result += fmt.Sprintf("'%s'=%s%% ", strings.ToLower(k), s.FormatFloat(v.(float64), 'f', 2, 64))
	}
	return strings.TrimSpace(result)
}

func getMemUsage(data JSONSkeleton) string {
	m := structs.Map(data.Memory)
	var result string
	for k, v := range m {
		result += fmt.Sprintf("'%s'=%s%% ", strings.ToLower(k), s.FormatFloat(v.(float64), 'f', 2, 64))
	}
	return strings.TrimSpace(result)
}

func getNICDownstream(data JSONSkeleton, nicname string) string {

	var m map[string]interface{}
	for i := range data.Network {
		if data.Network[i].Name == nicname {
			m = structs.Map(data.Network[i])
		}
	}

	return fmt.Sprintf("'%s'=%skbps ", strings.ToLower(m["Name"].(string)), s.FormatFloat(m["Rx"].(float64), 'f', 2, 64))
}

func getNICUpstream(data JSONSkeleton, nicname string) string {

	var m map[string]interface{}
	for i := range data.Network {
		if data.Network[i].Name == nicname {
			m = structs.Map(data.Network[i])
		}
	}

	return fmt.Sprintf("'%s'=%skbps ", strings.ToLower(m["Name"].(string)), s.FormatFloat(m["Tx"].(float64), 'f', 2, 64))
}

func main() {

	data := getJSON("127.0.0.1", 8888, "/test.json")

	fmt.Println(getCPUUsage(data))
	fmt.Println(getMemUsage(data))
	fmt.Println(getNICDownstream(data, "eth1"))
	fmt.Println(getNICUpstream(data, "eth1"))

}
