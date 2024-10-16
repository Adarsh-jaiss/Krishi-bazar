package admins

import (
	"fmt"
	"database/sql"

	"github.com/adarsh-jaiss/agrohub/types"
)

func ApproveUserStore(db *sql.DB, v types.Approve) error {
    query := `
    UPDATE farmers
    SET is_verified = $1
    WHERE user_id = $2;
    `

    updateUsersQuery := `
    UPDATE users
    SET updated_at = NOW()
    WHERE id = $2;
    `

    // Execute the farmers update query
    _, err := db.Exec(query, v.IsVerified, v.UserID)
    if err != nil {
        return fmt.Errorf("error updating is_verified field in userstore: %v", err)
    }

    // Execute the users update query
    _, err = db.Exec(updateUsersQuery, v.UserID)
    if err != nil {
        return fmt.Errorf("error updating updated_at field in userstore: %v", err)
    }

    return nil
}

func ApproveProductInStore(db *sql.DB, v types.ApproveProduct) error {
    q := `
    UPDATE products
    SET is_approved = $1, updated_at = NOW()
    WHERE product_id = $2;
    `

    // Execute the products update query
    _, err := db.Exec(q, v.IsVerified, v.ProductID)
    if err != nil {
        return fmt.Errorf("error updating is_approved field in products: %v", err)
    }

    return nil
}