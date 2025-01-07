package product

import (
	"database/sql"
	"net/http"
	"strconv"

	"fmt"

	"github.com/adarsh-jaiss/agrohub/types"
	"github.com/labstack/echo/v4"
)

func UpdateProductAvailability(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing product id :%v", err))
		}
		if err := UpdateProductAvailabilityInStore(db, id, false); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("failed to fetch products from store: %v", err))
		}

		return c.JSON(200, map[string]string{
			"message": "Product availability updated successfully!",
		})

	}
}

func ListAllProducts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := GetAllProductsFromStore(db)
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("failed to fetch products from store: %+v", err))
		}
		return c.JSON(200, res)
	}
}

func ListJariProducts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := GetAllMushroomAndJariProductsFromStore(db,"Jari")
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("failed to fetch jari products from store: %v", err))
		}
		return c.JSON(200, res)
	}
}

func ListMushroomProducts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := GetAllMushroomAndJariProductsFromStore(db,"Mushroom")
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("failed to fetch mushroom products from store: %v", err))
		}
		return c.JSON(200, res)
	}
}

func GetProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ProductID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing product id :%v", err))
		}

		res, err := GetProductFromStore(db, ProductID)
		if err != nil {
			echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("unable to fetch the product from store :%v", err))
		}

		return c.JSON(200, res)
	}
}

// TRY fetching UserID it from context --> Need to see frontend Implementation
func ListAllProductsOfFarmer(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		FarmerID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing farmer id :%v", err))
		}

		res, err := GetFarmersProductFromStore(db, FarmerID)
		if err != nil {
			echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("unable to fetch the products from store :%v", err))
		}

		return c.JSON(200, res)
	}
}

// set farmerID from the client side
// TRY fetching UserID it from context --> Need to see frontend Implementation
func CreateProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var p types.Product
		FarmerID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing farmer id :%v", err))
		}

		p.FarmerID = FarmerID

		if err := c.Bind(&p); err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("error parsing create request :%v", err))
		}
		if err := CreateProductInStore(db, &p); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error creating new product:%v", err))
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "product created successfully!"})
	}
}

// func UpdateProduct(db *sql.DB) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		ProductID, err := strconv.Atoi(c.Param("id"))
// 		if err!= nil {
// 			return echo.NewHTTPError(echo.ErrBadGateway.Code, fmt.Sprintf("error parsing update request"))
// 		}

// 		var p types.Product

// 		if err := UpdateProductInStore(db,ProductID)

// 		return c.JSON(200, "list of mushroom products")
// 	}
// }

func DeleteProduct(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ProductID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadGateway.Code, fmt.Sprintf("error parsing update request"))
		}

		if err := DeleteProductFromStore(db, ProductID); err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Sprintf("error deleting product :%v", err))
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "product deleted successfully!"})
	}
}
