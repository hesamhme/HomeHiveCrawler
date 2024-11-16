package repository

import (
	"CrawlerProject/internal/model"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{
		db,
	}
}

func (d *Database) Migrate() error {
	if err := d.AutoMigrate(&model.AdminLog{}, &model.CrawlerLog{}, &model.Filter{}, &model.Listing{}, &model.SearchHistory{}, &model.User{}); err != nil {
		return err
	}
	return nil
}
