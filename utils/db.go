package utils

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDBConnection() *gorm.DB {

	var db *gorm.DB
	var err error

	_, ok := os.LookupEnv("MACMAHOME_TEST")
	if !ok {
		dsn := "host=db user=admin password=admin dbname=macmahome port=5432"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open("../data/sqlite/test.db"), &gorm.Config{})
	}

	if err != nil {
		panic("Failed to connect database")
	}

	return db
}
