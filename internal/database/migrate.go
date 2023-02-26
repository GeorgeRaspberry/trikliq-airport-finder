package database

import (
	"bookbox-backend/internal/model"

	"gorm.io/gorm"
)

func Migrate(gormDB *gorm.DB) error {

	// Migrate ORM models.
	err := gormDB.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		return err
	}

	return nil
}
