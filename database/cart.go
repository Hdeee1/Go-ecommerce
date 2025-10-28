package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrCantFindProduct = errors.New("Can't find the product")
	ErrCantDecodeProduct = errors.New("Cant find the product")
	ErrUserIdNotValid = errors.New("This user is not valid")
	ErrCantUpdateUser = errors.New("Cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("cannot update the purchase")
)

func AddToCart() gin.HandlerFunc {
	
}

func RemoveItem() gin.HandlerFunc {
	
}

func GetItemFromCart() gin.HandlerFunc {
	
}

func BuyFromCart() gin.HandlerFunc {
	
}

func InstantBuy() gin.HandlerFunc {
	
}