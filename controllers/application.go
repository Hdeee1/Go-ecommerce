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

func (app *Application) AddToCart() func() {
	return func() {
		
	}
	
}

func (app *Application) RemoveItem() func() {
	return func() {
		
	}

	
}

func (app *Application) BuyFromCart() func() {
	return func() {
		
	}

	
}

func (app *Application) InstantBuy() func() {
	return func() {
		
	}

	
}