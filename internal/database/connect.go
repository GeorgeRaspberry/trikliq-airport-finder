package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		database.User,
		database.Password,
		database.Host,
		database.Port,
		database.Database,
	)

	// Open connection to the database.
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Run migrations.
	err = Migrate(gormDB)
	if err != nil {
		return err
	}

	DB = gormDB
	return nil
}
