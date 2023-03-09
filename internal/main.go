package main

import (
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"

  _ "embed"

  cron "github.com/go-co-op/gocron"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
)

var config = loadConfig()

// Delete files that are older than 1 week.
func cleanup(db *gorm.DB) {
  if config.AutoCleanUp != 0 {
    db.Where("created_at < ?", time.Now().Add(-1*24*time.Duration(config.AutoCleanUp)*time.Hour)).Delete(&Data{})
  } else {
    return
  }
}

func runCronJob(db *gorm.DB) {
  s := cron.NewScheduler(time.UTC)
  s.Every(1).Minute().Do(cleanup, db)
  s.StartAsync()
}

func main() {
  db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Silent),
  })

  if err != nil {
    log.Fatalln("Connection to database failed! Exiting.")
  }

  db.AutoMigrate(&Data{})

  defaultCheckers()

  runCronJob(db)

  e := echo.New()

  e.Use(middleware.Gzip())
  e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: "${method} - ${uri} - ${status} - ${latency_human}\n",
    Output: e.Logger.Output(),
  }))
  e.Use(middleware.Recover())
  e.Use(ipMiddleware())

  e.GET("/", func(c echo.Context) error {
    return c.HTML(http.StatusOK, fmt.Sprintf(`
<pre>
curl -F "file=@yourFile.txt" %s

WARNING: Each file gets deleted after 1 week.
</pre>`,
      fmt.Sprintf("%s://%s/\n", c.Scheme(), c.Request().Host)))
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
    return c.Blob(http.StatusOK, mime, bytes)
  })

  e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Host, config.Port)))
}