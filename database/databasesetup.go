package controllers

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := "root:dbmysql@tcp(127.0.0.1:3306)/db_ecommerce?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	fmt.Println("Connection Opened to Database")
}