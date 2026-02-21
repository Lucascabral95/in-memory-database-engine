package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/database"
)

func parseSeedFlags(args []string) (bool, error) {
	fs := flag.NewFlagSet("seed", flag.ContinueOnError)
	force := fs.Bool("force", false, "purge seed tables before inserting data")
	if err := fs.Parse(args); err != nil {
		return false, err
	}
	return *force, nil
}

func defaultSeedOptions(force bool) database.SeedOptions {
	return database.SeedOptions{
		Users:      24,
		Categories: 8,
		Products:   64,
		Orders:     50,
		Force:      force,
	}
}

func main() {
	force, err := parseSeedFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("seed flags error: %v", err)
	}

	cfg := config.LoadConfig()
	db, err := database.InitDB(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	defer func() {
		if err := database.CloseDB(db); err != nil {
			log.Printf("warning: close db: %v", err)
		}
	}()

	err = database.SeedDatabase(db, defaultSeedOptions(force))
	if err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Println("seed completed successfully")
}
