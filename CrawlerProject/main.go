package main


import (

	"CrawlerProject/internal/Model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=user password=QPdb7e3m dbname=crawler_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate schema
	db.AutoMigrate(
		&Model.User{},
		&Model.Listing{},
		&Model.Filter{},
		&Model.Bookmark{},
		&Model.Notification{},
		&Model.SearchHistory{},
		&Model.AdminLog{},
		&Model.CrawlerLog{},
	)

	log.Println("Database tables created/migrated successfully!")
}
