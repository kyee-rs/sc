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
	Address         string `json:"address"`
	Port            int16  `json:"port"`
	Size_limit      int16  `json:"size_limit"`
	DB_path         string `json:"db_path"`
	Blocklist_path  string `json:"blocklist_path"`
	Index_page_path string `json:"index_page_path"`
	Block_TOR       bool   `json:"block_tor"`
	Fake_SSL        bool   `json:"fake_ssl"` // Use this if you are using a reverse proxy with SSL enabled. There is no need to specify cert and key files.
	SSL_            bool   `json:"ssl_"`     // Use this to use real SSL on this executable.
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
		return config_struct{
			Address:         "127.0.0.1",
			Port:            3000,
			Size_limit:      10,
			DB_path:         "files.sqlite",
			Blocklist_path:  "blocklist.txt",
			Block_TOR:       true,
			Index_page_path: "index.html",
			Fake_SSL:        true,
			SSL_:            false,
			SSL_cert:        "",
			SSL_key:         "",
			Gzip_:           true,
		}
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
	case r.URL.Path == "/favicon.ico":
		return
	default:
		home(w, r)
	}
}

func main() {
	host = flag.String("h", "127.0.0.1", "Address to serve on.")
	port = flag.Uint64("p", 3000, "Port to listen on.")
	config_file = flag.String("c", "example.config.jsonc", "Config file path.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./ghost -p=3000 -c=config.json -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	glog.Flush()

	config := loadConfig(*config_file)
	checkRequiredFlags(config)

	if config.Fake_SSL && config.SSL_ {
		glog.Errorf("Fake SSL and Real SSL cannot be enabled at the same time.")
		os.Exit(1)
	}

	if _, err := os.Stat(config.Index_page_path); os.IsNotExist(err) {
		glog.Errorf("Index page not found. Creating a new one.")
		file, err := os.Create(config.Index_page_path)
		file.Write([]byte("This is a default index page."))
		if err != nil {
			glog.Errorf("Error creating index page.")
		}
		file.Close()
	}
	parsed_tmpl, err := template.ParseFiles(config.Index_page_path)
	if err != nil {
		glog.Errorf("Error parsing index page template.")
	}
	tmpl = template.Must(parsed_tmpl, err)

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			glog.Error("Failed to create a database file! Exiting.")
			os.Exit(1)
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		glog.Error("Connection to database failed! Exiting.")
		os.Exit(1)
	}
	db.AutoMigrate(&Data{})
	rand.Seed(time.Now().Unix())
	if config.Gzip_ {
		glog.Info("GZIP Enabled.")
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
	if config.Fake_SSL || config.SSL_ {
		glog.Infof("Secure HTTPS server running on https://%s:%d", *host, *port)
		glog.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *host, *port), config.SSL_cert, config.SSL_key, nil))
	} else {
		glog.Infof("HTTP server running on http://%s:%d", *host, *port)
		glog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
	}
}
