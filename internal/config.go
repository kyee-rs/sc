package main

type configuration struct {
	Port        int64          `hcl:"port"`
	DatabaseURL string         `hcl:"database_url"`
	Seed        uint64         `hcl:"seed"`
	Logger      loggerSettings `hcl:"logger,block"`
	Limits      limitations    `hcl:"limits,block"`
}

type limitations struct {
	MaxSize     int64    `hcl:"max_size"`
	BlockTor    bool     `hcl:"block_tor"`
	IPBlacklist []string `hcl:"ip_blacklist"`
}

type loggerSettings struct {
	ForceColors   bool `hcl:"force_colors"`
	FullTimestamp bool `hcl:"full_timestamp"`
}
