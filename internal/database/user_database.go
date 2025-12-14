package database

import (
	"errors"
	"lunch_menu/internal/models"

	"gorm.io/gorm"
)

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) (*models.User, error) {
	if err := DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByID retrieves a user by their ID
func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := DB.First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// UpdateUser updates an existing user in the database
func UpdateUser(user *models.User) error {
	return DB.Save(user).Error
}

// DeleteUser deletes a user from the database
func DeleteUser(userID uint) error {
	return DB.Delete(&models.User{}, userID).Error
}
