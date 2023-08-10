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
	"github.com/hashicorp/hcl/v2/hclsimple"
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
	config configuration
)

func cleanup(db *timescale.Queries) {
	if err := db.PurgeFiles(ctx); err != nil {
		log.Fatalln("Failed to execute cleanup")
	}
}

func runCronJob(db *timescale.Queries) {
	s := cron.NewScheduler(time.UTC)
	if _, err := s.Every(12).Hours().Do(cleanup, db); err != nil {
		log.Println("Failed to run a CronJob Task.")
		log.Fatal(err)
	}
	s.StartAsync()
}

func init() {
	err := hclsimple.DecodeFile("config.hcl", nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	log.SetFormatter(&log.TextFormatter{ForceColors: config.Logger.ForceColors, FullTimestamp: config.Logger.FullTimestamp, TimestampFormat: time.RFC822})
	log.SetOutput(os.Stdout)

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

	runCronJob(db)
}

func main() {
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
		torMiddleware(),
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(403)
	})

	app.Post("/", upload)

	app.Get("/:id", loadResponse)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.Server.Port)))
}
