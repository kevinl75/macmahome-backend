package utils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConnection() *gorm.DB {

	dsn := "host=db user=admin password=admin dbname=macmahome port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	return db
}
