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

	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatalf("Error loading .env.dev file: %v", err)
	}

	client_id := os.Getenv("DISCORD_CLIENT_ID")
	client_secret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirect_uri := os.Getenv("DISCORD_REDIRECT_URL")
	jwt_secret := os.Getenv("JWT_SECRET")
	secret_token := os.Getenv("SECRET_TOKEN")
	r := api.SetupRouter(client_id, client_secret, redirect_uri, jwt_secret, secret_token)
	r.Run(":4000")
}
