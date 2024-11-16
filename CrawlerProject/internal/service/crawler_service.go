package service

import (
	"CrawlerProject/internal/model"
	"gorm.io/gorm"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func SaveOrUpdateListing(db *gorm.DB, listing *model.Listing) error {
	// check for listing
	var existingListing model.Listing
	result := db.Where("url = ?", listing.Link).First(&existingListing)
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


func ReadFromJson(db *gorm.DB) {
    file, err := os.Open("C:/Users/Hesam/Desktop/project/numbr1/quera_gr11_project1/quera_g11_project1/CrawlerProject/internal/service/sample_data.json") 
    if err != nil {
        log.Fatalf("failed to open JSON file: %v", err)
    }
    defer file.Close()

    byteValue, _ := ioutil.ReadAll(file)

    // تبدیل JSON به ساختار داده Go
    var listings []model.Listing
    err = json.Unmarshal(byteValue, &listings)
    if err != nil {
        log.Fatalf("failed to unmarshal JSON: %v", err)
    }

    // ذخیره داده‌ها در دیتابیس
    for _, listing := range listings {
        err = db.Debug().Create(&listing).Error
        if err != nil {
            log.Printf("failed to save listing: %v", err)
        }
    }

    fmt.Println("Data inserted successfully!")
}
