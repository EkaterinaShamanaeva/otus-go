package main

import (
	"flag"
	"fmt"
	configuration "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"os"
)

var (
	flagSet    = flag.NewFlagSet("goose", flag.ExitOnError)
	dir        = flagSet.String("dir", ".", "dir with migration sql files")
	configFile = flagSet.String("config", "configs/config.yaml", "path with config file")
)

func main() {
	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if len(args) < 1 {
		flagSet.Usage()
	}

	config := configuration.NewConfig()

	if err := config.BuildConfig(*configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if config.Storage != "SQL" {
		log.Fatalf("in memory storage is used")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.Database.Username,
		config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name,
		config.Database.SSLMode)

	driver := "postgres"

	db, err := goose.OpenDBWithDriver(driver, dsn)
	if err != nil {
		log.Fatalf("failed to open DB with the error: %v", err)
	}

	var arguments []string
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err = goose.Run(args[0], db, *dir, arguments...); err != nil {
		log.Fatalf("goose migrator run: %v", err)
	}
}
