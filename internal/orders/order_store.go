package order

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/agrohub/internal/product"
	"github.com/adarsh-jaiss/agrohub/types"
)

func GetOrderFromStore(db *sql.DB, orderID int) (types.OrderSummary, error) {
	var order types.OrderSummary
	var expectedDeliveryDate sql.NullTime
	

	query := `
		SELECT 
			o.id, o.quantity_in_kg, o.total_price, o.status, o.mode_of_delivery, 
			o.expected_delivery_date, o.created_at, o.product_id, p.name,p.img,
			u.id, u.first_name, u.last_name, u.phone_number,
			o.delivery_address, o.delivery_city, o.delivery_address_zip,
			o.buyers_phone_number, p.farmers_phone_number
		FROM 
			orders o
		JOIN 
			products p ON o.product_id = p.id
		JOIN 
			users u ON o.buyer_id = u.id
		WHERE 
			o.id = $1
	`

	err := db.QueryRow(query, orderID).Scan(
		&order.OrderID, &order.QuantityInKg, &order.TotalPrice, &order.Status, &order.ModeOfDelivery,
		&expectedDeliveryDate, &order.OrderDate, &order.ProductID, &order.ProductName,&order.ProductImg,
		&order.UserID, &order.UserFirstName, &order.UserLastName, &order.UserPhoneNumber,
		&order.DeliveryAddress, &order.DeliveryCity, &order.DeliveryAddressZIP,
		&order.BuyersPhoneNumber, &order.FarmersPhoneNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("no order found with ID %d", orderID)
		}
		return order, fmt.Errorf("error querying order: %v", err)
	}

	if expectedDeliveryDate.Valid {
		order.ExpectedDeliveryDate = &expectedDeliveryDate.Time
	}

	return order, nil
}

func isValidOrderStatus(status string) bool {
	// Replace this with the actual allowed enum values from your database
	allowedStatuses := map[string]bool{
		"pending":   true,
		"approved":  true,
		"shipped":   true,  // Ensure this is the correct enum value
		"delivered": true,
		"cancelled": true,
	}

	// Check if the provided status exists in the allowed statuses
	_, valid := allowedStatuses[status]
	return valid
}


func UpdateOrderStatusInStore(db *sql.DB, orderID int, status string) error {
	if !isValidOrderStatus(status) {
		return fmt.Errorf("invalid status: %s", status)
	}

	query := `
		UPDATE orders 
		SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2
	`

	result, err := db.Exec(query, status, orderID)
	if err != nil {
		return fmt.Errorf("error updating order status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no order found with ID %d", orderID)
	}

	return nil
}

func CreateOrderInStore(db *sql.DB, order types.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// get prouduct details
	p, err := product.GetProductFromStore(db, order.ProductID)
	if err != nil {
		return fmt.Errorf("unable to fetch product :%v", err)
	}

	if p.Quantity < order.QuantityInKg {
		return fmt.Errorf("insufficient quantity available")
	}

	// Calculate total price
	order.TotalPrice = float64(order.QuantityInKg) * p.RatePerKg

	err = tx.QueryRow(`
		INSERT INTO orders (buyer_id, product_id, quantity_in_kg, total_price, status, mode_of_delivery, expected_delivery_date, delivery_address, delivery_city, delivery_address_zip, buyers_phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`, order.BuyerID, order.ProductID, order.QuantityInKg, order.TotalPrice, "pending", order.ModeOfDelivery, order.ExpectedDeliveryDate ,order.DeliveryAddress, order.DeliveryCity, order.DeliveryAddressZIP, order.BuyersPhoneNumber).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting order: %v", err)
	}

	// Update product quantity
	_, err = tx.Exec("UPDATE products SET quantity_in_kg = quantity_in_kg - $1 WHERE id = $2", order.QuantityInKg, order.ProductID)
	if err != nil {
		return fmt.Errorf("error updating product quantity: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil

}

// GetOrdersBasedOnUser fetches orders based on whether the user is a buyer or a farmer.
func GetOrdersBasedOnUser(db *sql.DB, userID int, userType string) ([]types.OrderStatus, error) {
	var query string

	// Define different queries based on whether the user is a buyer or farmer
	if userType == "farmer" {
		query = `
			SELECT 
				o.id, o.quantity_in_kg, o.total_price, o.status, o.mode_of_delivery, 
				o.expected_delivery_date, o.created_at, 
				p.id, p.name,p.img, 
				b.first_name, b.last_name, b.phone_number, o.delivery_address, o.delivery_city, o.delivery_address_zip
			FROM 
				orders o
			JOIN 
				products p ON o.product_id = p.id
			JOIN 
				users b ON o.buyer_id = b.id -- Join to get the buyer's details
			WHERE 
				p.farmer_id = $1
			ORDER BY 
				o.created_at DESC
		`
	} else if userType == "buyer" {
		query = `
			SELECT 
				o.id, o.quantity_in_kg, o.total_price, o.status, o.mode_of_delivery, 
				o.expected_delivery_date, o.created_at, 
				p.id, p.name,p.img
				f.first_name, f.last_name, f.phone_number,
				o.delivery_address, o.delivery_city, o.delivery_address_zip
			FROM 
				orders o
			JOIN 
				products p ON o.product_id = p.id
			JOIN 
				users f ON p.farmer_id = f.id -- Join to get the farmer's details
			WHERE 
				o.buyer_id = $1
			ORDER BY 
				o.created_at DESC
		`
	} else {
		return nil, fmt.Errorf("invalid user type")
	}

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}
	defer rows.Close()

	var orders []types.OrderStatus
	for rows.Next() {
		var o types.OrderStatus
		var expectedDeliveryDate sql.NullTime

		if userType == "farmer" {
			// Scan for farmer-specific data (including buyer details)
			err := rows.Scan(
				&o.OrderDetails.OrderID, &o.OrderDetails.QuantityInKg, &o.OrderDetails.TotalPrice, &o.OrderDetails.Status,
				&o.OrderDetails.ModeOfDelivery, &expectedDeliveryDate, &o.OrderDetails.OrderDate,
				&o.OrderDetails.ProductID, &o.OrderDetails.ProductName,&o.OrderDetails.ProductImg,
				&o.BuyersDetails.BuyerFirstName, &o.BuyersDetails.BuyerLastName,
				&o.BuyersDetails.BuyerPhoneNumber, &o.BuyersDetails.DeliveryAddress, &o.BuyersDetails.DeliveryCity, &o.BuyersDetails.DeliveryZIP,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan order: %v", err)
			}
		} else if userType == "buyer" {
			// Scan for buyer-specific data (no buyer details, just the order and product info)
			err := rows.Scan(
				&o.OrderDetails.OrderID, &o.OrderDetails.QuantityInKg, &o.OrderDetails.TotalPrice, &o.OrderDetails.Status,
				&o.OrderDetails.ModeOfDelivery, &expectedDeliveryDate, &o.OrderDetails.OrderDate,
				&o.OrderDetails.ProductID, &o.OrderDetails.ProductName,&o.OrderDetails.ProductImg,
				&o.SellerDetails.FarmerFirstName, &o.SellerDetails.FarmerLastName, &o.SellerDetails.FarmerPhoneNumber,
				&o.DeliveryAddress, &o.DeliveryCity, &o.DeliveryZIP,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan order: %v", err)
			}
		}

		// Handle the nullable time value
        if expectedDeliveryDate.Valid {
            o.OrderDetails.ExpectedDeliveryDate = &expectedDeliveryDate.Time
        } else {
            o.OrderDetails.ExpectedDeliveryDate = nil
        }

		orders = append(orders, o)
	}

	return orders, nil
}