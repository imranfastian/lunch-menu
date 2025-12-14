package models

import (
	"time"
)

// AuditLog represents an audit log entry for admin actions
type AuditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ActorUserID uint      `gorm:"index;not null" json:"actor_user_id"`
	Action      string    `gorm:"not null" json:"action"`
	TargetTable string    `gorm:"not null" json:"target_table"`
	TargetID    string    `gorm:"not null" json:"target_id"`
	Before      string    `json:"before"` // JSON string
	After       string    `json:"after"`  // JSON string
	IPAddress   string    `json:"ip_address"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName overrides the table name used by GORM (optional)
func (AuditLog) TableName() string {
	return "audit_logs"
}
