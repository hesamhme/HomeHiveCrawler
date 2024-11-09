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
        query = query.Where("city = ?", filters.City)
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
        query = query.Where("area >= ?", filters.AreaMin)
    }
    if filters.AreaMax > 0 {
        query = query.Where("area <= ?", filters.AreaMax)
    }

    // rooms filter
    if filters.RoomsMin > 0 {
        query = query.Where("rooms >= ?", filters.RoomsMin)
    }
    if filters.RoomsMax > 0 {
        query = query.Where("rooms <= ?", filters.RoomsMax)
    }

    // status filter
    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }

    // building age filter
    if filters.BuildingAgeMin > 0 {
        query = query.Where("building_age >= ?", filters.BuildingAgeMin)
    }
    if filters.BuildingAgeMax > 0 {
        query = query.Where("building_age <= ?", filters.BuildingAgeMax)
    }

    // filter by type
    if filters.PropertyType != "" {
        query = query.Where("property_type = ?", filters.PropertyType)
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
        query = query.Where("has_storage = ?", filters.HasStorage)
    }
    if filters.HasElevator {
        query = query.Where("has_elevator = ?", filters.HasElevator)
    }

    // filter base on ad date
    if !filters.CreationDateMin.IsZero() {
        query = query.Where("created_at >= ?", filters.CreationDateMin)
    }
    if !filters.CreationDateMax.IsZero() {
        query = query.Where("created_at <= ?", filters.CreationDateMax)
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
