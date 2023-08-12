package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	log "github.com/sirupsen/logrus"
	timescale "github.com/voxelin/sc/sqlc_gen"

	_ "embed"

	"github.com/bytedance/sonic"
	cron "github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
	"github.com/teris-io/shortid"
)

var (
	ctx    = context.Background()
	sid    *shortid.Shortid
	db     *timescale.Queries
	config Configuration
)

func cleanup() {
	if err := db.PurgeFiles(ctx); err != nil {
		log.Fatalln("Failed to execute cleanup")
	}
}

func runCronJob() {
	s := cron.NewScheduler(time.UTC)
	if _, err := s.Every(12).Hours().Do(cleanup); err != nil {
		log.Println("Failed to run a CronJob Task.")
		log.Fatal(err)
	}
	s.StartAsync()
}

func main() {
	config.load()

	log.SetFormatter(&log.TextFormatter{ForceColors: config.Logger.ForceColors, FullTimestamp: config.Logger.FullTimestamp, TimestampFormat: time.RFC822})
	log.SetOutput(os.Stdout)

	log.Println(config)

	dbInternal, err := sql.Open("postgres", config.Server.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open a TimescaleDB Connection: %v", err)
	}

	db = timescale.New(dbInternal)

	genSid, err := shortid.New(1, shortid.DefaultABC, config.Server.Seed)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "init(shortid.New)",
		}).Fatal(err)
	}

	sid = genSid

	if !fiber.IsChild() {
		runCronJob()
	}

	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ServerHeader: config.Server.ServerName,
		AppName:      config.Server.AppName,
		Prefork:      true,
		BodyLimit:    config.Limits.MaxSize * 1024 * 1024,
	})

	app.Use(
		compress.New(
			compress.Config{
				Level: compress.LevelBestSpeed,
			},
		),
		favicon.New(),
	)

	if config.Limits.BlockTor {
		app.Use(torMiddleware())
	}
	if len(config.Limits.IpBlacklist) > 0 {
		app.Use(ipMiddleware())
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(403)
	})

	app.Post("/", upload)

	app.Get("/:id", loadResponse)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.Server.Port)))
}
