package database

import (
	"errors"
	"lunch_menu/internal/models"

	"gorm.io/gorm"
)

// CreateMenuItem inserts a new menu item into the database
func CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	item.IsAvailable = true
	if err := DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// UpdateMenuItem updates an existing menu item by ID
func UpdateMenuItem(id uint, item *models.MenuItem) (*models.MenuItem, error) {
	menuItem, err := findMenuItemByID(id)
	if err != nil {
		return nil, err
	}
	menuItem.RestaurantID = item.RestaurantID
	menuItem.Name = item.Name
	menuItem.Description = item.Description
	menuItem.Price = item.Price
	menuItem.Category = item.Category
	menuItem.IsAvailable = item.IsAvailable

	if err := DB.Save(menuItem).Error; err != nil {
		return nil, err
	}
	return menuItem, nil
}

// DeleteMenuItem soft deletes a menu item by setting is_available to false
func DeleteMenuItem(id uint) error {
	_, err := findMenuItemByID(id)
	if err != nil {
		return err
	}
	result := DB.Model(&models.MenuItem{}).Where("id = ?", id).Update("is_available", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("menu item not found")
	}
	return nil
}

// GetMenuItems retrieves menu items for a restaurant with pagination
func GetMenuItems(restaurantID uint, limit, offset int) ([]models.MenuItem, int64, error) {
	var items []models.MenuItem
	var total int64

	// Query for menu items
	err := DB.Model(&models.MenuItem{}).
		Where("restaurant_id = ?", restaurantID).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

// GetMenuItemByID retrieves a menu item by its ID
func GetMenuItemByID(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := DB.Where("id = ? AND is_available = ?", id, true).First(&menuItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu item not found")
		}
		return nil, err
	}
	return &menuItem, nil
}

func findMenuItemByID(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := DB.First(&menuItem, id).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

// PartialUpdateMenuItem updates specific fields of an existing menu item by ID
func PartialUpdateMenuItem(id uint, input *models.MenuItemUpdateInput) (*models.MenuItem, error) {
	menuItem, err := findMenuItemByID(id)
	if err != nil {
		return nil, err
	}
	if input.RestaurantID != nil {
		if *input.RestaurantID == 0 {
			return nil, errors.New("invalid restaurant_id")
		}
		// Optionally, check if restaurant exists
		var restaurant models.Restaurant
		if err := DB.First(&restaurant, *input.RestaurantID).Error; err != nil {
			return nil, errors.New("restaurant_id does not exist")
		}
		menuItem.RestaurantID = *input.RestaurantID
	}
	if input.Name != nil {
		menuItem.Name = *input.Name
	}
	if input.Description != nil {
		menuItem.Description = *input.Description
	}
	if input.Price != nil {
		menuItem.Price = *input.Price
	}
	if input.Category != nil {
		menuItem.Category = *input.Category
	}
	if input.IsAvailable != nil {
		menuItem.IsAvailable = *input.IsAvailable
	}
	if err := DB.Save(menuItem).Error; err != nil {
		return nil, err
	}
	return menuItem, nil
}
