package main

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type ContainerInfo struct {
	Name       string `json:"Name"`
	IPAddress  string `json:"NetworkSettings.Networks.bridge.IPAddress"`
	MacAddress string `json:"NetworkSettings.MacAddress"`
}

func GetContainerIPs() ([]string, []ContainerInfo) {
	cmd := exec.Command("docker", "inspect", "-f", "{{.Name}}:{{.NetworkSettings.Networks.bridge.IPAddress}}", "$(docker ps -q)")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Ошибка выполнения команды: %v", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	containerInfos := make([]ContainerInfo, 0)
	ipAddresses := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		name := parts[0]
		ipAddress := parts[1]

		info := ContainerInfo{Name: name, IPAddress: ipAddress}
		containerInfos = append(containerInfos, info)
		ipAddresses = append(ipAddresses, ipAddress)
	}

	return ipAddresses, containerInfos
}
