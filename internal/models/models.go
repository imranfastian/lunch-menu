package models

import "time"

// Response Models

// RestaurantsResponse represents the response for listing restaurants
type RestaurantsResponse struct {
	Restaurants []Restaurant `json:"restaurants"`
	Total       int64        `json:"total"`
}

// RestaurantResponse represents the response for a single restaurant
type RestaurantResponse struct {
	Restaurant Restaurant `json:"restaurant"`
}

// MenuItemsResponse represents the response for restaurant menu items
type MenuItemsResponse struct {
	MenuItems []MenuItem `json:"menu_items"`
	Total     int64      `json:"total"`
}

// APIResponse represents basic API information
type APIResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// BusinessStatistics represents business analytics data
type BusinessStatistics struct {
	TotalRestaurants    int64                    `json:"total_restaurants"`
	ActiveRestaurants   int64                    `json:"active_restaurants"`
	InactiveRestaurants int64                    `json:"inactive_restaurants"`
	TotalMenuItems      int64                    `json:"total_menu_items"`
	AveragePrice        float64                  `json:"average_price"`
	RevenueByCateory    map[string]float64       `json:"revenue_by_category"`
	RestaurantDetails   []RestaurantBusinessData `json:"restaurant_details"`
}

// RestaurantBusinessData represents per-restaurant business data
type RestaurantBusinessData struct {
	RestaurantID   uint    `json:"restaurant_id"`
	RestaurantName string  `json:"restaurant_name"`
	MenuItemCount  int64   `json:"menu_item_count"`
	AveragePrice   float64 `json:"average_price"`
	TotalRevenue   float64 `json:"total_revenue"`
}

// SafeUser is used for API responses to hide sensitive fields
type SafeUser struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToSafeUser converts a User to a SafeUser
func (u *User) ToSafeUser() SafeUser {
	return SafeUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// StandardResponse represents a standard response format
type StandardResponse struct {
	Message string         `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}
