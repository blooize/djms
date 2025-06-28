package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func StartAPI(dg *discordgo.Session, backend_token string) {
	r := gin.Default()
	r.POST("/createMessage", func(c *gin.Context) {
		var data struct {
			EventID     uint   `json:"event_id"`
			ChannelID   string `json:"channel_id"`
			SecretToken string `json:"secret_token"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.String(400, "Invalid request payload")
			return
		}
		if data.SecretToken != backend_token {
			c.String(403, "Forbidden")
			return
		}
		_, err := CreateMessage(dg, data.ChannelID)
		if err != nil {
			fmt.Printf("Failed to Create Message: %v\n", err)
			c.String(500, "Failed to create message")
			return
		}
		c.String(200, "Message sent!")
	})

	r.PATCH("/updateMessage", func(c *gin.Context) {
		var data struct {
			EventID     uint   `json:"event_id"`
			ChannelID   string `json:"channel_id"`
			MessageID   string `json:"message_id"`
			SecretToken string `json:"secret_token"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.String(400, "Invalid request payload")
			return
		}
		if data.SecretToken != backend_token {
			c.String(403, "Forbidden")
			return
		}
		err := UpdateMessage(dg, data.ChannelID, data.MessageID)
		if err != nil {
			fmt.Printf("Failed to Refresh Message: %v\n", err)
			c.String(500, "Failed to refresh message")
			return
		}
		c.String(200, "Message refreshed!")
	})

	r.DELETE("/deleteMessage", func(c *gin.Context) {
		var data struct {
			EventID     uint   `json:"event_id"`
			ChannelID   string `json:"channel_id"`
			MessageID   string `json:"message_id"`
			SecretToken string `json:"secret_token"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.String(400, "Invalid request payload")
			return
		}
		if data.SecretToken != backend_token {
			c.String(403, "Forbidden")
			return
		}
		err := DeleteMessage(dg, data.ChannelID, data.MessageID)
		if err != nil {
			fmt.Printf("Failed to Delete Message: %v\n", err)
			c.String(500, "Failed to delete message")
			return
		}
		c.String(200, "Message deleted!")
	})

	r.Run(":6969")
}
