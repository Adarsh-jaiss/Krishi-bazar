package product

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func ListAllProducts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(200, "List of all products")
	}
}

func ListJariProducts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(200, "List of Jari products")
	}
}

func ListMushroomProducts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200,"list of mushroom products")
	}
}

func GetProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200,"list of mushroom products")
	}
}

func CreateProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200,"list of mushroom products")
	}
}

func UpdateProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200,"list of mushroom products")
	}
}

func DeleteProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200,"list of mushroom products")
	}
}




