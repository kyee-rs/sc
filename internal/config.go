package main

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     int
	DbPath   string
	BlockTor bool
	Gzip     bool
	CleanUp  int
	MaxSize  int
	Language string
}

func loadConfig() Config {
	v := viper.New()

	// Configure file loading
	v.AddConfigPath(".")
	v.AddConfigPath("cfg")
	v.AddConfigPath("ghost")
	v.AddConfigPath("/etc/ghost/")
	v.SetConfigName("cfg")

	// Configure environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace dashes and dots in env var names with underscores
	v.SetEnvPrefix("GS")                                         // Look for environment variables prefixed with GS_
	v.AutomaticEnv()                                             // Look for env vars for config keys
	v.AllowEmptyEnv(false)                                       // Consider defined environment variables with empty values

	// Set default values (in case none of the above config sources define a value for a certain key)
	v.SetDefault("host", "127.0.0.1")
	v.SetDefault("port", 8080)
	v.SetDefault("dbpath", "ghost_files.db")
	v.SetDefault("block_tor", false)
	v.SetDefault("gzip", true)
	v.SetDefault("cleanup", 0)
	v.SetDefault("maxsize", 0)
	v.SetDefault("language", "en")

	// Read and parse a config file
	// Ignore if error is File Not Found. Any other error is fatal.
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
