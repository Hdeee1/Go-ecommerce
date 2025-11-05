package database

import (
	"fmt"
	"log"

	"github.com/Hdeee1/go-ecommerce/config"
	"github.com/Hdeee1/go-ecommerce/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	fmt.Println("Connection Opened to Database")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	fmt.Println("Database migration complete")
}

func UserData() *gorm.DB {
	return DB
}

func ProductData() *gorm.DB {
	return DB
}