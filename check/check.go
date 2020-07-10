package main

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
	"github.com/ilkeskin/icinga-g3000/lib"
)

func getData(host string, port int, path string) (lib.JSONSkeleton, error) {
	var skel lib.JSONSkeleton

	resp, err := http.Get("http://" + host + ":" + s.Itoa(port) + path)
	if err != nil {
		return skel, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &skel)

	return skel, nil
}

func getCPUUsage(data lib.JSONSkeleton) (string, error) {
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

func getMemUsage(data lib.JSONSkeleton) (string, error) {
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

func getNICUsage(data lib.JSONSkeleton, nicname string) (string, error) {
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

func getPeerSecsSinceHS(data lib.JSONSkeleton, peerIndex int64) (string, error) {
	var result string

	for i := range data.Wireguard {
		ip := strings.Split(data.Wireguard[i].IntIPAddr, "/")
		ip = strings.Split(ip[0], ".")

		idx, err := s.ParseInt(ip[len(ip)-1], 10, 8)
		if err != nil {
			return result, err
		}

		if idx == peerIndex {
			result = fmt.Sprintf("'lasths'=%ds", time.Now().Unix()-data.Wireguard[i].LastHS)
			return result, nil
		}
	}
	return result, errors.New("Could not find peer with index " + s.FormatInt(peerIndex, 10) + ",empty response?")
}

func getPeerUsage(data lib.JSONSkeleton, peerIndex int64) (string, error) {
	var result string

	for i := range data.Wireguard {
		ip := strings.Split(data.Wireguard[i].IntIPAddr, "/")
		ip = strings.Split(ip[0], ".")

		idx, err := s.ParseInt(ip[len(ip)-1], 10, 8)
		if err != nil {
			return result, err
		}

		if idx == peerIndex {
			result = fmt.Sprintf("'down'=%.2fkbps 'up'=%.2fkbps", data.Wireguard[i].PeerRate.Rx, data.Wireguard[i].PeerRate.Tx)
			return result, nil
		}
	}
	return result, errors.New("Could not find peer with index " + s.FormatInt(peerIndex, 10) + ",empty response?")
}

func main() {
	data, err := getData("127.0.0.1", 8888, "/test.json")

	fmt.Println(getCPUUsage(data))
	fmt.Println(getMemUsage(data))
	fmt.Println(getNICUsage(data, "eth1"))
	fmt.Println(getPeerSecsSinceHS(data, 9))
	fmt.Println(getPeerUsage(data, 12))

	if err != nil {
		fmt.Println(err)
	}
}
