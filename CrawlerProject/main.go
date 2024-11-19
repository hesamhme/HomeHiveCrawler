package main

import (
    "CrawlerProject/internal/bot"
    "CrawlerProject/internal/repository"
    //"CrawlerProject/internal/service"
    "CrawlerProject/pkg/config"
    "CrawlerProject/pkg/logger"
    p "CrawlerProject/pkg/postgres"
    "log"
    "os"
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

    // // repositories
    // service.ReadFromJson(localDB)

    log.Println("Database tables created/migrated successfully!")

    // bot
    bot.SetDB(localDB)

    err = bot.SetupBot(config.TGToken)
    if err != nil {
        return
    }

}




