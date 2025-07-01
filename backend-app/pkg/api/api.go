package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"main/pkg/db"
	"strconv"

	"github.com/gin-contrib/cors"
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

func SetupRouter(client_id string, client_secret string, redirect_uri string, jwt_secret string, secret_token string) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000", "http://djms.praxis.club"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
	}))

	public := r.Group("/")
	private := r.Group("/")
	private.Use(AuthMiddleware(jwt_secret))

	conf := &oauth2.Config{
		Endpoint:     discord.Endpoint,
		Scopes:       []string{discord.ScopeIdentify},
		RedirectURL:  redirect_uri,
		ClientID:     client_id,
		ClientSecret: client_secret,
	}
	private.GET("/me", func(ctx *gin.Context) {
		discordID, exists := ctx.Get("discordID")
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		userData := db.GetUserByDiscordID(db.InitializeDatabase(), discordID.(string))
		ctx.JSON(200, gin.H{"user_id": discordID, "username": userData.Username, "avatar": userData.Avatar})
	})

	// routessss yessir
	public.GET("/auth/discord/login", func(ctx *gin.Context) {
		url := "https://discord.com/oauth2/authorize?client_id=1382418824237813911&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%3A4000%2Fauth%2Fdiscord%2Fcallback&scope=identify"
		ctx.Redirect(302, url)
	})

	public.GET("/auth/discord/callback", func(ctx *gin.Context) {
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
			log.Printf("Error unmarshalling user data: %v", err)
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
			log.Printf("Error signing JWT: %v", err)
			return
		}

		ctx.SetCookie("jwt", s, 600000, "/", "localhost", false, false)
		ctx.Redirect(302, "http://localhost:3000/")
		// theres more to do here with authentication/authorization but i want to do less annoying stuff now
	})

	private.GET("/api/clubs", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()

		type Club struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		var res struct {
			Clubs []Club `json:"clubs"`
		}

		user_id := ctx.MustGet("discordID").(string)
		userID := db.GetUserByDiscordID(connection, user_id)

		clubs, err := db.GetClubsOwnedByUserID(connection, strconv.Itoa(int(userID.ID))) //i believe .MustGet essentially forces authorization
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		for _, club := range clubs {
			res.Clubs = append(res.Clubs, Club{
				ID:   club.ID,
				Name: club.Name,
			})
		}

		ctx.JSON(200, res)
	})

	private.GET("/api/club", func(ctx *gin.Context) {
		var data struct {
			ClubID uint `json:"club_id"`
		}
		var res struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		club, found := db.GetClub(connection, data.ClubID)
		if !found {
			ctx.JSON(404, gin.H{"error": "Club not found"})
			return
		}

		res.ID = club.ID
		res.Name = club.Name
		ctx.JSON(200, res)
	})

	private.GET("/api/djslots", func(ctx *gin.Context) {
		var data struct {
			EventID uint `json:"event_id"`
		}

		type Talent struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		type TalentSlot struct {
			Date    uint64   `json:"date"`
			Talents []Talent `json:"talents"`
		}
		var res struct {
			EventID uint         `json:"event_id"`
			ClubID  uint         `json:"club_id"`
			Slots   []TalentSlot `json:"slots"`
		}
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

		slots := db.GetSlotsByEventID(connection, event.ID)
		if len(slots) == 0 {
			ctx.JSON(404, gin.H{"error": "No slots found for this event"})
			return
		}
		var talentSlot []TalentSlot

		for _, slot := range slots {
			talentsData := db.GetSlotTalents(connection, slot.ID)
			var talents []Talent
			for _, talent := range talentsData {
				talentData, found := db.GetTalent(connection, talent.TalentID)
				if found {
					talents = append(talents, Talent{
						ID:   talentData.ID,
						Name: talentData.Name,
					})
				}
			}
			talentSlot = append(talentSlot, TalentSlot{
				Date:    slot.Date,
				Talents: talents,
			})
		}

		res.EventID = event.ID
		res.ClubID = event.ClubID
		res.Slots = talentSlot

		ctx.JSON(200, res)
	})

	// this fucking sucks and would be better if i just used a join but im too dumb for that right now
	private.GET("/api/event/dancerslots", func(ctx *gin.Context) {
		var data struct {
			EventID uint `json:"event_id"`
		}

		type Dancer struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		type Slots struct {
			Date    uint64   `json:"date"`
			Dancers []Dancer `json:"dancers"`
		}

		var res struct {
			ClubID      uint    `json:"club_id"`
			EventID     uint    `json:"event_id"`
			DancerSlots []Slots `json:"dancer_slots"`
		}

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
		dancerslots := db.GetDancerSlots(connection, data.EventID)

		for _, slot := range dancerslots {
			slottalents := db.GetDancerSlotTalents(connection, slot.ID)
			var dancers []Dancer

			for _, talent := range slottalents {
				dancer, found := db.GetDancer(connection, talent.DancerID)
				if !found {
					ctx.JSON(500, gin.H{"error": "Internal Server Error"})
					return
				}
				dancers = append(dancers, Dancer{ID: dancer.ID, Name: dancer.Name})
			}

			res.DancerSlots = append(res.DancerSlots, Slots{
				Date:    slot.Date,
				Dancers: dancers,
			})
		}

		res.EventID = event.ID
		res.ClubID = event.ClubID

		ctx.JSON(200, res)
	})

	private.GET("/api/club/event", func(ctx *gin.Context) {
		var data struct {
			EventID uint `json:"event_id"`
		}

		var res struct {
			ID          uint            `json:"id"`
			Name        string          `json:"name"`
			ClubID      uint            `json:"club_id"`
			TalentSlots []db.Slot       `json:"talent_slots"`
			DancerSlots []db.DancerSlot `json:"dancer_slots"`
		}

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
		res.ID = event.ID
		res.Name = event.Name
		res.ClubID = event.ClubID

		res.TalentSlots = db.GetSlotsByEventID(connection, event.ID)
		res.DancerSlots = db.GetDancerSlots(connection, event.ID)
		ctx.JSON(200, res)
	})

	private.GET("/api/club/events", func(ctx *gin.Context) {

		type Event struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		var res struct {
			Events []Event `json:"events"`
		}

		connection := db.InitializeDatabase()

		clubID, err := strconv.Atoi(ctx.Query("club_id"))
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid club ID"})
			return
		}

		events := db.GetEventsByClubID(connection, uint(clubID))
		for _, event := range events {
			event, found := db.GetEvent(connection, event.ID)
			if !found {
				ctx.JSON(404, gin.H{"error": "Event not found"})
				return
			}
			res.Events = append(res.Events, Event{
				ID:   event.ID,
				Name: event.Name,
			})
		}

		ctx.JSON(200, res)
	})

	private.POST("/api/club", func(ctx *gin.Context) {
		var foo db.Club
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&foo); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		club := db.CreateClub(connection, db.Club{Name: foo.Name})
		user := db.GetUserByDiscordID(connection, ctx.MustGet("discordID").(string))
		db.CreateClubOwner(connection, club.ID, user.ID)

		ctx.JSON(201, club)
	})
	// discord bot only
	private.GET("/api/signupform", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()

		var data struct {
			EventID     uint `json:"event_id"`
			secretToken string
		}

		var res struct {
			MessageID string `json:"message_id"`
			ChannelID string `json:"channel_id"`
			GuildID   string `json:"guild_id"`
			EventID   uint   `json:"event_id"`
			ClubID    uint   `json:"club_id"`
		}

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		if data.secretToken != secret_token {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		event, found := db.GetEvent(connection, data.EventID)
		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		form, err := db.GetSignUpForm(connection, event.ID)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Sign up form not found"})
			return
		}
		res.MessageID = form.MessageID
		res.ChannelID = form.ChannelID
		res.GuildID = form.GuildID
		res.EventID = event.ID
		res.ClubID = event.ClubID
		ctx.JSON(200, res)
	})
	// discord bot only
	private.POST("/api/signupform", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()

		var data struct {
			MessageID   string `json:"message_id"`
			ChannelID   string `json:"channel_id"`
			GuildID     string `json:"guild_id"`
			EventID     uint   `json:"event_id"`
			ClubID      uint   `json:"club_id"`
			SecretToken string `json:"secret_token"`
		}

		var res struct {
			MessageID string `json:"message_id"`
			ChannelID string `json:"channel_id"`
			GuildID   string `json:"guild_id"`
			EventID   uint   `json:"event_id"`
			ClubID    uint   `json:"club_id"`
		}

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		if data.SecretToken != secret_token {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		event, found := db.GetEvent(connection, data.EventID)
		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		form := db.CreateSignUpForm(connection, data.MessageID, data.ChannelID, data.GuildID, event.ID, event.ClubID)
		if form.ID == 0 {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		res.MessageID = form.MessageID
		res.ChannelID = form.ChannelID
		res.GuildID = form.GuildID
		res.EventID = event.ID
		res.ClubID = event.ClubID
		ctx.JSON(201, res)
	})

	private.POST("/api/club/event", func(ctx *gin.Context) {
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

	private.POST("/api/talent", func(ctx *gin.Context) {
		var data struct {
			Name string `json:"name"`
		}

		var res struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		connection := db.InitializeDatabase()

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}
		talent := db.CreateTalent(connection, data.Name)
		res.ID = talent.ID
		res.Name = talent.Name
		ctx.JSON(201, res)

	})

	private.POST("/api/event/slot", func(ctx *gin.Context) {
		var data struct {
			EventID    uint     `json:"event_id"`
			Date       uint64   `json:"date"`
			TalentName []string `json:"talent_names"`
		}
		type Talent struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		var res struct {
			EventID uint     `json:"event_id"`
			Date    uint64   `json:"date"`
			SlotID  uint     `json:"slot_id"`
			Talents []Talent `json:"talents"`
		}

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

		slot := db.CreateSlot(connection, event.ID, data.Date)
		var talents []Talent
		for _, talentName := range data.TalentName {
			talentData := db.GetTalentByName(connection, talentName)
			if talentData.ID == 0 {
				talentData = db.CreateTalent(connection, talentName)
				if talentData.ID == 0 {
					ctx.JSON(500, gin.H{"error": "Internal Server Error"})
					return
				}
			}
			talents = append(talents, Talent{ID: talentData.ID, Name: talentData.Name})

			talent := db.CreateTalentSlot(connection, slot.ID, talentData.ID)
			if talent.ID == 0 {
				ctx.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
		}
		res.Talents = talents
		res.EventID = event.ID
		res.Date = data.Date
		res.SlotID = slot.ID

		ctx.JSON(201, res)
	})

	private.POST("/api/event/dancerslot", func(ctx *gin.Context) {

		connection := db.InitializeDatabase()

		var data struct {
			EventID     uint     `json:"event_id"`
			Date        uint64   `json:"date"`
			DancerNames []string `json:"dancer_names"`
		}

		type Dancer struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		var res struct {
			EventID uint     `json:"event_id"`
			Date    uint64   `json:"date"`
			SlotID  uint     `json:"slot_id"`
			Dancers []Dancer `json:"dancers"`
		}

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

		slot := db.CreateDancerSlot(connection, event.ID, data.Date)
		var dancers []Dancer

		for _, dancer := range data.DancerNames {
			dancerData, exists := db.GetDancerByName(connection, dancer)
			if !exists {
				dancerData = db.CreateDancer(connection, dancer)
				if dancerData.ID == 0 {
					ctx.JSON(500, gin.H{"error": "Internal Server Error"})
					return
				}
			}
			dancers = append(dancers, Dancer{ID: dancerData.ID, Name: dancerData.Name})

			talent := db.CreateDancerSlotTalent(connection, slot.ID, dancerData.ID)
			if talent.ID == 0 {
				ctx.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
		}
		res.Dancers = dancers
		res.EventID = event.ID
		res.Date = data.Date
		res.SlotID = slot.ID

		ctx.JSON(201, res)
	})

	private.POST("/api/dancer", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()

		var data struct {
			Name string `json:"name"`
		}

		var res struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		dancer := db.CreateDancer(connection, data.Name)
		res.ID = dancer.ID
		res.Name = dancer.Name
		ctx.JSON(201, res)
	})

	private.POST("/api/club/moderator", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data struct {
			ClubID uint `json:"club_id"`
			UserID uint `json:"user_id"`
		}

		var res struct {
			ID     uint `json:"id"`
			ClubID uint `json:"club_id"`
			UserID uint `json:"user_id"`
		}

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		//auth check to be sure
		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))
		authorized := db.CheckUserIsOwnerOfClub(connection, data.ClubID, user.ID)
		if !authorized {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		moderator := db.CreateClubModerator(connection, data.ClubID, data.UserID)
		res.ID = moderator.ID
		res.ClubID = moderator.ClubID
		res.UserID = moderator.UserID
		ctx.JSON(201, res)
	})

	private.PUT("/api/event/slot", func(ctx *gin.Context) {
		type Talent struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		var data struct {
			EventID     uint     `json:"event_id"`
			Date        uint64   `json:"date"`
			TalentNames []string `json:"talent_name"`
		}

		var res struct {
			EventID uint     `json:"event_id"`
			Date    uint64   `json:"date"`
			SlotID  uint     `json:"slot_id"`
			Talents []Talent `json:"talents"`
		}
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
		var talents []Talent

		for _, talent := range data.TalentNames {
			talentData := db.GetTalentByName(connection, talent)
			if talentData.ID == 0 {
				talentData = db.CreateTalent(connection, talent)
				if talentData.ID == 0 {
					ctx.JSON(500, gin.H{"error": "Internal Server Error"})
					return
				}
			}
			talents = append(talents, Talent{ID: talentData.ID, Name: talentData.Name})
			db.DeleteSlotTalents(connection, slot.ID)
			_ = db.CreateTalentSlot(connection, slot.ID, talentData.ID)

		}

		res.Talents = talents
		res.EventID = event.ID
		res.Date = data.Date
		res.SlotID = slot.ID

		ctx.JSON(202, res)

	})

	private.PUT("/api/event/dancerslot", func(ctx *gin.Context) {
		type Dancer struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}

		connection := db.InitializeDatabase()
		var data struct {
			EventID     uint     `json:"event_id"`
			Date        uint64   `json:"date"`
			DancerNames []string `json:"dancer_id"`
		}
		var res struct {
			EventID uint     `json:"event_id"`
			Date    uint64   `json:"date"`
			SlotID  uint     `json:"slot_id"`
			Dancers []Dancer `json:"dancers"`
		}

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

		db.DeleteDancerSlotTalents(connection, slot.ID)

		var dancers []Dancer

		for _, dancer := range data.DancerNames {
			dancerData, exists := db.GetDancerByName(connection, dancer)
			if !exists {
				dancerData = db.CreateDancer(connection, dancer)
				if dancerData.ID == 0 {
					ctx.JSON(500, gin.H{"error": "Internal Server Error"})
					return
				}
			}
			dancers = append(dancers, Dancer{ID: dancerData.ID, Name: dancerData.Name})

			talent := db.CreateDancerSlotTalent(connection, slot.ID, dancerData.ID)
			if talent.ID == 0 {
				ctx.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
		}
		res.Dancers = dancers
		res.EventID = event.ID
		res.Date = data.Date
		res.SlotID = slot.ID

		ctx.JSON(202, res)
	})

	private.DELETE("/api/club/moderator", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data struct {
			ID     uint `json:"id"`
			ClubID uint `json:"club_id"`
		}
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

		db.DeleteModerator(connection, data.ClubID, user.ID)
		ctx.JSON(200, gin.H{"message": "Moderator deleted"})
	})

	private.DELETE("/api/event", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data struct {
			ID uint `json:"id"`
		}

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}

		event, found := db.GetEvent(connection, data.ID)
		if !found {
			ctx.JSON(404, gin.H{"error": "Event not found"})
			return
		}

		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))
		authorized := db.CheckUserIsOwnerOfClub(connection, event.ClubID, user.ID)

		if !authorized {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		db.DeleteEvent(connection, data.ID)

		ctx.JSON(200, gin.H{"message": "Event deleted"})
	})

	private.DELETE("/api/club", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data struct {
			ID uint `json:"id"`
		}
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Bad Request"})
			log.Printf("Error binding JSON: %v", err)
			return
		}
		user := db.GetUserByDiscordID(connection, ctx.MustGet("userID").(string))
		authorized := db.CheckUserIsOwnerOfClub(connection, data.ID, user.ID)
		if !authorized {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}
		db.DeleteClub(connection, data.ID)
		ctx.JSON(200, gin.H{"message": "Club deleted"})
	})

	private.DELETE("/api/signupform", func(ctx *gin.Context) {
		connection := db.InitializeDatabase()
		var data struct {
			EventID     uint   `json:"event_id"`
			SecretToken string `json:"secret_token"`
		}

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

		if data.SecretToken != secret_token {
			ctx.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		if err := db.DeleteSignUpForm(connection, event.ID); err != nil {
			ctx.JSON(404, gin.H{"error": "Sign up form not found"})
			return
		}
		ctx.JSON(200, gin.H{"message": "Sign-up form deleted"})
	})

	return r
}
