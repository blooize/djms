package main

import (
	"djms-discord-bot/pkg/discord"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = discord.StartBot(os.Getenv("DISCORD_BOT_TOKEN"), os.Getenv("SECRET_TOKEN"))
	if err != nil {
		log.Fatalf("[DISCORD] Error starting Discord bot: %v", err)
	}
	log.Printf("[DISCORD] Discord bot started successfully\n")
}
