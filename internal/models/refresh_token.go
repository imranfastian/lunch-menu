package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token for a user session
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	TokenHash string     `gorm:"not null" json:"token_hash"`
	UserAgent string     `gorm:"not null" json:"user_agent"`
	IPAddress string     `gorm:"not null" json:"ip_address"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"` // pointer for nullable
}

// TableName overrides the table name used by GORM (optional)
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
