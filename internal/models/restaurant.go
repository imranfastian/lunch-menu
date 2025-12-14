package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/mail"
	"time"
)

// FloatArray represents a PostgreSQL array of floats for coordinates
type FloatArray []float64

// Value implements the driver.Valuer interface
func (fa FloatArray) Value() (driver.Value, error) {
	return json.Marshal(fa)
}

// Scan implements the sql.Scanner interface
func (fa *FloatArray) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, fa)
	case string:
		return json.Unmarshal([]byte(v), fa)
	default:
		return errors.New("cannot scan into FloatArray")
	}
}

// Restaurant represents a restaurant
// DB schema and GORM models
type Restaurant struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	Address     string     `json:"address"`
	Coordinate  FloatArray `gorm:"type:jsonb" json:"coordinate"` // Store as JSONB in Postgres
	Homepage    string     `json:"homepage"`
	Region      string     `json:"region"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName overrides the table name used by GORM (optional)
func (Restaurant) TableName() string {
	return "restaurants"
}

// RestaurantInput represents the input for creating or updating a restaurant
// DTOs	API input/output structs
type RestaurantInput struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Address     string    `json:"address" binding:"required"`
	Coordinate  []float64 `json:"coordinate" binding:"required,len=2"` // [lat, lng]
	Homepage    string    `json:"homepage"`
	Region      string    `json:"region"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email" binding:"required"`
}

// Validate checks the RestaurantInput for required fields and correct formats.
func (input *RestaurantInput) Validate() (bool, []string) {
	var invalidFields []string
	if input.Name == "" {
		invalidFields = append(invalidFields, "name")
	}
	if input.Address == "" {
		invalidFields = append(invalidFields, "address")
	}
	if input.Email == "" || !IsValidEmail(input.Email) {
		invalidFields = append(invalidFields, "email")
	}
	if len(input.Coordinate) != 2 {
		invalidFields = append(invalidFields, "coordinate")
	}
	return len(invalidFields) == 0, invalidFields
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// RestaurantUpdateInput represents the input for updating a restaurant
type RestaurantUpdateInput struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Address     *string    `json:"address,omitempty"`
	Coordinate  *[]float64 `json:"coordinate,omitempty"`
	Homepage    *string    `json:"homepage,omitempty"`
	Region      *string    `json:"region,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Email       *string    `json:"email,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}
