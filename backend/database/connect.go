package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ConnectionStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Tokyo",
		"postgres",
		"5432",
		"postgres",
		"postgres",
		"postgrespassword",
	)
	DB *gorm.DB
)

func Connect() {
	conn, err := gorm.Open(postgres.Open(ConnectionStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = conn
	conn.AutoMigrate()
}
