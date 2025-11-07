package controllers

import	"github.com/Hdeee1/go-ecommerce/repository"

type Application struct {
	UserRepo	*repository.UserRepository
	ProductRepo	*repository.ProductRepository
	OrderRepo	*repository.OrderRepository
}

func NewApplication(
	userRepo *repository.UserRepository,
	productRepo *repository.ProductRepository,
	orderRepo *repository.OrderRepository,
) *Application {
	return &Application{
		UserRepo: userRepo,
		ProductRepo: productRepo,
		OrderRepo: orderRepo,
	}
}