package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"html/template"
	"io"
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

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

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
	if _, err := os.Stat(db_path); os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create(db_path)
		if err != nil {
			panic("failed to create database file")
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
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
		fmt.Fprintf(os.Stderr, "USAGE: ./gh0.st -p=8080 -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	glog.Flush()

	// SSL support
	if ssl_ {
		glog.Infof("Serving on https://%s:%d", *host, *port)
		glog.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), ssl_cert, ssl_key, nil))
	} else {
		glog.Infof("Serving on http://%s:%d", *host, *port)
		glog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
	}

	// Use gzip compression
	if gzip_ {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				glog.Errorf("Error while compressing: %s", err)
				return
			}
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			router(gzr, r)
		})
	} else {
		http.HandleFunc("/", router)
	}
}
