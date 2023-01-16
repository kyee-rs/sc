package main

import (
	"compress/gzip"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tmpl *template.Template
var config = loadConfig()
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger.SetOutput(os.Stderr)
	WarningLogger.SetOutput(os.Stdout)
	InfoLogger.SetOutput(os.Stdout)
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func router(w http.ResponseWriter, r *http.Request) {
	ip, err := getIP(r)
	if err != nil {
		ErrorLogger.Printf("Error getting IP: %s", err)
	}
	if config.Blocklist_path != "" {
		if file, err := os.Open(config.Blocklist_path); err == nil {
			if isBlocked(ip, file) {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "403: Forbidden. Your IP (%s) has been blocked.", ip)
				return
			}
			defer file.Close()
		}
	}

	if config.Block_TOR {
		if isTorExitNode(ip) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "403: Forbidden. Your IP (%s) has been blocked.", ip)
			return
		}
	}
	// ---------------------------------------------------------------------------------------------------------------------

	switch {
	case strings.Contains(r.Header.Get("Content-type"), "multipart/form-data"):
		upload(w, r)
	case uuidMatch.MatchString(r.URL.Path):
		getFile(w, r)
	default:
		home(w, r)
	}
}

func main() {
	if config.Fake_SSL && config.SSL_ {
		ErrorLogger.Printf("Fake SSL and Real SSL cannot be enabled at the same time.")
		os.Exit(1)
	}

	if _, err := os.Stat(config.Index_path); os.IsNotExist(err) {
		ErrorLogger.Printf("Index page not found. Creating a new one.")
		file, err := os.Create(config.Index_path)
		file.Write([]byte("This is a default index page."))
		if err != nil {
			ErrorLogger.Printf("Error creating index page.")
		}
		file.Close()
	}
	parsed_tmpl, err := template.ParseFiles(config.Index_path)
	if err != nil {
		ErrorLogger.Printf("Error parsing index page template.")
	}

	tmpl = template.Must(parsed_tmpl, err)

	if _, err := os.Stat(config.DB_path); os.IsNotExist(err) {
		file, err := os.Create(config.DB_path)
		if err != nil {
			ErrorLogger.Printf("Failed to create a database file! Exiting.")
			os.Exit(1)
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		ErrorLogger.Printf("Connection to database failed! Exiting.")
		os.Exit(1)
	}
	db.AutoMigrate(&Data{})
	rand.Seed(time.Now().Unix())
	if config.Gzip_ {
		InfoLogger.Printf("GZIP Enabled.")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				ErrorLogger.Printf("Error while compressing: %s", err)
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
		InfoLogger.Printf("Secure HTTPS server running on https://%s:%d", config.Host, config.Port)
		http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Host, config.Port), config.SSL_cert, config.SSL_key, nil)
	} else if config.Fake_SSL {
		InfoLogger.Printf("Secure HTTPS server running on https://%s:%d", config.Host, config.Port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
	} else {
		InfoLogger.Printf("HTTP server running on http://%s:%d", config.Host, config.Port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
	}

}
