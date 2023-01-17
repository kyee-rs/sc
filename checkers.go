package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func defaultCheckers() {
	if _, err := os.Stat(config.Index_path); os.IsNotExist(err) {
		fmt.Printf("No index page found, creating a new one.\n")
		out, err := os.Create(config.Index_path)
		if err != nil {
			fmt.Printf("Error creating index page file.")
		}
		defer out.Close()
		resp, err := http.Get("https://raw.githubusercontent.com/voxelin/gh0.st/master/templates/index.html")
		if err != nil {
			fmt.Printf("Error loading default index page.\n")
		}
		defer resp.Body.Close()
		if _, err = io.Copy(out, resp.Body); err != nil {
			fmt.Printf("Error writing default index page.\n")
		}
	}

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			fmt.Printf("Failed to create a database file! Exiting.\n")
			os.Exit(1)
		}
		file.Close()
	}
}
