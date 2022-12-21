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

func isBlocked(ip string, blocklist_map *os.File) bool {
	data := make([]byte, 1024)
	count, err := blocklist_map.Read(data)
	if err != nil {
		glog.Errorf("Error reading blocklist file")
		return false
	}

	// Check if the IP is in the blocklist
	if strings.Contains(string(data[:count]), ip) {
		return true
	}
	return false
}

// Dead simple router that just does the **perform** the job
func router(w http.ResponseWriter, r *http.Request) {
	config := loadConfig(*config_file)
	file, err := os.Open(config.Blocklist_path)
	if err != nil {
		glog.Errorf("Error opening %s", config.Blocklist_path)
	}
	defer file.Close()

	if isBlocked(r.RemoteAddr, file) {
		glog.Errorf("Blocked IP: %s", r.RemoteAddr)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "403: Forbidden")
		return
	}
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
		fmt.Fprintf(os.Stderr, "USAGE: ./gh0.st -p=8080 -c=config.go -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	glog.Flush()

	config := loadConfig(*config_file)

	tmpl = template.Must(template.ParseFiles(config.Index_page_path))

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

	// Load the config file
	if _, err := os.Stat(*config_file); os.IsNotExist(err) {
		glog.Errorf("Config file %s does not exist.", *config_file)
		os.Exit(1)
	}
	file, err := os.Open(*config_file)
	if err != nil {
		glog.Errorf("Error opening %s", *config_file)
		os.Exit(1)
	}
	defer file.Close()
	// Read the file
	data := make([]byte, 1024)
	count, err := file.Read(data)
	if err != nil && err != io.EOF {
		glog.Errorf("Error reading %s", *config_file)
		os.Exit(1)
	}
	// Check if the config file is valid
	if !strings.Contains(string(data[:count]), "package main") {
		glog.Errorf("Config file %s is not valid.", *config_file)
		os.Exit(1)
	}
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
