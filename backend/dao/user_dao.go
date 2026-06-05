package dao

import (
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *domain.User) error {
	return db.Create(user).Error
}

func FindUserByEmail(db *gorm.DB, email string) (*domain.User, error) {
	var user domain.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
