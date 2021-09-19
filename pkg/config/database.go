package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	instance *gorm.DB
}

func DatabaseInitialize() DB {
	sslMode := "require"
	environment := os.Getenv("ENVIRONMENT")
	if environment == "local" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SCHEMA"),
		os.Getenv("DB_PORT"),
		sslMode,
	)

	instance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("DB Connection is successful")
	}

	return DB{
		instance: instance,
	}

}
