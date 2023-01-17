package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var config = loadConfig()

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

// Initialize loggers
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

func main() {
	defaultCheckers() // Run IF statements to check if index.html and db.sqlite3 exist

	router := gin.Default()

	if config.Gzip_ {
		router.Use(gzip.Gzip(gzip.DefaultCompression)) // Enable GZIP if config.Gzip_ is true
	}
	router.Use(ipMiddleware()) // Check if client ip is in blacklist or is a Tor exit node

	router.LoadHTMLFiles(config.Index_path) // Load index.html template
	var scheme string

	if config.SSL_ || config.Fake_SSL {
		scheme = "https"
	} else {
		scheme = "http"
	}

	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{}) // Connect to database file and run migrations.
	if err != nil {
		ErrorLogger.Printf("Connection to database failed! Exiting.")
		os.Exit(1)
	}
	db.AutoMigrate(&Data{})

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"host": scheme + "://" + c.Request.Host, // replace {{.host}} with the current url
		})
	})

	router.POST("/", func(c *gin.Context) {
		upload(c, db, scheme)
	})

	router.GET("/:uuid", func(c *gin.Context) {
		if c.Param("uuid") == "favicon.ico" {
			return
		}
		bytes, mime, filename := getFile(c.Param("uuid"))
		if bytes == nil {
			c.String(404, "404: File not found!\n")
			return
		}

		c.Header("Content-Disposition", "filename="+filename)
		c.Data(200, mime, bytes)
	})

	// Start server with TLS if config.SSL_ is true.
	if config.SSL_ {
		if err := router.RunTLS(fmt.Sprintf("%s:%d", config.Host, config.Port), config.SSL_cert, config.SSL_key); err != nil {
			ErrorLogger.Printf("Error while starting server: %s", err)
		}
	} else {
		if err := router.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
			ErrorLogger.Printf("Error while starting server: %s", err)
		}
	}
}
