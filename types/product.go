package types

import (
	"time"
)

type ApproveProduct struct {
	ProductID  string `json:"product_id"`
	IsVerified bool   `json:"is_verified"`
}

type Product struct {
	ID                 int        `json:"id" db:"id"`
	Img                string     `json:"img" db:"img"`
	FarmerID           int        `json:"farmer_id" db:"farmer_id"`
	Name               string     `json:"name" db:"name"`
	Type               string     `json:"type" db:"type"`
	Quantity           int        `json:"quantity_in_kg" db:"quantity_in_kg"`
	RatePerKg          float64    `json:"rate_per_kg" db:"rate_per_kg"`
	JariSize           string     `json:"jari_size,omitempty" db:"jari_size"`
	ExpectedDelivery   *time.Time `json:"expected_delivery,omitempty" db:"expected_delivery"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	FarmersPhoneNumber string     `json:"farmer_phone_number" db:"farmers_phone_number"`
	FarmerFirstName    string     `json:"farmers_first_name,omitempty"`
	FarmerLastName     string     `json:"farmers_last_name,omitempty"`
	IsAvailable        bool       `json:"is_available" db:"is_available"`
	IsVerifiedByAdmin  bool       `json:"is_verified_by_admin" db:"is_verified_by_admin"`
}


type OrderSummary struct {
	OrderID              int        `json:"order_id"`
	QuantityInKg         int        `json:"quantity_in_kg"`
	TotalPrice           float64    `json:"total_price"`
	Status               string     `json:"status"`
	ModeOfDelivery       string     `json:"mode_of_delivery"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
	OrderDate            time.Time  `json:"order_date"`
	ProductID            int        `json:"product_id"`
	ProductName          string     `json:"product_name"`
	UserID               int        `json:"user_id"`
	UserFirstName        string     `json:"user_first_name"`
	UserLastName         string     `json:"user_last_name"`
	UserPhoneNumber      string     `json:"user_phone_number"`
	DeliveryAddress      string     `json:"delivery_address"`
	DeliveryCity         string     `json:"delivery_city"`
	DeliveryAddressZIP   int        `json:"delivery_address_pin_code"`
	BuyersPhoneNumber    int        `json:"buyers_phone_number" db:"buyers_phone_number"`
	FarmersPhoneNumber   int        `json:"farmer_phone_number" db:"farmers_phone_number"`
	ProductImg           string     `json:"product_img" db:"img"`
}

// type OrderStatus struct {
// 	OrderDetails 
// 	BuyersDetails
// 	SEllerDetails
// }