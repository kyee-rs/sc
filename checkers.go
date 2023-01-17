package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func defaultCheckers() {
	if _, err := os.Stat(config.Index_path); os.IsNotExist(err) {

		out, err := os.Create(config.Index_path)
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()
		resp, err := http.Get("https://raw.githubusercontent.com/voxelin/gh0.st/master/templates/index.html")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		if _, err = io.Copy(out, resp.Body); err != nil {
			fmt.Println(err)
		}
	}

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		file.Close()
	}
}
