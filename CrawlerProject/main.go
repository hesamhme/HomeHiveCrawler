package main

import (
	"CrawlerProject/internal/repository"
	"CrawlerProject/internal/service"
	"CrawlerProject/pkg/config"
	"CrawlerProject/pkg/logger"
	p "CrawlerProject/pkg/postgres"
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/exp/rand"

	cr "CrawlerProject/internal/crawler"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error while initializing config")
		os.Exit(3)
	}
	db := p.NewDBConnection(config.DBHost, config.DBPort, config.DBName, config.DBUser, config.DBPassword)
	localDB, err := db.InitConnection()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error while initalise database connection")
		os.Exit(3)
	}
	d := repository.NewDatabase(localDB)
	if err := d.Migrate(); err != nil {
		logger.Logger.Error().Err(err).Msg("error while migrating tables")
	}

	// repositories

  log.Println("Database tables created/migrated successfully!")
  service.SetDefaultDB(localDB)

  // // repositories
  // service.ReadFromJson(localDB)

	// telegram bot

	// err = bot.SetupBot(config.TGToken)
	// if err != nil {
	// 	return
	// }
	// crawler
	rand.Seed(uint64(time.Now().UnixNano()))

	// Create crawler with default config
	crawlerConfig := cr.DefaultConfig()
	crawler := cr.NewCrawler(crawlerConfig)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the crawler

	if err := crawler.Start(ctx); err != nil {
		log.Fatal(err)
	}

}
