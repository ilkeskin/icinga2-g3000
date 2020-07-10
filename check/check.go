package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	s "strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/ilkeskin/icinga-g3000/lib"
)

func getData(host string, port int, path string) lib.JSONSkeleton {
	resp, err := http.Get("http://" + host + ":" + s.Itoa(port) + path)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var skel lib.JSONSkeleton
	json.Unmarshal(body, &skel)

	return skel
}

func getCPUUsage(data lib.JSONSkeleton) string {
	m := structs.Map(data.CPU)
	var result string
	for k, v := range m {
		result += fmt.Sprintf("'%s'=%s%% ",
			strings.ToLower(k),
			s.FormatFloat(v.(float64), 'f', 2, 64))
	}
	return strings.TrimSpace(result)
}

func getMemUsage(data lib.JSONSkeleton) string {
	m := structs.Map(data.Memory)
	var result string
	for k, v := range m {
		result += fmt.Sprintf("'%s'=%s%% ",
			strings.ToLower(k),
			s.FormatFloat(v.(float64), 'f', 2, 64))
	}
	return strings.TrimSpace(result)
}

func getNICUsage(data lib.JSONSkeleton, nicname string) string {
	var m map[string]interface{}
	for i := range data.Network {
		if data.Network[i].Name == nicname {
			m = structs.Map(data.Network[i])
		}
	}

	return fmt.Sprintf("'down'=%skbps 'up'=%skbps",
		s.FormatFloat(m["Rx"].(float64), 'f', 2, 64),
		s.FormatFloat(m["Tx"].(float64), 'f', 2, 64))
}

func getPeerSecsSinceHS(data lib.JSONSkeleton, peerIndex int64) string {
	for i := range data.Wireguard {
		ip := strings.Split(data.Wireguard[i].IntIPAddr, "/")
		ip = strings.Split(ip[0], ".")
		idx, _ := s.ParseInt(ip[len(ip)-1], 10, 8)
		if idx == peerIndex {
			return fmt.Sprintf("'lasths'=%ds", time.Now().Unix()-data.Wireguard[i].LastHS)
		}
	}
	return "Error"
}

func getPeerUsage(data lib.JSONSkeleton, peerIndex int64) string {
	for i := range data.Wireguard {
		ip := strings.Split(data.Wireguard[i].IntIPAddr, "/")
		ip = strings.Split(ip[0], ".")
		idx, _ := s.ParseInt(ip[len(ip)-1], 10, 8)
		if idx == peerIndex {
			return fmt.Sprintf("'down'=%.2fkbps 'up'=%.2fkbps",
				data.Wireguard[i].PeerRate.Rx,
				data.Wireguard[i].PeerRate.Tx)
		}
	}
	return "Error"
}

func main() {
	data := getData("127.0.0.1", 8888, "/test.json")

	fmt.Println(getCPUUsage(data))
	fmt.Println(getMemUsage(data))
	fmt.Println(getNICUsage(data, "eth1"))
	fmt.Println(getPeerSecsSinceHS(data, 9))
	fmt.Println(getPeerUsage(data, 12))
}
