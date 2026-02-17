package main

import (
	"flag"
	"log"

	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/database"
)

func main() {
	force := flag.Bool("force", false, "purge seed tables before inserting data")
	flag.Parse()

	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	defer func() {
		if err := database.CloseDB(db); err != nil {
			log.Printf("warning: close db: %v", err)
		}
	}()

	err := database.SeedDatabase(db, database.SeedOptions{
		Users:      24,
		Categories: 8,
		Products:   64,
		Orders:     50,
		Force:      *force,
	})
	if err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Println("seed completed successfully")
}
