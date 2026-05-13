package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {

	var db *gorm.DB
	var err error

	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// retry logic
	for i := 1; i <= 5; i++ {

		db, err = gorm.Open(
			postgres.Open(dsn),
			&gorm.Config{},
		)

		if err == nil {
			log.Println("Product database connected successfully")
			return db
		}

		log.Printf(
			"Waiting for product DB... retry %d",
			i,
		)

		time.Sleep(2 * time.Second)
	}

	log.Fatal("Failed to connect product DB:", err)

	return nil
}