package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"main/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

type DiscordRequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

func SetupRouter(client_id string, client_secret string, redirect_uri string) *gin.Engine {
	r := gin.Default()
	conf := &oauth2.Config{
		Endpoint:     discord.Endpoint,
		Scopes:       []string{discord.ScopeIdentify},
		RedirectURL:  redirect_uri,
		ClientID:     client_id,
		ClientSecret: client_secret,
	}
	// routessss yessir
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/auth/discord/login", func(ctx *gin.Context) {
		url := "https://discord.com/oauth2/authorize?client_id=1382418824237813911&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%3A4000%2Fauth%2Fdiscord%2Fcallback&scope=identify"
		ctx.Redirect(302, url)
	})

	r.GET("/auth/discord/callback", func(ctx *gin.Context) {
		code := ctx.Query("code")
		if code == "" {
			ctx.JSON(400, gin.H{"error": "Internal Server Error"})
			return
		}
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.Data(200, "application/json", body)

		var userData map[string]interface{}

		err = json.Unmarshal(body, &userData)

		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			log.Fatalf("Error unmarshalling user data: %v", err)
			return
		}

		discordID := userData["id"].(string)
		username := userData["username"].(string)
		avatar := userData["avatar"].(string)

		connection := db.InitializeDatabase()
		db.CheckUserExists(connection, username, discordID, avatar)

	})

	return r
}
