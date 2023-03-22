package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "embed"

	"github.com/fatih/color"
	cron "github.com/go-co-op/gocron"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed html/index.html
var html []byte

var tmpl = template.Must(template.New("index").Parse(string(html)))

var config = loadConfig()
var banner = `
   ________               __     _____                          
  / ____/ /_  ____  _____/ /_   / ___/___  ______   _____  _____
 / / __/ __ \/ __ \/ ___/ __/   \__ \/ _ \/ ___/ | / / _ \/ ___/
/ /_/ / / / / /_/ (__  ) /_    ___/ /  __/ /   | |/ /  __/ /    
\____/_/ /_/\____/____/\__/   /____/\___/_/    |___/\___/_/ 

`

func cleanup(db *gorm.DB) {
	if config.AutoCleanUp != 0 {
		db.Where("created_at < ?", time.Now().UTC().Add(-1*24*time.Duration(config.AutoCleanUp)*time.Hour)).Delete(&Data{})
	} else {
		return
	}
}

func runCronJob(db *gorm.DB) {
	s := cron.NewScheduler(time.UTC)
	if _, err := s.Every(1).Minute().Do(cleanup, db); err != nil {
		log.Fatalln("Failed to start cron job! Exiting.")
	}
	s.StartAsync()
}

func main() {
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalln("Connection to database failed! Exiting.")
	}

	if err := db.AutoMigrate(&Data{}); err != nil {
		log.Fatalln("Migration failed! Exiting.")
	}

	defaultCheckers()

	runCronJob(db)

	e := echo.New()

	e.HideBanner = true
	color.Cyan(banner)

	e.Use(middleware.Gzip())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} - ${uri} - ${status} - ${latency_human}\n",
		Output: e.Logger.Output(),
	}))
	e.Use(middleware.Recover())
	e.Use(ipMiddleware())

	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		
		return tmpl.Execute(c.Response(), map[string]interface{}{
			"host":      fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host),
			"retention": config.AutoCleanUp,
			"tor":       config.Block_TOR,
			"maxsize":   config.MaxSize,
		})
	})

	e.POST("/", func(c echo.Context) error {
		return upload(c, db)
	})

	e.GET("/favicon.ico", func(c echo.Context) error {
		if (c.Request().Header.Get("Accept")) == "application/json" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   true,
				"status":  http.StatusNotFound,
				"message": "404: File not found.",
			})
		}

		return c.Blob(http.StatusOK, "image/x-icon", []byte{})
	})

	e.GET("/:id", func(c echo.Context) error {
		bytes, filename, mime := getFile(c.Param("id"), db)
		if bytes == nil {
			return jsonOrString(c, http.StatusNotFound, "404: File not found.", true)
		}

		c.Response().Header().Set("Content-Disposition", "filename="+filename)
		c.Response().Header().Set("Content-Type", mime)
		
		return c.Blob(http.StatusOK, mime, bytes)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Host, config.Port)))
}
