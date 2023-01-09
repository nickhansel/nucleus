package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// get db connection from the DB_URL environment variable
	dbURL := "postgresql://postgres:december1963@database-2.c08a8epqacce.us-east-1.rds.amazonaws.com"

	dbInstance, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// generate struct for the db table using the gorm gen package
	// generateTable(db)

	DB = dbInstance
}
