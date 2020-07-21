package lib

import "time"

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

// DataModel defines the structure of the JSON response
type DataModel struct {
	Hostname  string        `json:"hostname"`
	Uptime    time.Duration `json:"uptime"`
	CPU       CPUUsage      `json:"cpu"`
	Memory    MemUsage      `json:"memory"`
	Network   []NetUsage    `json:"network"`
	Wireguard []WGPeer      `json:"wireguard"`
}

// ErrorRes defines the structure of the JSON response
type ErrorModel struct {
	Error string `json:"error"`
}
