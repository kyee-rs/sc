package main

import (
	"log"
	"os"
)

func defaultCheckers() {
	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			log.Fatalln("Failed to create database file")
		}
		file.Close()
	}
}
