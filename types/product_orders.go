package types

import (
	"time"
)

type ApproveProduct struct {
	ProductID   string `json:"product_id"`
	IsVerified  bool   `json:"is_verified"`
}

type Product struct {
	ID                 int        `json:"id" db:"id"`
	Img                string     `json:"product_img" db:"product_img"`
	FarmerID           int        `json:"farmer_id" db:"farmer_id"`
	Name               string     `json:"name" db:"name"`
	Type               string     `json:"type" db:"type"`
	Quantity           int        `json:"quantity" db:"quantity"`
	RatePerKg          float64    `json:"rate_per_kg" db:"rate_per_kg"`
	JariSize           *int       `json:"jari_size,omitempty" db:"jari_size"`
	ExpectedDelivery   *time.Time `json:"expected_delivery,omitempty" db:"expected_delivery"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	FarmersPhoneNumber int        `json:"farmer_phone_number" db:"farmers_phone_number"` // -> In case he's using another number for caling
	FarmerFirstName    string     `json:"farmers_first_name,omitempty"`
	FarmerLastName     string     `json:"farmers_last_name,omitempty"`
	IsAvailable        bool       `json:"is_available" db:"is_available"`
	IsVerifiedByAdmin  bool       `json:"is_verified_by_admin" db:"is_verified_by_admin"`
}

type Order struct {
	ID                   int        `json:"id" db:"id"`
	BuyerID              int        `json:"buyer_id" db:"buyer_id"`
	ProductID            int        `json:"product_id" db:"product_id"`
	Quantity             int        `json:"quantity" db:"quantity"`
	TotalPrice           float64    `json:"total_price" db:"total_price"`
	DeliveryAddress      string     `json:"delivery_address" db:"delivery_address"`
	DeliveryCity         string     `json:"delivery_city" db:"delivery_city"`
	DeliveryAddressZIP   int        `json:"delivery_address_zip" db:"delivery_address_pin_code"`
	Status               string     `json:"status" db:"status"`
	ModeOfDelivery       string     `json:"mode_of_delivery" db:"mode_of_delivery"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty" db:"expected_delivery_date"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	BuyersPhoneNumber    int        `json:"buyers_phone_number" db:"buyers_phone_number"` // -> In case he's using another number for caling
}

type OrderSummary struct {
	OrderID              int        `json:"order_id"`
	Quantity             int        `json:"quantity"`
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
}
