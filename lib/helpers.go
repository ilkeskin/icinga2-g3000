package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	s "strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
)

// QueryData issues a HTTP-GET request on a specified host and port and
// unmarshals the received JSON body into the shared data structure.
func QueryData(host string, port int, path string) (DataModel, error) {
	var skel DataModel

	resp, err := http.Get("http://" + host + ":" + s.Itoa(port) + path)
	if err != nil {
		return skel, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &skel)

	return skel, nil
}

// ParseCPUUsage parses the CPU usage related metrics retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParseCPUUsage(data DataModel) (string, error) {
	var result string

	if !structs.IsStruct(data.CPU) {
		return result, errors.New("data.CPU is not a struct, empty response?")
	}
	m := structs.Map(data.CPU)

	for k, v := range m {
		result += fmt.Sprintf("'%s'=%s%% ",
			strings.ToLower(k),
			s.FormatFloat(v.(float64), 'f', 2, 64))
	}
	return strings.TrimSpace(result), nil
}

// ParseMemUsage parses the memory usage related metrics retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParseMemUsage(data DataModel) (string, error) {
	var result string

	if !structs.IsStruct(data.Memory) {
		return result, errors.New("data.Memory is not a struct, empty response?")
	}
	m := structs.Map(data.Memory)

	for k, v := range m {
		result += fmt.Sprintf("'%s'=%s%% ",
			strings.ToLower(k),
			s.FormatFloat(v.(float64), 'f', 2, 64))
	}
	return strings.TrimSpace(result), nil
}

// ParseNetUsage parses the network usage related metrics of a specified NIC retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParseNetUsage(data DataModel, nicname string) (string, error) {
	var result string

	if len(data.Network) == 0 {
		return result, errors.New("data.Network holds no NICs, empty response?")
	}

	var m map[string]interface{}
	for i := range data.Network {
		if data.Network[i].Name == nicname {
			m = structs.Map(data.Network[i])
			result = fmt.Sprintf("'down'=%skbps 'up'=%skbps",
				s.FormatFloat(m["Rx"].(float64), 'f', 2, 64),
				s.FormatFloat(m["Tx"].(float64), 'f', 2, 64))
		}
	}

	return result, nil
}

// ParsePeer parses the Wireguard related metrics of a specified peer retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParsePeer(data DataModel, peerIndex int64) (string, error) {
	var result string

	for i := range data.Wireguard {
		ip := strings.Split(data.Wireguard[i].IntIPAddr, "/")
		ip = strings.Split(ip[0], ".")

		idx, err := s.ParseInt(ip[len(ip)-1], 10, 8)
		if err != nil {
			return result, err
		}

		if idx == peerIndex {
			result = fmt.Sprintf("'lasths'=%ds 'down'=%.2fkbps 'up'=%.2fkbps",
				time.Now().Unix()-data.Wireguard[i].LastHS,
				data.Wireguard[i].PeerRate.Rx,
				data.Wireguard[i].PeerRate.Tx)
			return result, nil
		}
	}
	return result, errors.New("Could not find peer with index " + s.FormatInt(peerIndex, 10) + ", empty response?")
}
