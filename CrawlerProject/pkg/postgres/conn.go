package postgres

import (
	"CrawlerProject/pkg/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnection struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func NewDBConnection(host string, port string, dbName string, user string, password string) *DBConnection {
	return &DBConnection{
		// "host=localhost user=user password=QPdb7e3m dbname=crawler_db port=5432 sslmode=disable"
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}
}

func (d *DBConnection) InitConnection() (*gorm.DB, error) {
	// func InitConnection(conf *config.config) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", d.Host, d.User, d.Password, d.DBName, d.Port),
		PreferSimpleProtocol: true,
	}), &gorm.Config{PrepareStmt: false})

	if err != nil {
		logger.Logger.Error().Err(err).Msg("error while openning database connection")
		return nil, err
	}
	return db, nil
}
