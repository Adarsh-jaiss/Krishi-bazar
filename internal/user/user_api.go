package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

func CreateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var u types.User
		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user data")
		}

		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
		u.LastLoginAt = time.Now()

		if _, err := CreateUserStore(db, u); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error creating user: %v", err))
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "user created successfully!"})
	}
}

// NOTE: Currently passing it in parameters, but not sure how will it behvave in app
func GetUserProfile(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user ID")
		}

		res, err := GetUserProfileFromStore(db, userID)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error getting user profile: %v", err))
		}

		return c.JSON(http.StatusOK, res)
	}
}

func UpdateProfile(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user ID")
		}

		var u types.User
		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user data")
		}

		u.UpdatedAt = time.Now()

		if err := UpdateProfileInStore(db, userID, u.PhoneNumber); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error updating user profile: %v", err))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "user profile updated successfully!"})
	}
}
