package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-ping/ping"
	"io"
	"log"
	//"net"
	"net/http"
	"time"
)

func getContainerIPs() ([]string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	ips := []string{}
	for _, container := range containers {
		for _, network := range container.NetworkSettings.Networks {
			ip := network.IPAddress
			if ip != "" {
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}

func pingContainers(ips []string) []PingInfo {
	var pingInfos []PingInfo
	for _, ip := range ips {
		info, err := pingSingleContainer(ip)
		if err != nil {
			continue
		}
		pingInfos = append(pingInfos, info)
	}
	return pingInfos
}

func pingSingleContainer(ip string) (PingInfo, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Printf("ERROR: cannot create pinger for %s: %v\n", ip, err)
		return PingInfo{}, err
	}
	pinger.Count = 1
	pinger.Timeout = 2 * time.Second

	pingTime := time.Now()
	pinger.Run()

	stats := pinger.Statistics()
	pingInfo := PingInfo{
		IPAddr:     stats.IPAddr.String(),
		PingTime:   pingTime.String(),
		PacketLoss: stats.PacketLoss,
	}
	fmt.Printf("Ping %s: Packet loss: %.2f%%, Avg: %v\n", ip, stats.PacketLoss, stats.AvgRtt)
	return pingInfo, nil
}

type PingInfo struct {
	IPAddr     string  `json:"ipAddr"`
	PingTime   string  `json:"pingTime"`
	PacketLoss float64 `json:"packetLoss"`
}

type Req struct {
	Stats []PingInfo `json:"stats"`
}

type Resp struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func SendRequest(url string, stats []PingInfo) (*Resp, error) {
	//requestData := Req{
	//	Stats: stats,
	//}
	//var requestData []PingInfo
	jsonBody, err := json.Marshal(stats)
	fmt.Println(string(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	body := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	clnt := &http.Client{}
	resp, err := clnt.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	fmt.Println(string(responseBody))

	var data Resp
	if err = json.Unmarshal(responseBody, &data); err != nil {
		return nil, err
	}
	return &data, nil

}

func main() {
	url := "http://backend:8080/pingInfo"
	ticker := time.NewTicker(2 * time.Minute)
	for range ticker.C {

		ips, err := getContainerIPs()
		if err != nil {
			log.Printf("ERROR: cannot get container IPs: %v\n", err)
		}
		if len(ips) == 0 {
			fmt.Println("No running containers found.")
		}
		stats := pingContainers(ips)
		//teststats := []PingInfo{{IPAddr: "172.17.0.2", PingTime: "2023-06-17T14:30:00Z", PacketLoss: 23.221}, {IPAddr: "102.37.1.5", PingTime: "2023-06-17T13:30:00Z", PacketLoss: 0.321}}

		resp, err := SendRequest(url, stats)
		if err != nil {
			log.Printf("ERROR: cannot send ping request: %v\n", err)
		}
		fmt.Println(resp)
	}

}
