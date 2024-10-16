package db

import (
	"database/sql"
	"fmt"
)

func DropTable(db *sql.DB, tableName string) error {
	// Prepare the SQL statement
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop table %s: %v", tableName, err)
	}

	fmt.Printf("Table %s dropped successfully\n", tableName)
	return nil
}

// TODO: Mark Mobile number as unique
func CreateTable() error {
	db, err := Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	createUserTypeEnum := `
	CREATE TYPE user_type AS ENUM ('buyer', 'farmer', 'admin');`

	createOrderStatusEnum := `DO $$ BEGIN
            CREATE TYPE order_status AS ENUM ('Pending', 'Processing', 'Shipped', 'Delivered', 'Cancelled', 'Refunded');
        EXCEPTION
            WHEN duplicate_object THEN null;
        END $$;`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
    	last_name VARCHAR(100) NOT NULL,
    	email VARCHAR(255) UNIQUE,
    	phone_number VARCHAR(20) NOT NULL,
		aadhar_number VARCHAR(12) UNIQUE NOT NULL,
    	user_type user_type NOT NULL,
		img TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login_at TIMESTAMP
	);`

	createFarmersTable := `
	CREATE TABLE IF NOT EXISTS farmers (
		user_id INT PRIMARY KEY REFERENCES users(id),
		is_verified_by_admin BOOLEAN DEFAULT FALSE,
		farm_size FLOAT NOT NULL, -- Value in acres
		address TEXT NOT NULL,
		city VARCHAR(100) NOT NULL,
		state VARCHAR(100) NOT NULL,
		pin_code VARCHAR(10) NOT NULL	
	);
	`

	createBuyersTable := `
	CREATE TABLE IF NOT EXISTS buyers (
    	user_id INT PRIMARY KEY REFERENCES users(id),
    	address TEXT NOT NULL,
    	city VARCHAR(100) NOT NULL,
    	state VARCHAR(100) NOT NULL,
    	pin_code VARCHAR(10) NOT NULL
	);`

	createAdminsTable := `
	CREATE TABLE IF NOT EXISTS admins (
    	user_id INT PRIMARY KEY REFERENCES users(id),
    	admin_level INT NOT NULL
	);`

	createAuthTable := `
	CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    phone_number VARCHAR(15) NOT NULL,
    verification_code VARCHAR(10),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	verified_at TIMESTAMP,
    last_login_at TIMESTAMP
);`

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,
	farmer_id INT NOT NULL REFERENCES users(id),
	name VARCHAR(255) NOT NULL,
	type VARCHAR(100) NOT NULL,
	img  TEXT NOT NULL,
	quantity INT NOT NULL,
	rate_per_kg DECIMAL(10, 2) NOT NULL,
	jari_size INT,
	expected_delivery DATE,
	farmers_phone_number VARCHAR(15) NOT NULL,
	is_available BOOLEAN DEFAULT TRUE,
	is_verified_by_admin BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	buyer_id INT NOT NULL REFERENCES users(id),
	product_id INT NOT NULL REFERENCES products(id),
	buyers_phone_number VARCHAR(15) NOT NULL,
	quantity INT NOT NULL,
	total_price DECIMAL(10, 2) NOT NULL,
	status order_status NOT NULL,
	mode_of_delivery VARCHAR(100),
	expected_delivery_date DATE,
	delivery_address TEXT NOT NULL,
	delivery_city VARCHAR(100) NOT NULL,
	delivery_address_zip VARCHAR(10) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	indexes := `
	CREATE INDEX idx_farmer_id ON products(farmer_id);
	CREATE INDEX idx_buyer_id ON orders(buyer_id);
	CREATE INDEX idx_product_id ON orders(product_id);
	`
	enums := []string{createUserTypeEnum, createOrderStatusEnum}
	for i := 0; i < len(enums); i++ {
		_, err = db.Exec(enums[i])
		if err != nil {
			return fmt.Errorf("failed to create user type enum: %v", err)
		}
	}

	tables := []string{createUsersTable, createFarmersTable, createBuyersTable, createAdminsTable, createAuthTable, createProductsTable, createOrdersTable}
	for i := 0; i < len(tables); i++ {
		_, err := db.Exec(tables[i])
		if err != nil {
			return fmt.Errorf("error creating %v: %v", tables[i], err)
		}
	}

	_, err = db.Exec(indexes)
	if err != nil {
		return fmt.Errorf("failed to created indexes:%v", err)
	}

	// Create trigger for updating 'updated_at' automatically
	_, err = db.Exec(`
        CREATE OR REPLACE FUNCTION update_modified_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = now();
            RETURN NEW;
        END;
        $$ language 'plpgsql';

        CREATE TRIGGER update_product_modtime
            BEFORE UPDATE ON products
            FOR EACH ROW
            EXECUTE FUNCTION update_modified_column();

        CREATE TRIGGER update_order_modtime
            BEFORE UPDATE ON orders
            FOR EACH ROW
            EXECUTE FUNCTION update_modified_column();
    `)
	if err != nil {
		return fmt.Errorf("failed to create update trigger: %v", err)
	}

	fmt.Println("Tables,enums,indexes and triggers created successfully!")
	return nil
}
