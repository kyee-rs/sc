package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/prophittcorey/tor"
)

func isTorExitNode(address string) bool {
	res, err := tor.IsExitNode(address)
	if err != nil {
		WarningLogger.Printf("Error checking if %s is a Tor exit node: %s", address, err)
	}
	if res {
		WarningLogger.Printf("%s is a Tor exit node. Access denied.", address)
		return true
	}
	return false
}

func isBlocked(ip string, blocklist_map *os.File) bool {
	data := make([]byte, 1024)
	count, err := blocklist_map.Read(data)
	if err != nil {
		return false
	}

	// Check if the IP is in the blocklist
	if strings.Contains(string(data[:count]), ip) {
		WarningLogger.Printf("%s is in a block-list.", ip)
		return true
	}
	return false
}

func getIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("NO VALID IP FOUND")
}
