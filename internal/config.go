package main

type configuration struct {
	Server serverSettings      `hcl:"server,block"`
	Logger loggerSettings      `hcl:"logger,block"`
	Limits limitationsSettings `hcl:"limits,block"`
}

type serverSettings struct {
	ServerName  string `hcl:"server_name,label"`
	AppName     string `hcl:"app_name,label"`
	Port        int    `hcl:"port"`
	DatabaseURL string `hcl:"database_url"`
	Seed        uint64 `hcl:"seed"`
}

type limitationsSettings struct {
	MaxSize     int      `hcl:"max_size"`
	BlockTor    bool     `hcl:"block_tor"`
	IPBlacklist []string `hcl:"ip_blacklist"`
}

type loggerSettings struct {
	ForceColors   bool `hcl:"force_colors"`
	FullTimestamp bool `hcl:"full_timestamp"`
}
