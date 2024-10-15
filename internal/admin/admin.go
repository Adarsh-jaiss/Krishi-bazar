package admins

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

// TODO: add user ID from client side only from the list of all unapproved profiles.
func ApproveUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var a types.Approve
		if err := c.Bind(&a); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user data")
		}

		if err := ApproveUserStore(db, a); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error creating user: %v", err))
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "user approved successfully!"})
	}
}
