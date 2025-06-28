package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func StartBot(token string, backend_token string) error {
	dg, err := discordgo.New("Bot " + token)
	go StartAPI(dg, backend_token)

	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}
	defer dg.Close()
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	return nil
}

func CreateMessage(dg *discordgo.Session, channelID string) (string, error) {
	message, err := dg.ChannelMessageSend(channelID, "bla")
	if err != nil {
		return "", fmt.Errorf("error sending message: %w", err)
	}
	fmt.Printf("Message sent with ID: %s\n", message.ID)
	return message.ID, nil
}

func RefreshMessage(dg *discordgo.Session, channelID, messageID string) error {
	// This function is a placeholder for refreshing a message.
	// Discord API does not support refreshing messages directly.
	// You might want to edit the message or delete and resend it.
	_, err := dg.ChannelMessageEdit(channelID, messageID, "Updated message content")
	if err != nil {
		return fmt.Errorf("error editing message: %w", err)
	}
	fmt.Printf("Message with ID %s has been refreshed.\n", messageID)
	return nil
}
