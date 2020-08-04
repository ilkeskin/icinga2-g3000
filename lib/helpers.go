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

/*// QueryData issues a HTTP-GET request on a specified host and port and
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
}*/

// QueryData issues a HTTP-GET request on a specified host and port and
// unmarshals the received JSON body into the shared data structure.
func QueryData(host string, port int, path string, timeout int) (interface{}, error) {
	var result interface{}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get("http://" + host + ":" + s.Itoa(port) + path)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Could not read body from HTTP response")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return result, nil
	} else if resp.StatusCode == 500 {
		return nil, errors.New(result.(ErrorModel).Error)
	}

	return nil, errors.New("Unknown status code in HTTP response")
}

/*// ParseCPUUsage parses the CPU usage related metrics retrieved from the agent
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
}*/

// ParseCPUUsage parses the CPU usage related metrics retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParseCPUUsage(data CPUUsage) (string, error) {
	var result string

	if !structs.IsStruct(data) {
		return result, errors.New("data is not a struct")
	}

	return fmt.Sprintf("'user'=%s%% 'system'=%s%% 'idle'=%s%%",
		s.FormatFloat(data.User, 'f', 2, 64),
		s.FormatFloat(data.System, 'f', 2, 64),
		s.FormatFloat(data.Idle, 'f', 2, 64)), nil
}

/*// ParseMemUsage parses the memory usage related metrics retrieved from the agent
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
}*/

// ParseMemUsage parses the memory usage related metrics retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParseMemUsage(data MemUsage) (string, error) {
	var result string

	if !structs.IsStruct(data) {
		return result, errors.New("data is not a struct")
	}

	return fmt.Sprintf("'used'=%s%% 'cached'=%s%% 'free'=%s%%",
		s.FormatFloat(data.Used, 'f', 2, 64),
		s.FormatFloat(data.Cached, 'f', 2, 64),
		s.FormatFloat(data.Free, 'f', 2, 64)), nil
}

/*// ParseNetUsage parses the network usage related metrics of a specified NIC retrieved from the agent
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
}*/

// ParseNetUsage parses the network usage related metrics of a specified NIC retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParseNetUsage(data []NetUsage, nicname string) ([2]string, error) {
	var result [2]string

	if len(data) == 0 {
		return result, errors.New("Data holds no devices")
	}

	for i := range data {
		if data[i].Name == nicname {
			result[0] = fmt.Sprintf("'upstream'=%skbps", s.FormatFloat(data[i].Tx, 'f', 2, 64))
			result[1] = fmt.Sprintf("'downstream'=%skbps", s.FormatFloat(data[i].Rx, 'f', 2, 64))
			return result, nil
		}
	}

	return result, errors.New("Could not find device with name " + nicname)
}

/*// ParsePeer parses the Wireguard related metrics of a specified peer retrieved from the agent
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
}*/

// ParsePeer parses the Wireguard related metrics of a specified peer retrieved from the agent
// into a format that is understood by Icingas API and returns them as a string.
func ParsePeer(data []WGPeer, index int64) ([3]string, error) {
	var result [3]string

	peer, err := GetPeerByIndex(data, index)
	if err != nil {
		return result, err
	}

	result[0] = fmt.Sprintf("'lasths'=%ds", time.Now().Unix()-peer.LastHS)
	result[1] = fmt.Sprintf("'upstream'=%.2fkbps", peer.PeerRate.Tx)
	result[2] = fmt.Sprintf("'downstream'=%.2fkbps", peer.PeerRate.Rx)

	return result, nil
}

// GetPeerByIndex returns peer with given index based on the last octet of its internal IP address
func GetPeerByIndex(peerArr []WGPeer, index int64) (WGPeer, error) {
	var result WGPeer

	for i := range peerArr {
		ip := strings.Split(peerArr[i].IntIPAddr, "/")
		ip = strings.Split(ip[0], ".")

		idx, err := s.ParseInt(ip[len(ip)-1], 10, 8)
		if err != nil {
			return result, err
		}

		//fmt.Printf("Testing %d against given index %d: %t\n", idx, index, idx == index)

		if idx == index {
			result = peerArr[i]
			// fmt.Println(result)
			return result, nil
		}
	}

	return result, errors.New("Could not find peer with index " + s.FormatInt(index, 10))
}

/*func isJSONArray(data []byte) bool {
	x := bytes.TrimLeft(data, " \t\r\n")
	//isArray := len(x) > 0 && x[0] == '['
	//isObject := len(x) > 0 && x[0] == '{'
	return len(x) > 0 && x[0] == '['
}*/
