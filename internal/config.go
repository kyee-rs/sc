package main

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host        string
	Port        int
	DB_path     string
	Block_TOR   bool
	Gzip        bool
	AutoCleanUp int
	MaxSize     int
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
	v.AllowEmptyEnv(false)                                       // Consider defined environment variables with empty values

	// Set default values (in case none of the above config sources define a value for a certain key)
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("port", 8080)
	v.SetDefault("db_path", "files.db")
	v.SetDefault("block_tor", false)
	v.SetDefault("gzip", true)
	v.SetDefault("autocleanup", 0)
  v.SetDefault("maxsize", 0)

	// Read and parse a config file
	// Ignore file not found errors
	err := v.ReadInConfig()
	if _, configFileNotFound := err.(viper.ConfigFileNotFoundError); err != nil && !configFileNotFound {
		log.Fatalln(err)
	}

	var config Config

	if err := v.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	return config
}
