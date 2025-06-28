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
		log.Fatal("Error loading .env.dev file")
	}
	discord_bot_token := os.Getenv("DISCORD_BOT_TOKEN")
	backend_token := os.Getenv("SECRET_TOKEN")

	err = discord.StartBot(discord_bot_token, backend_token)
	if err != nil {
		log.Fatalf("[DISCORD] Error starting Discord bot: %v", err)
	}
	log.Printf("[DISCORD] Discord bot started successfully\n")
}
