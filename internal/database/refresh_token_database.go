package database

import (
	"crypto/sha256"
	"encoding/hex"
	"lunch_menu/internal/models"
	"time"
)

// SaveRefreshToken saves a hashed refresh token in the database.
func SaveRefreshToken(userID uint, token, userAgent, ip string, expiresAt time.Time) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	if userAgent == "" {
		userAgent = "unknown"
	}
	if ip == "" {
		ip = "unknown"
	}

	// Upsert using GORM
	var rt models.RefreshToken
	result := DB.Where("user_id = ? AND user_agent = ? AND ip_address = ?", userID, userAgent, ip).First(&rt)
	if result.Error == nil {
		// Update existing
		rt.TokenHash = tokenHash
		rt.ExpiresAt = expiresAt
		rt.RevokedAt = nil
		return DB.Save(&rt).Error
	} else if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}

	// Create new
	rt = models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		UserAgent: userAgent,
		IPAddress: ip,
		ExpiresAt: expiresAt,
	}
	return DB.Create(&rt).Error
}

// GetRefreshToken finds a refresh token by its raw value and user ID.
func GetRefreshToken(userID int, token string) (*models.RefreshToken, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	var rt models.RefreshToken
	err := DB.Where(
		"user_id = ? AND token_hash = ? AND revoked_at IS NULL AND expires_at > ?", userID, tokenHash, time.Now(),
	).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// DeleteRefreshToken deletes a refresh token for a user
func DeleteRefreshToken(userID uint, token string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	return DB.Where("user_id = ? AND token_hash = ?", userID, tokenHash).Delete(&models.RefreshToken{}).Error
}
