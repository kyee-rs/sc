package main

import (
	"compress/gzip"
	"encoding/json"
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
var config_file *string

type config_struct struct {
	Size_limit      int16  `json:"size_limit"`
	DB_path         string `json:"db_path"`
	Blocklist_path  string `json:"blocklist_path"`
	Index_page_path string `json:"index_page_path"`
	SSL_            bool   `json:"ssl_"`
	SSL_cert        string `json:"ssl_cert"`
	SSL_key         string `json:"ssl_key"`
	Gzip_           bool   `json:"gzip_"`
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func loadConfig(file string) config_struct {
	var config config_struct
	jsonFile, err := os.Open(file)
	if err != nil {
		glog.Errorf("Error opening %s. Consider specifying the config file using -c flag.", file)
		os.Exit(1)
	}
	bv, _ := io.ReadAll(jsonFile)
	json.Unmarshal(bv, &config)
	return config
}

func checkRequiredFlags(config config_struct) {
	if len(config.Index_page_path) <= 0 || len(config.Blocklist_path) <= 0 || len(config.DB_path) <= 0 {
		glog.Errorf("Some required flags are missing. Please check the config file.")
		os.Exit(1)
	}
}

func router(w http.ResponseWriter, r *http.Request) {
	config := loadConfig(*config_file)
	file, err := os.Open(config.Blocklist_path)
	if err != nil {
		glog.Errorf("Error opening %s", config.Blocklist_path)
	}
	defer file.Close()

	// Check if the IP is blocked -----------------------------------------------
	ip, err := getIP(r)
	if err != nil {
		glog.Errorf("Error getting IP: %s", err)
	}
	if isBlocked(ip, file) || isTorExitNode(ip) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "403: Forbidden. Your IP (%s) has been blocked.", ip)
		return
	}
	// --------------------------------------------------------------------------

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
	// Flags for the leveled logging
	host = flag.String("h", "0.0.0.0", "Address to serve on")
	port = flag.Uint64("p", 8000, "port")
	config_file = flag.String("c", "ghost.config.json", "config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./gh0.st -p=8080 -c=config.json -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	glog.Flush()

	config := loadConfig(*config_file)
	checkRequiredFlags(config)

	parsed_tmpl, err := template.ParseFiles(config.Index_page_path)
	if err != nil {
		glog.Errorf("Error parsing index page template. Make sure that %s exists.", config.Index_page_path)
		os.Exit(1)
	}
	tmpl = template.Must(parsed_tmpl, err)

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create(config.DB_path)
		if err != nil {
			panic("failed to create database file")
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Data{})

	// Random seed
	rand.Seed(time.Now().Unix())

	// Gzip support
	if config.Gzip_ {
		glog.Info("Gzip module loaded.")
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

	// SSL support
	if config.SSL_ {
		glog.Infof("SSL-HTTP server running on https://%s:%d", *host, *port)
		glog.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), config.SSL_cert, config.SSL_key, nil))
	} else {
		glog.Infof("HTTP server running on http://%s:%d", *host, *port)
		glog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
	}
}
