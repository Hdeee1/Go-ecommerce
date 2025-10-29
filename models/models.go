package models

import (
	// "database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_Name			*string			`json:"first_name" validate:"required,min=2,max=30"`
	Last_Name			*string			`json:"last_name" validate:"required,min=2,max=30"`
	Password			*string			`json:"password" validate:"required,min=6"`
	Email				*string			`json:"email" validate:"required,email" gorm:"unique"`
	Phone				*string			`json:"phone" validate:"required" gorm:"unique"`
	Token				*string			`json:"token"`
	Refresh_Token		*string			`json:"refresh_token"`
	User_ID				string			`json:"user_id" gorm:"unique;index"`

	Address_Details		[]Address		`json:"address" gorm:"foreignKey:UserID"`
	Order_Status		[]Order			`json:"orders" gorm:"foreignKey:UserID"`
}

// Generate UserID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.User_ID == "" {
		u.User_ID = uuid.New().String()
	}

	return nil
}

type Product struct {
	gorm.Model
	Product_Name	*string		`json:"product_name" validate:"required"`
	Price			*uint64		`json:"price" validate:"required"`
	Rating			*uint8		`json:"rating"`
	Image			*string		`json:"image"`
}

type Address struct {
	gorm.Model
	UserID		uint		`json:"user_id"`
	House		*string		`json:"house" validate:"required"`
	Street		*string		`json:"street" validate:"required"`
	City		*string		`json:"city" validate:"required"`
	Pincode		*string		`json:"pincode" validate:"required"`
}

type OrderItem struct {
    gorm.Model
    OrderID             uint            `json:"order_id"`
    ProductID           uint            `json:"product_id"`
    Product_Name        *string         `json:"product_name"`
    Price               uint64          `json:"price"`
    Rating              *uint8          `json:"rating"`
    Image               *string         `json:"image"`
    Quantity            uint            `json:"quantity"validate:"required,min=1"`
}

type Order struct {
	gorm.Model
	UserID			uint			`json:"user_id"`
	OrderItems		[]OrderItem		`json:"order_item" gorm:"foreignKey:OrderID"`
	Order_At		time.Time		`json:"order_at"`
	Price			int				`json:"price"`
	Discount		*int			`json:"discount"`
	Payment_Method	Payment			`json:"payment_method" gorm:"embedded"`
}

type Payment struct {
	Digital		bool	`json:"digital"`
	COD			bool	`json:"cod"`
}