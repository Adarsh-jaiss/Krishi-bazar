package types

import "time"

type Order struct {
	ID                   int        `json:"id" db:"id"`
	BuyerID              int        `json:"buyer_id" db:"buyer_id"`
	ProductID            int        `json:"product_id" db:"product_id"`
	QuantityInKg         int        `json:"quantity_in_kg" db:"quantity_in_kg"`
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

// OrderDetails struct contains the details related to the order itself.
type OrderDetails struct {
	OrderID              int        `json:"order_id"`
	QuantityInKg         int        `json:"quantity_in_kg"`
	TotalPrice           float64    `json:"total_price"`
	Status               string     `json:"status"`
	ModeOfDelivery       string     `json:"mode_of_delivery"`
	ExpectedDeliveryDate time.Time `json:"expected_delivery_date,omitempty"`
	OrderDate            time.Time  `json:"order_date"`
	ProductID            int        `json:"product_id"`
	ProductName          string     `json:"product_name"`
}

// BuyersDetails struct contains information related to the buyer.
type BuyersDetails struct {
	BuyerFirstName    string `json:"buyer_first_name"`
	BuyerLastName     string `json:"buyer_last_name"`
	BuyerPhoneNumber  string `json:"buyer_phone_number"`
	DeliveryAddress   string `json:"delivery_address"`
	DeliveryCity      string `json:"delivery_city"`
	DeliveryZIP       int    `json:"delivery_zip"`
}

// SellerDetails struct contains information related to the seller (farmer).
type SellerDetails struct {
	FarmerFirstName   string `json:"farmer_first_name"`
	FarmerLastName    string `json:"farmer_last_name"`
	FarmerPhoneNumber string `json:"farmer_phone_number"`
}

// OrderStatus struct combines OrderDetails, BuyersDetails, and SellerDetails.
type OrderStatus struct {
	OrderDetails  `json:"order_details"`
	BuyersDetails `json:"buyer_details"`
	SellerDetails `json:"seller_details"`
}
