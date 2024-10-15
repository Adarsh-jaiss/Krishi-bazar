package authy

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	users "github.com/adarsh-jaiss/agrohub/internal/user"
	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

// Save the user data in the temporary store at client side
func HandleSignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		var u types.User
		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid user data")
		}

		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
		u.LastLoginAt = time.Now()

		// Call Authenticate function to send verification code
		if err := Authenticate(u.PhoneNumber); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error sending verification code: %v", err))
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "verification code sent successfully!"})
	}
}

func HandleCompleteSignup(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req types.CompleteSignupRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
		}

		// Verify the code
		if err := VerifyCode(req.User.PhoneNumber, req.VerificationCode); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Verification failed: %v", err))
		}

		// Create user in database
		req.User.CreatedAt = time.Now()
		req.User.UpdatedAt = time.Now()
		req.User.LastLoginAt = time.Now()

		userID, err := users.CreateUserStore(db, req.User)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
		}

		// Create auth record
		if err := CreateAuthRecord(db, userID, req.User.PhoneNumber); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error creating auth record: %v", err))
		}

		// Update auth verification
		if err := UpdateAuthVerification(db, userID, true); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error updating auth verification: %v", err))
		}

		// Get user type
		userType := "buyer"
		if !req.User.IsFarmer {
			userType = "farmer"
		}

		// Generate JWT token
		token, err := GenerateToken(userID, userType)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error generating token: %v", err))
		}

		// return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully!"})
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "User created successfully!",
			"user":    req.User.ID,
			"token":   token,
		})
	}
}

func HandleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req string
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
		}

		// Call Authenticate function to send verification code
		if err := Authenticate(req); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error sending verification code: %v", err))
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "verification code sent successfully!"})

	}

}

func HandleCompleteLogin(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req types.LoginRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
		}

		// Verify the code
		if err := VerifyCode(req.PhoneNumber, req.VerificationCode); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Verification failed: %v", err))
		}

		// Get user from database
		u, err := GetUserByPhoneNumberAndAadharNo(db, req)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
		}

		// Update last login in both auth and users tables
		userID, err := strconv.Atoi(u.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error converting user ID: %v", err))
		}
		if err := UpdateLastLogin(db, userID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error updating last login: %v", err))
		}

		userType := "buyer"
		if u.IsFarmer {
			userType = "farmer"
		}

		// Generate JWT token
		token, err := GenerateToken(userID, userType)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error generating token: %v", err))
		}

		// return c.JSON(http.StatusOK, map[string]string{"message": "user logged in successfully!"})
	
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "User logged in successfully!",
			"token":   token,
			"user":    userID,
		})

	}
}
