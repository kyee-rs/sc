package main

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host           string
	Port           int
	Size_limit     int
	DB_path        string
	Blocklist_path string
	Index_path     string
	Block_TOR      bool
	Fake_SSL       bool // Use this if you are using a reverse proxy with SSL enabled. There is no need to specify cert and key files.
	SSL_           bool `mapstructure:"enable_ssl"` // Use this to use real SSL on this executable.
	SSL_cert       string
	SSL_key        string
	Gzip_          bool `mapstructure:"enable_gzip"`
}

func loadConfig() Config {
	v := viper.New()
	// Configure file loading
	v.AddConfigPath(".")
	v.AddConfigPath("config")
	v.AddConfigPath("/etc/ghost/config")
	v.SetConfigName("config")
	// Configure environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace dashes and with underscores
	v.SetEnvPrefix("ghost")                                      // Look for environment variables prefixed with GHOST_
	v.AutomaticEnv()                                             // Look for env vars for config keys
	v.AllowEmptyEnv(true)                                        // Consider defined environment variables with empty values

	// Set default values (in case none of the above config sources define a value for a certain key)
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("port", 8080)
	v.SetDefault("size_limit", 10)
	v.SetDefault("db_path", "files.db")
	v.SetDefault("blocklist_path", "blocklist.txt")
	v.SetDefault("index_path", "index.html")
	v.SetDefault("block_tor", true)
	v.SetDefault("fake_ssl", true)
	v.SetDefault("enable_ssl", false)
	v.SetDefault("ssl_cert", "cert.pem")
	v.SetDefault("ssl_key", "key.pem")
	v.SetDefault("enable_gzip", true)

	// Read and parse a config file
	// Ignore file not found errors
	err := v.ReadInConfig()
	if _, configFileNotFound := err.(viper.ConfigFileNotFoundError); err != nil && !configFileNotFound {
		panic(err)
	}

	var config Config

	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}
