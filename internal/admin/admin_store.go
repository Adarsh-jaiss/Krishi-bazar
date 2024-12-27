package admins

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

func GetAdminByID(db *sql.DB, id string) (Admin, error) {
    query := `
    SELECT id,password from admins WHERE username = $1`

    rows, err := db.Query(query, id)
    if err != nil {
        return Admin{}, echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("failed to fetch rows from store: %v", err))
    }
    defer rows.Close()

    var a Admin
    if rows.Next() {
        err = rows.Scan(&a.AdminID,&a.Password)
        if err != nil {
            return Admin{}, fmt.Errorf("error scanning row: %v", err)
        }
        return a, nil
    }
    return Admin{}, fmt.Errorf("admin not found")
}

func GetAllUnapprovedFarmersFromStore(db *sql.DB) ([]types.User, error) {
    query := `
        SELECT u.id, u.img, u.first_name, u.last_name, u.aadhar_number, u.email, u.created_at
        FROM users u 
        JOIN farmers f ON u.id = f.user_id 
        WHERE f.is_verified_by_admin = false;`

    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch unverified farmers: %v", err)
    }
    defer rows.Close()

    var users []types.User
    for rows.Next() {
        var user types.User
        var nullableImage sql.NullString
        err = rows.Scan(&user.ID, &nullableImage, &user.FirstName, &user.LastName, &user.AadharNumber, &user.Email, &user.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }
        user.Image = nullableImage.String
        users = append(users, user)
    }

    if len(users) == 0 {
        return nil, fmt.Errorf("no unverified farmers found")
    }

    return users, nil
}

func ApproveUserStore(db *sql.DB, userID int) error {
    query := `
    UPDATE farmers
    SET is_verified_by_admin = true
    WHERE user_id = $1;
    `

    updateUsersQuery := `
    UPDATE users
    SET updated_at = NOW()
    WHERE id = $1;
    `

    // Execute the farmers update query
    _, err := db.Exec(query, userID)
    if err != nil {
        return fmt.Errorf("error updating is_verified field in userstore: %v", err)
    }

    // Execute the users update query
    _, err = db.Exec(updateUsersQuery, userID)
    if err != nil {
        return fmt.Errorf("error updating updated_at field in userstore: %v", err)
    }

    return nil
}

func ApproveProductInStore(db *sql.DB, v types.ApproveProduct) error {
    q := `
    UPDATE products
    SET is_verified_by_admin = $1, updated_at = NOW()
    WHERE product_id = $2;
    `

    // Execute the products update query
    _, err := db.Exec(q, v.IsVerified, v.ProductID)
    if err != nil {
        return fmt.Errorf("error updating is_approved field in products: %v", err)
    }

    return nil
}

func GetUserFromStore(db *sql.DB, userID int) (types.User, error) {
    var user types.User
    var nullableImage sql.NullString

    // First, get the user type
    var userType string
    err := db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&userType)
    if err != nil {
        return types.User{}, fmt.Errorf("error finding user type: %v", err)
    }

    switch userType {
    case "farmer":
        query := `
        SELECT
            u.id, u.first_name, u.last_name, u.email, u.phone_number, u.aadhar_number,
            u.user_type, u.img, u.created_at, u.updated_at,
            f.is_verified_by_admin, f.farm_size, f.address, f.city, f.state, f.pin_code
        FROM users u
        LEFT JOIN farmers f ON u.id = f.user_id
        WHERE u.id = $1`

        err = db.QueryRow(query, userID).Scan(
            &user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.AadharNumber,
            &user.UserType, &nullableImage, &user.CreatedAt, &user.UpdatedAt,
            &user.IsVerified, &user.FarmSize, &user.Address, &user.City, &user.State, &user.PinCode,
        )
    }

    if err != nil {
        return types.User{}, fmt.Errorf("error finding user: %v", err)
    }

    user.Image = nullableImage.String
    user.IsFarmer = (userType == "farmer")

    return user, nil
}
