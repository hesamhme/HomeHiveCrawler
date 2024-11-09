package service

import (
	"CrawlerProject/internal/model"
	"gorm.io/gorm"
)

func SaveOrUpdateListing(db *gorm.DB, listing *model.Listing) error {
	// check for listing
	var existingListing model.Listing
	result := db.Where("url = ?", listing.URL).First(&existingListing)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	// update existing or create new
	if existingListing.ListingID != 0 {
		listing.ListingID = existingListing.ListingID
		return db.Save(listing).Error
	}
	return db.Create(listing).Error
}
