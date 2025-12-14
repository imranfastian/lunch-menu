package database

import (
	"lunch_menu/internal/models"
	"time"
)

// Add token to blacklist
func BlacklistToken(token string, expiresAt time.Time) error {
	return DB.Create(&models.BlacklistedToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}).Error
}

// Check if token is blacklisted
func IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := DB.Model(&models.BlacklistedToken{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error
	return count > 0, err
}
