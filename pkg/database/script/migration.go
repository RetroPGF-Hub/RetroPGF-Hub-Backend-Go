package main

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/database/migration"
	"context"
	"log"
	"os"
)

func main() {

	ctx := context.Background()
	_ = ctx
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 1 {
			log.Fatal("Error: .env path is invalid")
		}
		return os.Args[1]
	}())

	switch cfg.App.Name {
	case "project":
		// migration.Commentigrate(ctx, &cfg)
		// migration.FavMigrate(ctx, &cfg)
		migration.ProjectMigrate(ctx, &cfg)
	case "users":
		migration.UsersMigrate(ctx, &cfg)

	}

}
