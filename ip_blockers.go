package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prophittcorey/tor"
)

func isTorExitNode(address string) bool {
	res, err := tor.IsExitNode(address)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if res {
		return true
	}
	return false
}

func isBlacklisted(ip string, blacklist_map *os.File) bool {
	data := make([]byte, 1024)
	count, err := blacklist_map.Read(data)
	if err != nil {
		return false
	}

	// Check if the IP is in the blacklist
	if strings.Contains(string(data[:count]), ip) {
		return true
	}
	return false
}

func ipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if config.Block_TOR {
			if isTorExitNode(ip) {
				c.AbortWithStatus(403)
				return
			}
		}
		if config.Blacklist_path != "" {
			if file, err := os.Open(config.Blacklist_path); err == nil {
				if isBlacklisted(ip, file) {
					c.AbortWithStatus(403)
					return
				}
				defer file.Close()
			}
		}
		c.Next()
	}
}
