package main

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var config = loadConfig()

func main() {
	defaultCheckers() // Run IF statements to check if index.html and db.sqlite3 exist

	router := gin.New()

	if config.Gzip_ {
		router.Use(gzip.Gzip(gzip.DefaultCompression)) // Enable GZIP if config.Gzip_ is true
	}
	router.Use(ipMiddleware()) // Check if client ip is in blacklist or is a Tor exit node
	router.Use(logger.SetLogger())
	router.Use(gin.Recovery())

	if len(config.Allowed_IPs) > 0 {
		router.SetTrustedProxies(config.Allowed_IPs) // Set trusted proxies
	} else {
		router.SetTrustedProxies(nil)
	}
	if config.Trusted_Platform != "" {
		switch config.Trusted_Platform {
		case "cloudflare":
			router.TrustedPlatform = gin.PlatformCloudflare
		case "google":
			router.TrustedPlatform = gin.PlatformGoogleAppEngine
		default:
			router.TrustedPlatform = config.Trusted_Platform
		}
	}

	router.LoadHTMLFiles(config.Index_path) // Load index.html template
	var scheme string

	if config.SSL_ || config.Fake_SSL {
		scheme = "https"
	} else {
		scheme = "http"
	}

	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{
		Logger:                                   gormlogger.Default.LogMode(gormlogger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	}) // Connect to database file and run migrations.
	if err != nil {
		panic("Connection to database failed! Exiting.\n")
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
			fmt.Println(err)
		}
	} else {
		if err := router.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
			fmt.Println(err)
		}
	}
}
