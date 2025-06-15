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
	Name string
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

type Slot struct {
	gorm.Model
	EventID uint
	DJID    uint
	Date    uint64 // using uint64 to store time.Time as Unix timestamp, idk why go doesnt allow me idk
}

func InitializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("djms.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// migrate the schema
	db.AutoMigrate(&User{}, &Club{}, &Event{}, &DJ{}, &VrcdnLink{}, &ClubModerator{}, &ClubOwner{}, &Slot{})

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

// please tell me theres a better way to do this, UPDATE:  yea, because this does not f-ing work
func FindClubsOwnedByUserID(db *gorm.DB, userID string) ([]Club, error) {
	var clubs []Club
	err := db.Joins("JOIN club_owners ON club_owners.club_id = clubs.id").
		Where("club_owners.user_id = ?", userID).
		Find(&clubs).Error
	return clubs, err
}

func FindClubByID(db *gorm.DB, id string) (Club, bool) {
	var club Club
	db.First(&club, "id = ?", id)

	if club.ID == 0 {
		return club, false
	} else {
		return club, true
	}
}

func FindEventByID(db *gorm.DB, id string) (Event, bool) {
	var event Event
	db.First(&event, "id = ?", id)

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

func FindUserByDiscordID(db *gorm.DB, discordID string) User {
	var user User
	db.First(&user, "discord_id = ?", discordID)

	if user.ID == 0 {
		user = User{DiscordID: discordID}
		db.Create(&user)
	}
	return user
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

func CheckIfSlotExists(db *gorm.DB, event_id uint, date uint64) (Slot, bool) {
	var slot Slot
	db.First(&slot, "event_id = ? AND date = ?", event_id, date)

	if slot.ID == 0 {
		return slot, false
	} else {
		return slot, true
	}
}

func CreateClub(db *gorm.DB, club Club) Club {
	db.Create(&club)
	return club
}

func CreateEvent(db *gorm.DB, name string, clubID uint) Event {
	event := Event{Name: name, ClubID: clubID}
	db.Create(&event)
	return event
}

func CreateDJ(db *gorm.DB, name string) DJ {
	dj := DJ{Name: name}
	db.Create(&dj)
	return dj
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

func CreateSlot(db *gorm.DB, eventID uint, djID uint, date uint64) Slot {
	slot := Slot{EventID: eventID, DJID: djID, Date: date}
	db.Create(&slot)
	return slot
}

func UpdateSlot(db *gorm.DB, slot Slot) Slot {
	db.Save(&slot)
	return slot
}
