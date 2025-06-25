package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func API(dg *discordgo.Session) {
	r := gin.Default()
	r.POST("/update/talent", func(c *gin.Context) {
		var data struct {
			EventID uint `json:"event_id"`
		}
		// Todo: need to connect to the db to get the messageID of the event signup
		// will probably be multiple messages and need to get just talent/dance slots
		if err := c.ShouldBindJSON(&data); err != nil {
			c.String(400, "Invalid request payload")
			return
		}
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

	r.POST("/update/dancer", func(c *gin.Context) {

	})
	r.Run(":6969")
}
