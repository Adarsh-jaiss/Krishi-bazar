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

	// createUserTypeEnum := `
	// CREATE TYPE user_type AS ENUM ('buyer', 'farmer', 'admin');`

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
    mobile_number VARCHAR(15) NOT NULL,
    verification_code VARCHAR(10),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	verified_at TIMESTAMP,
    last_login_at TIMESTAMP
);`

	// _, err = db.Exec(createUserTypeEnum)
	// if err != nil {
	// 	return fmt.Errorf("failed to create user type enum: %v", err)
	// }

	_, err = db.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)

	}

	_, err = db.Exec(createFarmersTable)
	if err != nil {
		return fmt.Errorf("failed to create farmers table: %v", err)
	}

	_, err = db.Exec(createBuyersTable)
	if err != nil {
		return fmt.Errorf("failed to create buyers table: %v", err)
	}

	_, err = db.Exec(createAdminsTable)
	if err != nil {
		return fmt.Errorf("failed to create admins table: %v", err)
	}

	_, err = db.Exec(createAuthTable)
	if err != nil {
		return fmt.Errorf("failed to create auth table: %v", err)
	}

	fmt.Println("Tables created successfully!")
	return nil
}
