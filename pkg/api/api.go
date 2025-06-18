package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"main/pkg/db"
	"strconv"

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
		connection := db.InitializeDatabase()

		user_id := ctx.MustGet("userID").(string)
		userID := db.GetUserByDiscordID(connection, user_id)

		clubs, err := db.GetClubsOwnedByUserID(connection, strconv.Itoa(int(userID.ID))) //i believe .MustGet essentially forces authorization
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.JSON(200, clubs)
	})

	r.GET("/api/club", func(ctx *gin.Context) {
		var data db.Club
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		club, found := db.GetClub(connection, data.ID)
		if !found {
			ctx.JSON(404, gin.H{"error": "Club not found"})
			return
		}
		ctx.JSON(200, club)
	})

	r.GET("/api/djslots", func(ctx *gin.Context) {
		var data struct {
			ID uint `json:"id"`
		}
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		slots := db.GetSlotsByEventID(connection, data.ID)

		ctx.JSON(200, slots)
	})

	r.GET("/api/dancerslots", func(ctx *gin.Context) {
		var data struct {
			ID uint `json:"id"`
		}
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		slots := db.GetDancerSlotsByEventID(connection, data.ID)

		ctx.JSON(200, slots)
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
		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))
		db.CreateClubOwner(connection, club.ID, user.ID)

		ctx.JSON(201, club)
	})

	r.POST("/api/club/event", func(ctx *gin.Context) {
		var foo db.Event
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&foo); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		club, found := db.GetClub(connection, foo.ClubID)

		if !found {
			ctx.JSON(404, gin.H{"error": "Club not found"})
			return
		}

		event := db.CreateEvent(connection, foo.Name, club.ID)
		ctx.JSON(201, event)
	})
	r.POST("/api/dj", func(ctx *gin.Context) {
		var foo db.DJ
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&foo); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}
		dj := db.CreateDJ(connection, foo.Name)
		ctx.JSON(201, dj)
	})

	r.POST("/api/event/slot", func(ctx *gin.Context) {
		var data db.Slot
		connection := db.InitializeDatabase()
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		event, found := db.GetEvent(connection, data.EventID)

		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		_, exists := db.GetSlot(connection, event.ID, data.Date)

		if exists {
			ctx.JSON(500, gin.H{"error": "Slot already exists"})
			return
		}

		slot := db.CreateSlot(connection, event.ID, data.DJID, data.Date)
		ctx.JSON(201, slot)
	})

	r.POST("/api/event/dancerslot", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data db.DancerSlot

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		event, found := db.GetEvent(connection, data.EventID)

		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}
		_, exists := db.GetDancerSlot(connection, event.ID, data.Date)

		if exists {
			ctx.JSON(500, gin.H{"error": "Slot already exists"})
			return
		}

		slot := db.CreateDancerSlot(connection, event.ID, data.DancerID, data.Date)

		ctx.JSON(201, slot)
	})

	r.POST("/api/dancer", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data db.Dancer

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		dancer := db.CreateDancer(connection, data.Name)
		ctx.JSON(201, dancer)
	})

	r.POST("/api/club/moderator", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data db.ClubModerator

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))

		authorized := db.CheckUserIsOwnerOfClub(connection, data.ClubID, user.ID)
		if !authorized {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		moderator := db.CreateClubModerator(connection, data.ClubID, data.UserID)
		ctx.JSON(201, moderator)
	})

	r.PUT("/api/event/slot", func(ctx *gin.Context) {
		var data db.Slot
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}
		event, found := db.GetEvent(connection, data.EventID)

		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		slot, exists := db.GetSlot(connection, event.ID, data.Date)
		if !exists {
			ctx.JSON(404, gin.H{"error": "Slot not found"})
			return

		}

		slot.DJID = data.DJID
		slot.Date = data.Date
		slot = db.UpdateSlot(connection, slot)
		ctx.JSON(202, slot)

	})

	r.PUT("/api/event/dancerslot", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data db.DancerSlot

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		event, found := db.GetEvent(connection, data.EventID)
		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}
		slot, exists := db.GetDancerSlot(connection, event.ID, data.Date)

		if !exists {
			ctx.JSON(404, gin.H{"error": "Dancer slot not found"})
			return
		}

		slot.DancerID = data.DancerID
		slot.Date = data.Date
		slot = db.UpdateDancerSlot(connection, slot)
		ctx.JSON(202, slot)
	})

	r.DELETE("/api/club/moderator", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data db.ClubModerator
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}
		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))
		authorized := db.CheckUserIsOwnerOfClub(connection, data.ClubID, user.ID)

		if !authorized {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		db.DeleteModerator(connection, data.ClubID, data.UserID)
		ctx.JSON(200, gin.H{"message": "Moderator deleted"})
	})

	r.DELETE("/api/event", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data db.Event
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))
		authorized := db.CheckUserIsOwnerOfClub(connection, data.ClubID, user.ID)

		if !authorized {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		db.DeleteEvent(connection, data.ID)
		ctx.JSON(200, gin.H{"message": "Event deleted"})
	})

	return r
}
