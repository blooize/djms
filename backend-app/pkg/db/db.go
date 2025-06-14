package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	DiscordID string
	Avatar    string
}

type Club struct {
	gorm.Model
	Name string
}

type Event struct {
	gorm.Model
	Name   string
	ClubID uint
}

type DJ struct {
	gorm.Model
	Name   string
	UserID uint
}

type VrcdnLink struct {
	gorm.Model
	RTSPT string
	DJID  uint
}

type ClubModerator struct {
	gorm.Model
	ClubID uint
	UserID uint
}

type ClubOwner struct {
	gorm.Model
	ClubID uint
	UserID uint
}

type EventDJ struct {
	gorm.Model
	EventID uint
	DJID    uint
}

func InitializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("djms.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// migrate the schema
	db.AutoMigrate(&User{}, &Club{}, &Event{}, &DJ{}, &VrcdnLink{}, &ClubModerator{}, &ClubOwner{}, &EventDJ{})

	return db
}

// find a DJ by name, if not found, create a new DJ with the given name
// note: not quite sure if collisions actually matter here since name isnt unique
// and one person can have multiple names in practice
// impersonation is possible, but idgaf
func CheckUserExists(db *gorm.DB, username string, discordID string, avatar string) {
	var user User

	db.First(&user, "discord_id = ?", discordID)

	if user.ID == 0 {
		user = User{Username: username, DiscordID: discordID, Avatar: avatar}
		db.Create(&user)
	} else {
		user.Username = username
		user.Avatar = avatar
		db.Save(&user)
	}
}

func FindDJByName(db *gorm.DB, name string) DJ {
	var dj DJ
	db.First(&dj, "name = ?", name)

	if dj.ID == 0 {
		dj = DJ{Name: name}
		db.Create(&dj)
	}
	return dj
}

func FindDJByID(db *gorm.DB, id uint) (DJ, bool) {
	var dj DJ
	db.First(&dj, id)

	if dj.ID == 0 {
		return dj, false
	} else {
		return dj, true
	}
}
func FindVrcdnByDJID(db *gorm.DB, id uint) (VrcdnLink, bool) {
	var vrcdn VrcdnLink
	db.First(&vrcdn, "dj_id = ?", id)

	if vrcdn.ID == 0 {
		return vrcdn, false
	} else {
		return vrcdn, true
	}
}

// if link cant be found for given link and DJ ID, create a new link with the given DJ
// note: 1:n relationship between one DJ and multiple links (i think thats how it worked???)

func FindVrcdnByLink(db *gorm.DB, rtspt string, dj DJ) VrcdnLink {
	var vrcdn VrcdnLink
	db.First(&vrcdn, "RTSPT = ?", rtspt, "DJID = ?", dj.ID)

	if vrcdn.ID == 0 {
		vrcdn = VrcdnLink{RTSPT: rtspt, DJID: dj.ID}
		db.Create(&vrcdn)
	}
	return vrcdn
}

func FindClubByUserID(db *gorm.DB, id uint) (Club, bool) {
	var club Club
	db.First(&club, "user_id = ?", id)

	if club.ID == 0 {
		return club, false
	} else {
		return club, true
	}
}

// please tell me theres a better way to do this
func FindClubsOwnedByUserID(db *gorm.DB, userID uint) ([]Club, error) {
	var clubs []Club
	err := db.Joins("JOIN club_owners ON club_owners.club_id = clubs.id").
		Where("club_owners.user_id = ?", userID).
		Find(&clubs).Error
	return clubs, err
}

func FindEventByID(db *gorm.DB, id uint) (Event, bool) {
	var event Event
	db.First(&event, id)

	if event.ID == 0 {
		return event, false
	} else {
		return event, true
	}
}

func FindEventsByClubID(db *gorm.DB, id uint) []Event {
	var events []Event
	db.Find(&events, "club_id = ?", id)
	return events
}

func FindUserByID(db *gorm.DB, id uint) (User, bool) {
	var user User
	db.First(&user, id)

	if user.ID == 0 {
		return user, false
	} else {
		return user, true
	}
}

func FindClubModeratorByUserID(db *gorm.DB, id uint) (ClubModerator, bool) {
	var moderator ClubModerator
	db.First(&moderator, "user_id = ?", id)

	if moderator.ID == 0 {
		return moderator, false
	} else {
		return moderator, true
	}
}

func FindEventDJsByEventID(db *gorm.DB, id uint) ([]EventDJ, bool) {
	var eventDJs []EventDJ
	db.Find(&eventDJs, "event_id = ?", id)

	if len(eventDJs) == 0 {
		return eventDJs, false
	} else {
		return eventDJs, true
	}
}

func FindEventDJByDJID(db *gorm.DB, id uint) (EventDJ, bool) {
	var eventDJ EventDJ
	db.First(&eventDJ, "dj_id = ?", id)

	if eventDJ.ID == 0 {
		return eventDJ, false
	} else {
		return eventDJ, true
	}
}

func CreateClub(db *gorm.DB, name string) Club {
	club := Club{Name: name}
	db.Create(&club)
	return club
}

func CreateEvent(db *gorm.DB, name string, clubID uint) Event {
	event := Event{Name: name, ClubID: clubID}
	db.Create(&event)
	return event
}

func CreateDJ(db *gorm.DB, name string, userID uint) DJ {
	dj := DJ{Name: name, UserID: userID}
	db.Create(&dj)
	return dj
}

func CreateEventDJ(db *gorm.DB, eventID uint, djID uint) EventDJ {
	eventDJ := EventDJ{EventID: eventID, DJID: djID}
	db.Create(&eventDJ)
	return eventDJ
}

func CreateClubOwner(db *gorm.DB, clubID uint, userID uint) ClubOwner {
	owner := ClubOwner{ClubID: clubID, UserID: userID}
	db.Create(&owner)
	return owner
}

func CreateClubModerator(db *gorm.DB, clubID uint, userID uint) ClubModerator {
	moderator := ClubModerator{ClubID: clubID, UserID: userID}
	db.Create(&moderator)
	return moderator
}
