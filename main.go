package main

import (
	"log"
	"os"

	"github.com/adarsh-jaiss/agrohub/db"
	admins "github.com/adarsh-jaiss/agrohub/internal/admin"
	"github.com/adarsh-jaiss/agrohub/internal/auth"
	"github.com/adarsh-jaiss/agrohub/internal/orders"
	"github.com/adarsh-jaiss/agrohub/internal/product"
	users "github.com/adarsh-jaiss/agrohub/internal/user"
	"github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	// "fmt"
)

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// }

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	if err = db.CreateTable(); err != nil {
		log.Printf("error creating table : %v", err)
		panic(err)
	}

	defer conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())

	// e.Use(middleware.Recover())
	api := e.Group("/api")

	// Public routes
	auth := api.Group("/auth")
	auth.POST("/signup", authy.HandleSignUp())
	auth.POST("/complete-signup", authy.HandleCompleteSignup(conn))
	auth.POST("/login", authy.HandleLogin())
	auth.POST("/complete-login", authy.HandleCompleteLogin(conn))

	// Admin routes
	admin := api.Group("/admin/approve")
	admin.POST("/user", admins.ApproveUser(conn))
	admin.POST("/product", admins.ApproveProduct(conn))

	// protected routes
	v1 := api.Group("/v1")
	v1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	v1.Use(authy.ExtractUserID)

	// User routes --> store userId locally which is being returned by login
	user := v1.Group("/user")
	user.GET("/:id", users.GetUserProfile(conn)) // -> user/123/profile , for req body i have to use post method and directly send the req body
	user.PUT("/:id/profile", users.UpdateProfile(conn))
	user.POST("/:id/newproduct", product.CreateProduct(conn), authy.IsFarmer)
	user.POST("", users.CreateUser(conn))
	// user.GET("/farmers",users.ListAllFarmers(conn))  //-> see all farmers with their contact details , product and eDOD
	// user.GET("/farmers/:id",users.ListAllFarmers(conn))  -> see a farmer with their contact details and eDOD

	// Product routes
	products := v1.Group("/product")
	products.GET("", product.ListAllProducts(conn))
	products.GET("/farmer/:id", product.ListAllProductsOfFarmer(conn))
	products.GET("/jari", product.ListJariProducts(conn))
	products.GET("/mushroom", product.ListMushroomProducts(conn))
	products.GET("/:id", product.GetProduct(conn))
	products.GET("/:id/mark-unavailable", product.UpdateProductAvailability(conn)) // --> Marks unavailable  --> Manage availabilty and is verified on client side

	products.DELETE("/:id", product.DeleteProduct(conn), authy.IsFarmer)

	// Order routes
	products.POST("/:id/order", order.CreateOrder(conn))
	orders := v1.Group("/orders")
	user.GET("/:id/orders", order.GetOrders(conn))	// -> GET ALL ORDERS
	orders.GET("/:id", order.GetOrdersByID(conn))	// -> GET ORDER BY ID
	orders.PUT("/:id/status", order.UpdateOrderStatus(conn))	

	e.Logger.Fatal(e.Start(":8080"))
}
