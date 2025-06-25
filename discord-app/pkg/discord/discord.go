package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func StartBot(token string) error {
	dg, err := discordgo.New("Bot " + token)
	go API(dg)

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

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	return nil
}
