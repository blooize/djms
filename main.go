package main

import (
	"fmt"
	"log"
	"main/pkg/api"
	"main/pkg/db"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the database
	connection := db.InitializeDatabase()

	fmt.Printf("Database initialized: %v\n", connection != nil)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client_id := os.Getenv("DISCORD_CLIENT_ID")
	client_secret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirect_uri := os.Getenv("DISCORD_REDIRECT_URL")

	r := api.SetupRouter(client_id, client_secret, redirect_uri)
	r.Run(":4000")
}
