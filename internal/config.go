package main

import (
	"strings"

	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	Server serverSettings      `koanf:"server"`
	Logger loggerSettings      `koanf:"logger"`
	Limits limitationsSettings `koanf:"limits"`
}

type serverSettings struct {
	ServerName  string `koanf:"serverName"`
	AppName     string `koanf:"appName"`
	Port        int    `koanf:"port"`
	DatabaseURL string `koanf:"databaseUrl"`
	Seed        uint64 `koanf:"seed"`
}

type limitationsSettings struct {
	MaxSize     int      `koanf:"maxSize"`
	BlockTor    bool     `koanf:"blockTor"`
	IpBlacklist []string `koanf:"ipBlacklist"`
}

type loggerSettings struct {
	ForceColors   bool `koanf:"forceColors"`
	FullTimestamp bool `koanf:"fullTimestamp"`
}

var k = koanf.New(".")

func (c *Configuration) load() {
	k.Load(confmap.Provider(map[string]interface{}{ //nolint:errcheck
		"server.serverName":    "Simple Cache v1.2.4",
		"server.appName":       "Simple Cache",
		"server.port":          8080,
		"server.databaseURL":   "postgres://user:password@localhost:5432/db",
		"server.seed":          "3719",
		"logger.forceColors":   false,
		"logger.fullTimestamp": true,
		"limits.maxSize":       10,
		"limits.blockTor":      false,
	}, "."), nil)

	k.Load(file.Provider("./config.hcl"), hcl.Parser(true)) //nolint:errcheck

	k.Load(env.Provider("SC_", ".", func(s string) string { //nolint:errcheck
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "SC_")), "_", ".", -1)
	}), nil)

	if err := k.Unmarshal("", &c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Failed to unmarshal the configuration.")
	}
}
