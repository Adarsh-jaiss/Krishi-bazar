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
			o.id, o.quantity, o.total_price, o.status, o.mode_of_delivery, 
			o.expected_delivery_date, o.created_at, o.product_id, p.name,
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
		&order.OrderID, &order.Quantity, &order.TotalPrice, &order.Status, &order.ModeOfDelivery,
		&expectedDeliveryDate, &order.OrderDate, &order.ProductID, &order.ProductName,
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

func UpdateOrderStatusInStore(db *sql.DB, orderID int, status string) error {
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

	if p.Quantity < order.Quantity {
		return fmt.Errorf("insufficient quantity available")
	}

	// Calculate total price
	order.TotalPrice = float64(order.Quantity) * p.RatePerKg

	err = tx.QueryRow(`
		INSERT INTO orders (buyer_id, product_id, quantity, total_price, status, delivery_address, delivery_city, delivery_address_zip, buyers_phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`, order.BuyerID, order.ProductID, order.Quantity, order.TotalPrice, "Pending", order.DeliveryAddress, order.DeliveryCity, order.DeliveryAddressZIP, order.BuyersPhoneNumber).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting order: %v", err)
	}

	// Update product quantity
	_, err = tx.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", order.Quantity, order.ProductID)
	if err != nil {
		return fmt.Errorf("error updating product quantity: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil

}

func GetFarmerOrdersFromStore(db *sql.DB, farmerID int) ([]types.OrderSummary, error) {
	query := `
        SELECT 
            o.id, o.quantity, o.total_price, o.status, o.mode_of_delivery, 
            o.expected_delivery_date, o.created_at, 
            p.id, p.name, 
            u.id, u.first_name, u.last_name
        FROM 
            orders o
        JOIN 
            products p ON o.product_id = p.id
        JOIN 
            users u ON o.buyer_id = u.id
        WHERE 
            p.farmer_id = $1
        ORDER BY 
            o.created_at DESC
    `

	rows, err := db.Query(query, farmerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.OrderSummary
	for rows.Next() {
		var o types.OrderSummary
		err := rows.Scan(
			&o.OrderID, &o.Quantity, &o.TotalPrice, &o.Status, &o.ModeOfDelivery,
			&o.ExpectedDeliveryDate, &o.OrderDate,
			&o.ProductID, &o.ProductName,
			&o.UserID, &o.UserFirstName, &o.UserLastName,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func GetBuyerOrdersFromStore(db *sql.DB, buyerID int) ([]types.OrderSummary, error) {
	query := `
        SELECT 
            o.id, o.quantity, o.total_price, o.status, o.mode_of_delivery, 
            o.expected_delivery_date, o.created_at, 
            p.id, p.name, 
            u.id, u.first_name, u.last_name
        FROM 
            orders o
        JOIN 
            products p ON o.product_id = p.id
        JOIN 
            users u ON p.farmer_id = u.id
        WHERE 
            o.buyer_id = $1
        ORDER BY 
            o.created_at DESC
    `

	rows, err := db.Query(query, buyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.OrderSummary
	for rows.Next() {
		var o types.OrderSummary
		err := rows.Scan(
			&o.OrderID, &o.Quantity, &o.TotalPrice, &o.Status, &o.ModeOfDelivery,
			&o.ExpectedDeliveryDate, &o.OrderDate,
			&o.ProductID, &o.ProductName,
			&o.UserID, &o.UserFirstName, &o.UserLastName,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
