package main

import (
	"fmt"
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

	// "github.com/joho/godotenv"
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

	fmt.Printf("DB_USER: %s\n DB_PASSWORD:%ss\n DB_HOST: %s\n DB_PORT=%s\n DB_NAME:%s\n TWILLo_ACCOUNT_SID: %s\n TWILIO_AUTH_TOKEN:%s\n TWILIO_VERIFY_SID: %s\n JWT_SECRET: %s\n", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("TWILLIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), os.Getenv("TWILIO_VERIFY_SID"), os.Getenv("JWT_SECRET"))

	os.Setenv("DB_USER", os.Getenv("DB_USER"))
	os.Setenv("DB_PASSWORD", os.Getenv("DB_PASSWORD"))
	os.Setenv("DB_HOST", os.Getenv("DB_HOST"))
	os.Setenv("DB_PORT", os.Getenv("DB_PORT"))
	os.Setenv("DB_NAME", os.Getenv("DB_NAME"))
	os.Setenv("TWILLIO_ACCOUNT_SID", os.Getenv("TWILLIO_ACCOUNT_SID"))
	os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))
	os.Setenv("TWILIO_VERIFY_SID", os.Getenv("TWILIO_VERIFY_SID"))
	os.Setenv("JWT_SECRET", os.Getenv("JWT_SECRET"))

	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	// tables := []string {"users", "farmers", "buyers", "admins", "auth", "products", "orders"}
	// for i := 0; i < len(tables); i++ {
	// 	if err := db.DropTable(conn, tables[i]); err!= nil {
	// 		panic(err)
	// 	}
	// }

	if err = db.CreateTable(); err != nil {
		log.Printf("error creating table : %v", err)
		panic(err)
	}

	defer conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://krishi-bazar-xi.vercel.app", "*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		ExposeHeaders: []string{"Content-Length"},
		MaxAge: 86400,
	}))


	// e.Use(middleware.Recover())
	api := e.Group("/api")
	// Public routes
	auth := api.Group("/auth")
	auth.POST("/signup", authy.HandleSignUp())
	auth.POST("/complete-signup", authy.HandleCompleteSignup(conn))
	auth.POST("/login", authy.HandleLogin())
	auth.POST("/complete-login", authy.HandleCompleteLogin(conn))

	// Admin routes
	admin := api.Group("/admin")
	admin.POST("/login", admins.AdminLogin(conn))

	adminv1 := admin.Group("/v1")
	adminv1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	adminv1.GET("/dashboard", admins.GetAllUnapprovedFarmers(conn))
	adminv1.GET("/users/:id", admins.GetUserProfile(conn))
	adminv1.POST("/user/:id/approve", admins.ApproveUser(conn))
	adminv1.POST("/approve-product", admins.ApproveProduct(conn))

	// protected routes
	v1 := api.Group("/v1")
	v1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	v1.Use(authy.ExtractUserID)

	// User routes --> store userId locally which is being returned by login
	user := v1.Group("/user")
	user.GET("/:id", users.GetUserProfile(conn)) // -> user/123/ , for req body i have to use post method and directly send the req body
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
	user.GET("/:id/orders", order.GetOrders(conn)) // -> GET ALL ORDERS
	orders.GET("/:id", order.GetOrdersByID(conn))  // -> GET ORDER BY ID
	orders.PUT("/:id/status", order.UpdateOrderStatus(conn))

	e.Logger.Fatal(e.Start(":8080"))
}
