package authy

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/agrohub/types"
)


func GetUserByPhoneNumberAndAadharNo(db *sql.DB, v types.LoginRequest) (types.User, error) {
    query := `SELECT * FROM users WHERE aadhar_number = $1 AND phone_number = $2;`

    var user types.User
    err := db.QueryRow(query, v.AadharNumber, v.PhoneNumber).Scan(
        &user.ID,
        &user.FirstName,
        &user.LastName,
        &user.Email,
        &user.PhoneNumber,
        &user.IsFarmer,
        &user.Image,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.LastLoginAt,
    )
    if err != nil {
        return user, fmt.Errorf("error finding user: %v", err)
    }

    return user, nil
}

func UpdateLastLogin(db *sql.DB, userID int) error {
    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("error starting transaction: %v", err)
    }
    defer tx.Rollback() // Rollback the transaction if it hasn't been committed

    // Update auth table
    authQuery := `
    UPDATE auth
    SET last_login_at = CURRENT_TIMESTAMP
    WHERE user_id = $1
    `
    _, err = tx.Exec(authQuery, userID)
    if err != nil {
        return fmt.Errorf("error updating auth table: %v", err)
    }

    // Update users table
    userQuery := `
    UPDATE users
    SET last_login_at = CURRENT_TIMESTAMP
    WHERE id = $1
    `
    _, err = tx.Exec(userQuery, userID)
    if err != nil {
        return fmt.Errorf("error updating users table: %v", err)
    }

    // Commit the transaction
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("error committing transaction: %v", err)
    }

    return nil
}

func CreateAuthRecord(db *sql.DB, userID int, mobileNumber string) error {
    query := `
    INSERT INTO auth (user_id, mobile_number)
    VALUES ($1, $2)
    `
    _, err := db.Exec(query, userID, mobileNumber)
    return err
}

func UpdateAuthVerification(db *sql.DB, userID int, isVerified bool) error {
    query := `
    UPDATE auth
    SET is_verified = $2, verified_at = CASE WHEN $2 = true THEN CURRENT_TIMESTAMP ELSE NULL END
    WHERE user_id = $1
    `
    _, err := db.Exec(query, userID, isVerified)
    return err
}
