package main

import (
	"io"
	"net/http"
	"os"
)

func defaultCheckers() {
	if _, err := os.Stat(config.Index_path); os.IsNotExist(err) {
		ErrorLogger.Printf("No index page found, creating a new one.")
		out, err := os.Create(config.Index_path)
		if err != nil {
			ErrorLogger.Printf("Error creating index page file.")
		}
		defer out.Close()
		resp, err := http.Get("https://raw.githubusercontent.com/voxelin/gh0.st/master/templates/index.html")
		if err != nil {
			ErrorLogger.Printf("Error loading default index page.")
		}
		defer resp.Body.Close()
		if _, err = io.Copy(out, resp.Body); err != nil {
			ErrorLogger.Printf("Error writing default index page.")
		}
	}

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			ErrorLogger.Printf("Failed to create a database file! Exiting.")
			os.Exit(1)
		}
		file.Close()
	}
}
