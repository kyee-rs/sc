package main

import (
	"log"
	"os"
)

func defaultCheckers() {
	if _, err := os.Stat(config.DbPath); os.IsNotExist(err) {
		file, err := os.Create(config.DbPath)
		if err != nil {
			log.Println(ts.DatabaseErrors.CreationFailed)
			log.Fatalln(err)
		}
		if err := file.Close(); err != nil {
			log.Println(ts.DatabaseErrors.CloseFailed)
			log.Fatalln(err)
		}
	}
}
