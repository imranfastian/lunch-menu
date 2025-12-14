package database

import (
	"errors"
	"lunch_menu/internal/models"

	"gorm.io/gorm"
)

// CreateRestaurant inserts a new restaurant into the database
func CreateRestaurant(r *models.Restaurant) (*models.Restaurant, error) {
	r.IsActive = true
	if err := DB.Create(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateRestaurant updates an existing restaurant by ID
func UpdateRestaurant(id uint, r *models.Restaurant) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := DB.First(&restaurant, id).Error; err != nil {
		return nil, err
	}
	// Update fields
	restaurant.Name = r.Name
	restaurant.Description = r.Description
	restaurant.Address = r.Address
	restaurant.Coordinate = r.Coordinate
	restaurant.Homepage = r.Homepage
	restaurant.Region = r.Region
	restaurant.Phone = r.Phone
	restaurant.Email = r.Email
	restaurant.IsActive = r.IsActive

	if err := DB.Save(&restaurant).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// DeleteRestaurant soft deletes a restaurant by setting is_active to false
func DeleteRestaurant(id uint) error {
	return DB.Model(&models.Restaurant{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// GetRestaurants retrieves all active restaurants with pagination
func GetRestaurants(limit, offset int) ([]models.Restaurant, int64, error) {
	var restaurants []models.Restaurant
	var total int64

	if err := DB.Model(&models.Restaurant{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := DB.Where("is_active = ?", true).Order("id").Limit(limit).Offset(offset).Find(&restaurants).Error; err != nil {
		return nil, 0, err
	}

	return restaurants, total, nil
}

// GetRestaurantByID retrieves a single active restaurant by ID
func GetRestaurantByID(id uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := DB.Where("id = ? AND is_active = ?", id, true).First(&restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("restaurant not found")
		}
		return nil, err
	}
	return &restaurant, nil
}

// PartialUpdateRestaurant updates only the non-zero fields of an existing restaurant
func PartialUpdateRestaurant(id uint, input *models.RestaurantUpdateInput) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := DB.First(&restaurant, id).Error; err != nil {
		return nil, err
	}
	// Only update fields that are not nil
	if input.Name != nil {
		restaurant.Name = *input.Name
	}
	if input.Description != nil {
		restaurant.Description = *input.Description
	}
	if input.Address != nil {
		restaurant.Address = *input.Address
	}
	if input.Coordinate != nil {
		restaurant.Coordinate = *input.Coordinate
	}
	if input.Homepage != nil {
		restaurant.Homepage = *input.Homepage
	}
	if input.Region != nil {
		restaurant.Region = *input.Region
	}
	if input.Phone != nil {
		restaurant.Phone = *input.Phone
	}
	if input.Email != nil {
		restaurant.Email = *input.Email
	}
	if input.IsActive != nil {
		restaurant.IsActive = *input.IsActive
	}
	if err := DB.Save(&restaurant).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}
