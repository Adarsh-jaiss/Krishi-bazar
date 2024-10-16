package main

import (
	"fmt"
	"log"

	"github.com/adarsh-jaiss/agrohub/db"
)

func main() {

	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	tables := []string{"users", "farmers", "buyers", "admins", "auth", "products", "orders"}
	for i := 0; i < len(tables); i++ {
		if err := db.DropTable(conn, tables[i]); err != nil {
			fmt.Println("error dropping tables...")
			return
		}
	}

	if err = db.CreateTable(); err != nil {
		log.Printf("error creating table : %v", err)
		panic(err)
	}

}
