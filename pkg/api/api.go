package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"main/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func SetupRouter(client_id string, client_secret string, redirect_uri string, jwt_secret string) *gin.Engine {
	r := gin.Default()
	r.Use(AuthMiddleware(jwt_secret))

	conf := &oauth2.Config{
		Endpoint:     discord.Endpoint,
		Scopes:       []string{discord.ScopeIdentify},
		RedirectURL:  redirect_uri,
		ClientID:     client_id,
		ClientSecret: client_secret,
	}
	// routessss yessir
	r.GET("/", func(c *gin.Context) {
		foo, err := c.Get("userID")
		if !err {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			log.Fatalf("Error getting userID from context: %v", err)
			return
		}

		c.String(200, foo.(string))
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
		// ctx.Data(200, "application/json", body)

		var userData map[string]any

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
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": discordID,
		})
		s, err := t.SignedString([]byte(jwt_secret))

		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			log.Fatalf("Error signing JWT: %v", err)
			return
		}

		ctx.SetCookie("jwt", s, 600000, "/", "localhost", false, true)
		ctx.Redirect(302, "http://localhost:4000/")
		// theres more to do here with authentication/authorization but i want to do less annoying stuff now
	})

	r.GET("/api/clubs", func(ctx *gin.Context) {
		clubs, err := db.FindClubsOwnedByUserID(db.InitializeDatabase(), ctx.MustGet("userID").(string)) //i believe .MustGet essentially forces authorization
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(200, clubs)
	})

	r.GET("/api/club/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		club, found := db.FindClubByID(db.InitializeDatabase(), id)
		if !found {
			ctx.JSON(404, gin.H{"error": "Club not found"})
			return
		}
		ctx.JSON(200, club)
	})

	// creates the club and also adds the user as the owner
	r.POST("/api/club/add", func(ctx *gin.Context) {
		var foo db.Club
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&foo); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		club := db.CreateClub(connection, db.Club{Name: foo.Name})
		user := db.FindUserByDiscordID(connection, ctx.MustGet("userID").(string))
		db.CreateClubOwner(connection, club.ID, user.ID)

		ctx.JSON(201, club)
	})

	return r
}
