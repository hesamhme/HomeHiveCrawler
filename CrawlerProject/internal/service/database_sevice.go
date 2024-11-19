package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gorm.io/gorm"

	"CrawlerProject/internal/model"
)

var defaultDB *gorm.DB

func SetDefaultDB(db *gorm.DB) {
	defaultDB = db
}

// GetListings fetches all listings from the database.
func GetListings(db *gorm.DB) ([]model.Listing, error) {
	if db == nil {
		db = defaultDB
	}
	var listings []model.Listing
	if err := db.Find(&listings).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch listings: %w", err)
	}
	return listings, nil
}

// StoreListing saves or updates a single listing in the database.
func StoreListing(db *gorm.DB, listing model.Listing) error {
	if db == nil {
		db = defaultDB
	}
	// Check for an existing listing with the same Link.
	var existingListing model.Listing
	err := db.Where("link = ?", listing.URL).First(&existingListing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to query existing listing: %w", err)
	}

	if existingListing.ListingID != 0 {
		// Update existing listing.
		listing.ListingID = existingListing.ListingID
		if err := db.Save(&listing).Error; err != nil {
			return fmt.Errorf("failed to update listing: %w", err)
		}
	} else {
		// Create a new listing.
		if err := db.Create(&listing).Error; err != nil {
			return fmt.Errorf("failed to create listing: %w", err)
		}
	}
	return nil
}

// StoreAllListings reads listings from a JSON file and saves/updates them in the database.
func StoreAllListings(db *gorm.DB, filePath string) error {
	if db == nil {
		db = defaultDB
	}
	// Open the JSON file.
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the JSON into a slice of Listing.
	var listings []model.Listing
	if err := json.NewDecoder(file).Decode(&listings); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Iterate through each listing and save/update it.
	for _, listing := range listings {
		if err := StoreListing(db, listing); err != nil {
			fmt.Printf("failed to store listing with link %s: %v\n", listing.URL, err)
		}
	}
	return nil
}
