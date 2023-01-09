package config

import (
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// get db connection from the DB_URL environment variable
	err := godotenv.Load("../.env")

	url := os.Getenv("DB_URL")

	fmt.Println(os.Getenv("DB_URL"))

	dbInstance, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = dbInstance
}
