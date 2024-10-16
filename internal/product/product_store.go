package product

import (
	"database/sql"

	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

func UpdateProductAvailabilityInStore(db *sql.DB,ProductID int, availabilty bool) error {
	q := `
	UPDATE products
	SET is_available = $1
	WHERE id = $2;`

	_, err := db.Exec(q, availabilty, ProductID)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to update product availability in store :%v", err)
	}
	return nil
}

func GetAllProductsFromStore(db *sql.DB) ([]types.Product, error) {
	q := `
		SELECT p.id, p.farmer_id, p.name, p.type, p.img, p.quantity, 
    	p.rate_per_kg, p.jari_size, p.expected_delivery, 
    	p.farmers_phone_number, p.created_at, p.updated_at,
    	u.first_name AS farmer_first_name, u.last_name AS farmer_last_name
		FROM 
		    products p
		JOIN 
		    users u ON p.farmer_id = u.id
		ORDER BY 
		    p.created_at DESC;`

	rows, err := db.Query(q)
	if err != nil {
		return nil, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to fetch rows from store :%v", err)
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(
			&p.ID, &p.FarmerID, &p.Name, &p.Type, &p.Img, &p.Quantity,
			&p.RatePerKg, &p.JariSize, &p.ExpectedDelivery,
			&p.FarmersPhoneNumber, &p.CreatedAt, &p.UpdatedAt,
			&p.FarmerFirstName, &p.FarmerLastName,
		); err != nil {
			return nil, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to scan rows: %v", err)
		}
		products = append(products, p)
	}

	return products, nil

}

func GetAllMushroomAndJariProductsFromStore(db *sql.DB, productType string) ([]types.Product, error) {
	q := `
		SELECT p.id, p.farmer_id, p.name, p.type, p.img, p.quantity, 
    	p.rate_per_kg, p.jari_size, p.expected_delivery, 
    	p.farmers_phone_number, p.created_at, p.updated_at,
    	u.first_name AS farmer_first_name, u.last_name AS farmer_last_name
		FROM 
		    products p
		JOIN 
		    users u ON p.farmer_id = u.id
		ORDER BY 
		    p.created_at DESC;`

	switch productType {
	case "jari":
		q += " WHERE p.type = 'jari'"
	case "mushroom":
		q += " WHERE p.type = 'mushroom'"
	}
	q += " ORDER BY p.created_at DESC"

	rows, err := db.Query(q)
	if err != nil {
		return nil, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to fetch rows from store :%v", err)
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(
			&p.ID, &p.FarmerID, &p.Name, &p.Type, &p.Img, &p.Quantity,
			&p.RatePerKg, &p.JariSize, &p.ExpectedDelivery,
			&p.FarmersPhoneNumber, &p.CreatedAt, &p.UpdatedAt,
			&p.FarmerFirstName, &p.FarmerLastName,
		); err != nil {
			return nil, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to scan rows: %v", err)
		}
		products = append(products, p)
	}

	return products, nil

}

func GetProductFromStore(db *sql.DB, ProductID int) (types.Product, error) {
	q := `
	SELECT p.id, p.farmer_id, p.name, p.type, p.img, p.quantity, 
	p.rate_per_kg, p.jari_size, p.expected_delivery, 
	p.farmers_phone_number, p.created_at, p.updated_at,
	u.first_name AS farmer_first_name, u.last_name AS farmer_last_name
	FROM 
		products p
	JOIN 
		users u ON p.farmer_id = u.id
	WHERE 
		p.id = $1;`

	var p types.Product
	if err := db.QueryRow(q, ProductID).Scan(
		&p.ID, &p.FarmerID, &p.Name, &p.Type, &p.Img, &p.Quantity,
		&p.RatePerKg, &p.JariSize, &p.ExpectedDelivery,
		&p.FarmersPhoneNumber, &p.CreatedAt, &p.UpdatedAt,
		&p.FarmerFirstName, &p.FarmerLastName,
	); err != nil {
		return types.Product{}, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to scan rows: %v", err)
	}

	return p, nil
}

func GetFarmersProductFromStore(db *sql.DB, FarmerID int) ([]types.Product, error) {
	q := `
	SELECT p.id, p.farmer_id, p.name, p.type, p.img, p.quantity, 
	p.rate_per_kg, p.jari_size, p.expected_delivery, 
	p.farmers_phone_number, p.created_at, p.updated_at,
	u.first_name AS farmer_first_name, u.last_name AS farmer_last_name
	FROM 
		products p
	JOIN 
		users u ON p.farmer_id = u.id
	WHERE
		p.farmer_id = $1
	ORDER BY 
		p.created_at DESC;`

	rows, err := db.Query(q, FarmerID)
	if err != nil {
		return nil, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to fetch rows from store :%v", err)
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(
			&p.ID, &p.FarmerID, &p.Name, &p.Type, &p.Img, &p.Quantity,
			&p.RatePerKg, &p.JariSize, &p.ExpectedDelivery,
			&p.FarmersPhoneNumber, &p.CreatedAt, &p.UpdatedAt,
			&p.FarmerFirstName, &p.FarmerLastName,
		); err != nil {
			return nil, echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to scan rows: %v", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func CreateProductInStore(db *sql.DB, p types.Product) error {
	q := `
	INSERT INTO products (farmer_id, name, type, img, quantity, rate_per_kg, jari_size, expected_delivery, farmers_phone_number)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	_, err := db.Exec(q, p.FarmerID, p.Name, p.Type, p.Img, p.Quantity, p.RatePerKg, p.JariSize, p.ExpectedDelivery, p.FarmersPhoneNumber)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to insert product in store :%v", err)
	}
	return nil
}

func DeleteProductFromStore(db *sql.DB, ProductID int) error {
	q := `
	DELETE FROM products
	WHERE id = $1;`

	_, err := db.Exec(q, ProductID)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to delete product from store :%v", err)
	}
	return nil
}