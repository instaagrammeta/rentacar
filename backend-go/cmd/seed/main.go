// Command seed initialises the database with default users, settings and demo
// data. Usage: `go run ./cmd/seed` (add -minimal for admin + settings only).
package main

import (
	"flag"
	"log"

	"github.com/instaagrammeta/rentacar/backend-go/internal/config"
	"github.com/instaagrammeta/rentacar/backend-go/internal/database"
	"github.com/instaagrammeta/rentacar/backend-go/internal/seed"
)

func main() {
	minimal := flag.Bool("minimal", false, "Only seed the admin account and settings")
	flag.Parse()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration: %v", err)
	}

	seed.Run(db, !*minimal)
	log.Println("Готово. База данных инициализирована.")
}
