package controllers

import	"github.com/Hdeee1/go-ecommerce/repository"

type Application struct {
	UserRepo	*repository.UserRepository
	ProductRepo	*repository.ProductRepository
	OrderRepo	*repository.OrderRepository
	AddressRepo	*repository.AddressRepository
}

func NewApplication(
	userRepo *repository.UserRepository,
	productRepo *repository.ProductRepository,
	orderRepo *repository.OrderRepository,
	addressRepo *repository.AddressRepository,
) *Application {
	return &Application{
		UserRepo: userRepo,
		ProductRepo: productRepo,
		OrderRepo: orderRepo,
		AddressRepo: addressRepo,
	}
}