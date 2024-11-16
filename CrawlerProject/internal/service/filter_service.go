package service

import (
	"CrawlerProject/internal/model"

	"gorm.io/gorm"
)

func GetFilteredListings(db *gorm.DB, filters model.Filter) ([]model.Listing, error) {
	var listings []model.Listing
	query := db.Model(&model.Listing{})

	// filter by city
	if filters.City != "" {
        query = query.Where("city LIKE ?", "%"+filters.City+"%")
    }

	// filter by nighbor
	if filters.Neighborhood != "" {
		query = query.Where("neighborhood = ?", filters.Neighborhood)
	}

	// price range filter
	if filters.PriceMin > 0 {
		query = query.Where("price >= ?", filters.PriceMin)
	}
	if filters.PriceMax > 0 {
		query = query.Where("price <= ?", filters.PriceMax)
	}

	// area range filter
	if filters.AreaMin > 0 {
		query = query.Where("meterage  >= ?", filters.AreaMin)
	}
	if filters.AreaMax > 0 {
		query = query.Where("meterage  <= ?", filters.AreaMax)
	}

	// rooms filter
	if filters.RoomsMin > 0 {
		query = query.Where("bedrooms >= ?", filters.RoomsMin)
	}
	if filters.RoomsMax > 0 {
		query = query.Where("bedrooms <= ?", filters.RoomsMax)
	}

	// status filter
	if filters.Status != "" {
		if filters.Status == "اجاره و رهن" {
			// Handle "اجاره و رهن" as a special case to return both "اجاره" and "رهن"
			query = query.Where("ad_type IN (?)", []string{"اجاره", "رهن"})
		} else {
			query = query.Where("ad_type = ?", filters.Status)
		}
	}

	// building age filter
	if filters.BuildingAgeMin > 0 {
		query = query.Where("age >= ?", filters.BuildingAgeMin)
	}
	if filters.BuildingAgeMax > 0 {
		query = query.Where("age <= ?", filters.BuildingAgeMax)
	}

	// filter by type
	if filters.PropertyType != "" {
		query = query.Where("houseType = ?", filters.PropertyType)
	}

	// floor filter
	if filters.FloorMin != nil {
		query = query.Where("floor >= ?", *filters.FloorMin)
	}
	if filters.FloorMax != nil {
		query = query.Where("floor <= ?", *filters.FloorMax)
	}

	// elevator and storage
	if filters.HasStorage {
		query = query.Where("warehouse = ?", filters.HasStorage)
	}
	if filters.HasElevator {
		query = query.Where("elevator = ?", filters.HasElevator)
	}

	// filter base on ad date
	if !filters.CreationDateMin.IsZero() {
		query = query.Where("ad_create_date >= ?", filters.CreationDateMin)
	}
	if !filters.CreationDateMax.IsZero() {
		query = query.Where("ad_dreate_date <= ?", filters.CreationDateMax)
	}

	// filter base radius
	if filters.Latitude != 0 && filters.Longitude != 0 && filters.Radius > 0 {
		// eg
		query = query.Where(`
            ST_Distance_Sphere(
                point(longitude, latitude),
                point(?, ?)
            ) <= ?`, filters.Longitude, filters.Latitude, filters.Radius)
	}

	// run and return results
	result := query.Find(&listings)
	return listings, result.Error
}
