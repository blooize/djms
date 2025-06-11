package main

import (
	"fmt"
	"main/pkg/db"
)

func main() {
	// Initialize the database
	connection := db.InitializeDatabase()
	// Create test data
	dj := db.DJ{Name: "bloo"}
	connection.Create(&dj)

	fmt.Printf("%d\n", dj.ID)
}
