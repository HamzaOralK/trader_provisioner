package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type DB struct {
	instance *gorm.DB
}

func databaseInitialize() DB {
	// dsn := "host=postgres-postgresql user=postgres password=12345aaa dbname=trader port=5432 sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		"trader",
	)

	instance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("DB Connection is successful")
	}

	return DB {
		instance: instance,
	}

}
