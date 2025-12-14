package models

import (
	"net/mail"
	"time"
)

// User represents an application user (admin or customer)
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserInput struct {
	Username     string `json:"username" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Role         string `json:"role" binding:"required"`
}

func (input *UserInput) Validate() (bool, []string) {
	var invalidFields []string
	if input.Username == "" {
		invalidFields = append(invalidFields, "username")
	}
	if input.PasswordHash == "" {
		invalidFields = append(invalidFields, "password_hash")
	}
	if input.Email == "" || !IsValid_Email(input.Email) {
		invalidFields = append(invalidFields, "email")
	}
	if input.Role == "" {
		invalidFields = append(invalidFields, "role")
	}
	return len(invalidFields) == 0, invalidFields
}

// IsValidEmail checks if the given email has a valid format.
func IsValid_Email(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type UserLoginInput struct {
	Username     string `json:"username" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}
