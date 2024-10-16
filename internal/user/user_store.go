package users

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/agrohub/types"
)

func UpdateProfileInStore(db *sql.DB, userID int, newPhoneNumber string) error {
    query := `
        UPDATE users
        SET phone_number = $1, updated_at = NOW()
        WHERE id = $2
    `

    _, err := db.Exec(query, newPhoneNumber, userID)
    if err != nil {
        return fmt.Errorf("error updating phone number: %v", err)
    }

    return nil
}
func GetUserProfileFromStore(db *sql.DB, userID int) (types.User, error) {
    var user types.User

    // First, get the user type
    var userType string
    err := db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&userType)
    if err != nil {
        return types.User{}, fmt.Errorf("error finding user type: %v", err)
    }

    // Base query for user information
    baseQuery := `
    SELECT
        u.id, u.first_name, u.last_name, u.email, u.phone_number, u.aadhar_number,
        u.user_type, u.img, u.created_at, u.updated_at, u.last_login_at`

    // Additional fields and join based on user type
    var additionalFields string
    var joinClause string
    var scanArgs []interface{}

    switch userType {
    case "farmer":
        additionalFields = `, f.is_verified_by_admin, f.farm_size, f.address, f.city, f.state, f.pin_code`
        joinClause = ` LEFT JOIN farmers f ON u.id = f.user_id`
        scanArgs = []interface{}{
            &user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.AadharNumber,
            &user.UserType, &user.Image, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
            &user.IsVerified, &user.FarmSize, &user.Address, &user.City, &user.State, &user.PinCode,
        }
    case "buyer":
        additionalFields = `, b.address, b.city, b.state, b.pin_code`
        joinClause = ` LEFT JOIN buyers b ON u.id = b.user_id`
        scanArgs = []interface{}{
            &user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.AadharNumber,
            &user.UserType, &user.Image, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
            &user.Address, &user.City, &user.State, &user.PinCode,
        }
    default:
        // For admin or any other user type, we just use the base query
        scanArgs = []interface{}{
            &user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.AadharNumber,
            &user.UserType, &user.Image, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
        }
    }

    // Combine the query parts
    fullQuery := baseQuery + additionalFields + " FROM users u" + joinClause + " WHERE u.id = $1"

    // Execute the query
    err = db.QueryRow(fullQuery, userID).Scan(scanArgs...)
    if err != nil {
        return types.User{}, fmt.Errorf("error finding user: %v", err)
    }

    // Set IsFarmer based on user type
    user.IsFarmer = (userType == "farmer")

    return user, nil
}


// TODO : manage user_type in api as farmer or buyer
func CreateUserStore(db *sql.DB, user types.User) (int, error) {
    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        return 0, fmt.Errorf("error starting transaction: %v", err)
    }
    defer tx.Rollback() // Rollback the transaction if it hasn't been committed

    // Insert into users table and get the new user ID
    insertUserQuery := `
    INSERT INTO users (
        first_name, last_name, email, phone_number, aadhar_number,
        user_type, img, created_at, updated_at, last_login_at
    ) VALUES (
        $1, $2, $3, $4, $5,
        CASE WHEN $6 THEN 'farmer'::user_type ELSE 'buyer'::user_type END,
        $7, $8, $9, $10
    )
    RETURNING id;
    `

    var newUserID int
    err = tx.QueryRow(insertUserQuery,
        user.FirstName,
        user.LastName,
        user.Email,
        user.PhoneNumber,
        user.AadharNumber,
        user.IsFarmer,
        user.Image,
        user.CreatedAt,
        user.UpdatedAt,
        user.LastLoginAt,
    ).Scan(&newUserID)
    if err != nil {
        return 0, fmt.Errorf("error creating user in userstore: %v", err)
    }

    // Insert into farmers or buyers table based on user type
    if user.IsFarmer {
        insertFarmerQuery := `
        INSERT INTO farmers (
            user_id, is_verified_by_admin, farm_size,
            address, city, state, pin_code
        ) VALUES (
            $1, $2, $3::FLOAT,
            $4, $5, $6, $7
        );
        `
        _, err = tx.Exec(insertFarmerQuery,
            newUserID,
            user.IsVerified,
            user.FarmSize,
            user.Address,
            user.City,
            user.State,
            user.PinCode,
        )
        if err != nil {
            return 0, fmt.Errorf("error creating farmer in userstore: %v", err)
        }
    } else {
        insertBuyerQuery := `
        INSERT INTO buyers (
            user_id, address, city, state, pin_code
        ) VALUES (
            $1, $2, $3, $4, $5
        );
        `
        _, err = tx.Exec(insertBuyerQuery,
            newUserID,
            user.Address,
            user.City,
            user.State,
            user.PinCode,
        )
        if err != nil {
            return 0, fmt.Errorf("error creating buyer in userstore: %v", err)
        }
    }

    // Commit the transaction
    if err = tx.Commit(); err != nil {
        return 0, fmt.Errorf("error committing transaction: %v", err)
    }

    return newUserID, nil
}



