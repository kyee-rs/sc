package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/prophittcorey/tor"
)

func isTorExitNode(address string) bool {
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "isTorExitNode",
		"file":    "ip_blockers.go",
	})
	res, err := tor.IsExitNode(address)
	if err != nil {
		logger.Warnf("Error checking if %s is a Tor exit node: %s", address, err)
	}
	if res {
		logger.Warnf("%s is a Tor exit node. Acess denied.", address)
		return true
	}
	return false
}

func isBlocked(ip string, blocklist_map *os.File) bool {
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "isBlocked",
		"file":    "ip_blockers.go",
	})
	data := make([]byte, 1024)
	count, err := blocklist_map.Read(data)
	if err != nil {
		return false
	}

	// Check if the IP is in the blocklist
	if strings.Contains(string(data[:count]), ip) {
		logger.Warnf("%s is in a block-list.", ip)
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
