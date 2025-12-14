package database

import (
	"fmt"
	"log"
	"lunch_menu/internal/config"
	"lunch_menu/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	host := config.AppConfig.DBHost
	port := config.AppConfig.DBPort
	user := config.AppConfig.DBUser
	password := config.AppConfig.DBPassword
	dbname := config.AppConfig.DBName
	sslmode := config.AppConfig.DBSSLmode
	timezone := config.AppConfig.DBTimeZone

	// fmt.Printf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s\n", host, port, user, password, dbname, sslmode, timezone)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	DB = db

	log.Println("Database connection established successfully")
	return nil
}

func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Auto-migrate all models
func Migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Restaurant{},
		&models.MenuItem{},
		&models.RefreshToken{},
		&models.BlacklistedToken{},
		&models.AuditLog{},
	)
}

// GetBusinessStatistics retrieves business analytics data
func GetBusinessStatistics() (*models.BusinessStatistics, error) {
	stats := &models.BusinessStatistics{}

	// Get restaurant counts
	var counts struct {
		Total    int64
		Active   int64
		Inactive int64
	}

	err := DB.Table("restaurants").Select(`
		COUNT(*) as total,
		COUNT(CASE WHEN is_active = true THEN 1 END) as active,
		COUNT(CASE WHEN is_active = false THEN 1 END) as inactive
	`).Scan(&counts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant stats: %w", err)
	}

	stats.TotalRestaurants = counts.Total
	stats.ActiveRestaurants = counts.Active
	stats.InactiveRestaurants = counts.Inactive

	// Get menu item count and average price
	var menuStats struct {
		TotalMenuItems int64
		AveragePrice   float64
	}

	err = DB.Table("menu_items").Select(`
		COUNT(*) as total_menu_items,
		COALESCE(AVG(price), 0) as average_price
	`).Where("is_available = true").Scan(&menuStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get menu stats: %w", err)
	}

	stats.TotalMenuItems = menuStats.TotalMenuItems
	stats.AveragePrice = menuStats.AveragePrice

	// Get revenue by category
	rows, err := DB.Table("menu_items").Select(`
		category, SUM(price * 100) as estimated_revenue 
	`).Where("is_available = true").Group("category").Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue by category: %w", err)
	}
	defer rows.Close()

	stats.RevenueByCateory = make(map[string]float64)
	for rows.Next() {
		var category string
		var revenue float64
		if err := rows.Scan(&category, &revenue); err != nil {
			return nil, fmt.Errorf("failed to scan revenue data: %w", err)
		}
		stats.RevenueByCateory[category] = revenue
	}

	// Get detailed per-restaurant business data
	restaurantRows, err := DB.Table("restaurants r").
		Select(`
			r.id, r.name,
			COUNT(m.id) as menu_count,
			COALESCE(AVG(m.price), 0) as avg_price,
			COALESCE(SUM(m.price * 150), 0) as estimated_total_revenue
		`).
		Joins("LEFT JOIN menu_items m ON r.id = m.restaurant_id AND m.is_available = true").
		Where("r.is_active = true").
		Group("r.id, r.name").
		Order("estimated_total_revenue DESC").
		Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant business data: %w", err)
	}
	defer restaurantRows.Close()

	var restaurantDetails []models.RestaurantBusinessData
	for restaurantRows.Next() {
		var detail models.RestaurantBusinessData
		if err := restaurantRows.Scan(&detail.RestaurantID, &detail.RestaurantName,
			&detail.MenuItemCount, &detail.AveragePrice, &detail.TotalRevenue); err != nil {
			return nil, fmt.Errorf("failed to scan restaurant business data: %w", err)
		}
		restaurantDetails = append(restaurantDetails, detail)
	}
	stats.RestaurantDetails = restaurantDetails

	return stats, nil
}
