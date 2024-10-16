package order

import (
	"fmt"
	"net/http"
	"strconv"

	"database/sql"

	users "github.com/adarsh-jaiss/agrohub/internal/user"
	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

// buyerID is sent from frontend!
// NOTE: Extracted UserID from context
func CreateOrder(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var o types.Order
		ProductID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid farmer ID"})
		}

		o.ProductID = ProductID
		o.BuyerID = c.Get("user_id").(int) // used context to fetch userID

		if err := CreateOrderInStore(db, o); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error creating new order:%v", err))
		}

		return c.JSON(200, map[string]string{"message": "order placed successfully!"})
	}
}

func GetOrders(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id").(int)

		u, err := users.GetUserProfileFromStore(db, userID)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error fetching user profile:%v", err))
		}

		if u.UserType == "farmer" {
			orderSummaries, err := GetFarmerOrdersFromStore(db, userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
			}
			return c.JSON(http.StatusOK, orderSummaries)
		} else {
			orderSummaries, err := GetBuyerOrdersFromStore(db, userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
			}

			return c.JSON(http.StatusOK, orderSummaries)
		}
	}
}

func GetOrdersByID(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		orderID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing order id:%v", err))
		}

		res, err := GetOrderFromStore(db, orderID)
		if err != nil {
			echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error fetching order from the store", err))
		}

		return c.JSON(200, res)

	}
}

func UpdateOrderStatus(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
        orderID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing order id:%v", err))
		}
		var orderStatus string
		if err := c.Bind(&orderStatus); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing update request:%v", err))
		}

        if err := UpdateOrderStatusInStore(db,orderID, orderStatus); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error updating order status:%v", err))
		}

        return c.JSON(200, map[string]string{"message":"order status updated successfully!"})

	}
}
