package main

import (
	"log"

	"github.com/adarsh-jaiss/agrohub/db"
	admins "github.com/adarsh-jaiss/agrohub/internal/admin"
	"github.com/adarsh-jaiss/agrohub/internal/auth"
	"github.com/adarsh-jaiss/agrohub/internal/product"
	users "github.com/adarsh-jaiss/agrohub/internal/user"
	"github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	// tables := []string{"users","farmers","buyers","admins","auth"}
	// for i := 0; i < len(tables); i++ {
	// 	if err := db.DropTable(conn,tables[i]); err!= nil{
	// 		fmt.Println("error dropping tables...")
	// 		return
	// 	}
	// }

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

	// // Protected routes
	// jwtMiddleware := middleware.JWT(middleware.JWTConfig{
	// 	SigningKey: []byte(os.Getenv("JWT_SECRET")),
	// })

	// Admin routes
	// admin := api.Group("/admin", jwtMiddleware, auth.IsAdmin)
	admin := api.Group("/admin")
	admin.POST("/approve", admins.ApproveUser(conn))

	// protected routes
	v1 := api.Group("/v1", )
	v1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
	}))

	// User routes --> store userId locally which is being returned by login 
	user := v1.Group("/user")
	user.GET(":id/profile", users.GetUserProfile(conn))		// -> user/123/profile , for req body i have to use post method and directly send the req body
	user.PUT(":id/profile", users.UpdateProfile(conn))

	// Product routes
	products := v1.Group("/products")
	products.GET("", product.ListAllProducts(conn))
	products.GET("/jari", product.ListJariProducts(conn))
	products.GET("/mushroom", product.ListMushroomProducts(conn))
	products.GET("/:id", product.GetProduct(conn))
	products.POST("", product.CreateProduct(conn), authy.IsFarmer)
	products.PUT("/:id", product.UpdateProduct(conn), authy.IsFarmer)
	products.DELETE("/:id", product.DeleteProduct(conn), authy.IsFarmer)

	// // Order routes
	// orders := v1.Group("/orders", jwtMiddleware)
	// orders.POST("", users.CreateOrder(conn))
	// orders.GET("", users.ListOrders(conn))
	// orders.GET("/:id", users.GetOrder(conn))

	e.Logger.Fatal(e.Start(":8080"))
}

func Hello(c echo.Context) error {
	msg := map[string]string{"msg": "hello bhaiya"}
	return c.JSON(200, msg)
}
