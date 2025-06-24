package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func StartBot(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	// Register the slash command handler
	// dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	if i.Type == discordgo.InteractionApplicationCommand {
	// 		switch i.ApplicationCommandData().Name {
	// 		case "hello":
	// 			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 				Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 				Data: &discordgo.InteractionResponseData{
	// 					Content: "Hello from the bot!",
	// 				},
	// 			})
	// 		}
	// 	}
	// })

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}
	defer dg.Close()

	// _, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", &discordgo.ApplicationCommand{
	// 	Name:        "hello",
	// 	Description: "Say hello!",
	// })
	// if err != nil {
	// 	return fmt.Errorf("cannot create slash command: %w", err)
	// }

	go func() {
		r := gin.Default()
		r.POST("/update", func(c *gin.Context) {
			jwt, err := c.Cookie("jwt")
			if err != nil || jwt == "" {
				c.String(401, "Missing or invalid JWT cookie")
				return
			}
			channelID := "1386866380233510962"
			_, err = dg.ChannelMessageSend(channelID, "Pong from API!")
			if err != nil {
				fmt.Printf("Failed to Update Message: %v\n", err)
				c.String(500, "Failed to update message")
				return
			}
			c.String(200, "Message sent!")
		})
		r.Run(":6969")
	}()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	return nil
}
