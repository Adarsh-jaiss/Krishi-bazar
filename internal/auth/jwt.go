package authy

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GenerateToken(userID int, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 360).Unix(), // Token expires in 72 hours
	}

	fmt.Println(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ExtractUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))
		
		// Store user ID in context
		c.Set("user_id", userID)
		
		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := c.Get("user").(*jwt.Token)
        claims := user.Claims.(jwt.MapClaims)
        userType := claims["user_type"].(string)

        if userType != "admin" {
            return echo.ErrUnauthorized
        }
        return next(c)
    }
}

func IsFarmer(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := c.Get("user")
        if user == nil {
            return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid JWT token")
        }

        token, ok := user.(*jwt.Token)
        if !ok {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT claims")
        }

        userType, ok := claims["user_type"].(string)
        if !ok {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user type in JWT")
        }

        if userType != "farmer" {
            return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf( "Access denied: User is not a farmer : %v",userType))
        }

        return next(c)
    }
}
