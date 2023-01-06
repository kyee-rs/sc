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

	"github.com/apex/log"
	"github.com/tailscale/hujson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tmpl *template.Template
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
func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

func loadConfig(file string) config_struct {
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "loadConfig",
		"file":    "main.go",
	})
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
	jsonc, err := io.ReadAll(jsonFile)
	if err != nil {
		logger.Errorf("Error reading config file: %s", err)
		os.Exit(1)
	}
	bv, _ := standardizeJSON(jsonc)
	json.Unmarshal(bv, &config)
	defer jsonFile.Close()
	return config
}

func checkRequiredFlags(config config_struct) {
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "checkRequiredFlags",
		"file":    "main.go",
	})
	if len(config.Index_page_path) <= 0 || len(config.Blocklist_path) <= 0 || len(config.DB_path) <= 0 {
		logger.Errorf("Some required flags are missing. Please check the config file.")
		os.Exit(1)
	}
}

func router(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "router",
		"file":    "main.go",
	})
	config := loadConfig(*config_file)
	file, err := os.Open(config.Blocklist_path)
	if err != nil {
		logger.Errorf("Error opening %s", config.Blocklist_path)
	}
	defer file.Close()

	// Check if the IP is blocked -----------------------------------------------
	ip, err := getIP(r)
	if err != nil {
		logger.Errorf("Error getting IP: %s", err)
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
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "main",
		"file":    "main.go",
	})
	config_file = flag.String("c", "config.jsonc", "Config file path.")
	level := flag.Uint64("level", 2, "Log level. 1 - DebugLevel 2 - InfoLevel 3 - WarnLevel 4 - ErrorLevel 5 - FatalLevel")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./ghost -c=config.jsonc -level=2\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	switch *level {
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.ErrorLevel)
	case 5:
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	config := loadConfig(*config_file)
	checkRequiredFlags(config)

	if config.Fake_SSL && config.SSL_ {
		logger.Errorf("Fake SSL and Real SSL cannot be enabled at the same time.")
		os.Exit(1)
	}

	if _, err := os.Stat(config.Index_page_path); os.IsNotExist(err) {
		logger.Errorf("Index page not found. Creating a new one.")
		file, err := os.Create(config.Index_page_path)
		file.Write([]byte("This is a default index page."))
		if err != nil {
			logger.Errorf("Error creating index page.")
		}
		file.Close()
	}
	parsed_tmpl, err := template.ParseFiles(config.Index_page_path)
	if err != nil {
		logger.Errorf("Error parsing index page template.")
	}
	tmpl = template.Must(parsed_tmpl, err)

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			logger.Error("Failed to create a database file! Exiting.")
			os.Exit(1)
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		logger.Error("Connection to database failed! Exiting.")
		os.Exit(1)
	}
	db.AutoMigrate(&Data{})
	rand.Seed(time.Now().Unix())
	if config.Gzip_ {
		logger.Info("GZIP Enabled.")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				logger.Errorf("Error while compressing: %s", err)
				return
			}
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			router(gzr, r)
		})
	} else {
		http.HandleFunc("/", router)
	}
	if config.SSL_ {
		logger.Infof("Secure HTTPS server running on https://%s:%d", config.Address, config.Port)
		http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Address, config.Port), config.SSL_cert, config.SSL_key, nil)
	} else if config.Fake_SSL {
		logger.Infof("Secure HTTPS server running on https://%s:%d", config.Address, config.Port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.Address, config.Port), nil)
	} else {
		logger.Infof("HTTP server running on http://%s:%d", config.Address, config.Port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.Address, config.Port), nil)
	}
}
