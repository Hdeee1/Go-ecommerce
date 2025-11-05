package controllers

import "gorm.io/gorm"

type Application struct {
	ProductDB 	*gorm.DB
	UserDB		*gorm.DB
}

func NewApplication(productDB, userDB *gorm.DB) *Application {
	return &Application{
		ProductDB: productDB,
		UserDB: userDB,
	}
}