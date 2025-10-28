package models

import (
	// "database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_Name			*string			`json:"first_name"`
	Last_Name			*string			`json:"last_name"`
	Password			*string			`json:"password"`
	Email				*string			`json:"email" gorm:"unique"`
	Phone				*string			`json:"phone" gorm:"unique"`
	Token				*string			`json:"token"`
	Refresh_Token		*string			`json:"refresh_token"`
	User_ID				string			`json:"user_id" gorm:"unique;index"`

	Address_Details		[]Address		`json:""`
	Order_Status		[]Order			`json:""`
}

type Product struct {
	gorm.Model
	Product_Name	*string		`json:"product_name"`
	Price			*uint64		`json:"price"`
	Rating			*uint8		`json:"rating"`
	Image			*string		`json:"image"`
}

type Address struct {
	gorm.Model
	UserID		uint		`json:"user_id"`
	House		*string		`json:"house"`
	Street		*string		`json:"street"`
	City		*string		`json:"city"`
	Pincode		*string		`json:"pincode"`
}

type OrderItem struct {
    gorm.Model
    OrderID             uint            `json:"order_id"`
    ProductID           uint            `json:"product_id"`
    Product_Name        *string         `json:"product_name"`
    Price               uint64          `json:"price"`
    Rating              *uint8          `json:"rating"`
    Image               *string         `json:"image"`
    Quantity            uint            `json:"quantity"`
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