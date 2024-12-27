package admins

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	authy "github.com/adarsh-jaiss/agrohub/internal/auth"
	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

type Admin struct {
	AdminID  int `json:"admin_id,omitempty" db:"id"`
	UserName string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func AdminLogin(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var a Admin
		if err := c.Bind(&a); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid Request")
		}

		res, err := GetAdminByID(db, a.UserName)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error finding admin: %v", err))
		}

		fmt.Printf("Login attempt - Username: %s, Admin details: %+v\n", a.UserName, a.Password)

		fmt.Printf("%v:%v\n%v:%v\n",a.UserName,res.UserName,a.Password,res.Password)

		if a.Password == res.Password {
			// Generate JWT token
			token, err := authy.GenerateToken(res.AdminID, "admin")
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error generating token: %v", err))
			}

			// return c.JSON(http.StatusOK, map[string]string{"message": "user logged in successfully!"})

			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "User logged in successfully!",
				"token":   token,
				"user":    res.AdminID,
			})

		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}
}

func GetUserProfile(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))		// using path parameter -> for query use c.query("id")
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("Invalid user ID: %v", err))
		}

		res, err := GetUserFromStore(db, userID)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error getting user profile: %v", err))
		}

		return c.JSON(http.StatusOK, res)
	}
}

func GetAllUnapprovedFarmers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := GetAllUnapprovedFarmersFromStore(db)
		if err!= nil {
			return  echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error fetching users: %v", err))
		}

		return c.JSON(http.StatusFound, res)
	}
}

func ApproveUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))	
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user data")
		}

		if err := ApproveUserStore(db, userID); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error approving user: %v", err))
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "user approved successfully!"})
	}
}

func ApproveProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var a types.ApproveProduct
		if err := c.Bind(&a); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user data")
		}

		a.IsVerified = true

		if err := ApproveProductInStore(db, a); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error approving product: %v", err))
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "product approved successfully!"})
	}
}
