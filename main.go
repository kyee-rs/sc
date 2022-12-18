package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var port *uint64
var tmpl *template.Template
var host *string

// Dead simple router that just does the **perform** the job
func router(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.Header.Get("Content-type"), "multipart/form-data"):
		upload(w, r)
	case uuidMatch.MatchString(r.URL.Path):
		getFile(w, r)
	default:
		home(w, r)
	}
}

// Route handling, logging and application serving
func main() {
	if _, err := os.Stat("./db/files.sqlite"); os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create("./db/files.sqlite")
		if err != nil {
			panic("failed to create database file")
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open("./db/files.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Data{})

	// Random seed
	rand.Seed(time.Now().Unix())

	// Home template initalization
	tmpl = template.Must(template.ParseFiles("./templates/index.html"))
	// Flags for the leveled logging

	host = flag.String("h", "0.0.0.0", "Address to serve on")
	port = flag.Uint64("p", 8000, "port")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./0xg0.st -p=8080 -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	glog.Flush()

	// Routing
	http.HandleFunc("/", router)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
}
