package cmd

import (
	"context"
	"github.com/h44z/lightmigrate"
	"github.com/h44z/lightmigrate-mongodb/mongodb"
	"github.com/spf13/cobra"
	"linnoxlewis/tg-bot-eng-dictionary/internal/config"
	"linnoxlewis/tg-bot-eng-dictionary/internal/db"
	"log"
	"os"
	"time"
)

var path = "migrations"

var mgrCmd = &cobra.Command{
	Use: "migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30))
		defer cancel()
		mongoDb := db.Init(ctx, cfg)
		mongoClient := mongoDb.GetClient()

		fsys := os.DirFS("/app/data")

		source, err := lightmigrate.NewFsSource(fsys, path)
		if err != nil {
			log.Fatalf("unable to setup source: %v", err)
		}
		defer source.Close()

		driver, err := mongodb.NewDriver(mongoClient, cfg.GetMongoDatabase(),
			mongodb.WithLocking(mongodb.LockingConfig{
				Enabled: true,
			}))
		if err != nil {
			log.Fatalf("unable to setup driver: %v", err)
		}
		defer driver.Close()

		migrator, err := lightmigrate.NewMigrator(source, driver, lightmigrate.WithVerboseLogging(true))
		if err != nil {
			log.Fatalf("unable to setup migrator: %v", err)
		}

		err = migrator.Migrate(2)
		if err != nil {
			log.Fatalf("migration error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mgrCmd)
}
