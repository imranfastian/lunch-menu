package models

import "time"

// MenuItem represents a menu item
type MenuItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	Price        float64   `gorm:"not null" json:"price"`
	Category     string    `json:"category"`
	IsAvailable  bool      `gorm:"default:true" json:"is_available"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName overrides the table name used by GORM (optional)
func (MenuItem) TableName() string {
	return "menu_items"
}

// MenuItemInput represents the input for creating or updating a menu item
type MenuItemInput struct {
	RestaurantID uint    `json:"restaurant_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	Category     string  `json:"category"`
	IsAvailable  *bool   `json:"is_available"` // pointer to allow "not set"
}

type MenuItemUpdateInput struct {
	RestaurantID *uint    `json:"restaurant_id,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	Category     *string  `json:"category,omitempty"`
	IsAvailable  *bool    `json:"is_available,omitempty"`
}
