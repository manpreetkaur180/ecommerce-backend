package database

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "fmt"
)

func Connect() *gorm.DB {
  dsn := "host=localhost user=manpreet password=2004 dbname=ecommerce port=5432 sslmode=disable"
  

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    
    if err != nil {
        log.Fatal("DB connection failed")
    }
    fmt.Println("DATABASE CONNECTED SUCCESSFULLY")

    return db
}