package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
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

type Talent struct {
	gorm.Model
	Name string
}

type VrcdnLink struct {
	gorm.Model
	RTSPT    string
	TalentID uint
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
	Date    uint64 // using uint64 to store time.Time as Unix timestamp, idk why go doesnt allow me idk
}

type SlotTalent struct {
	gorm.Model
	SlotID   uint
	TalentID uint
}

type Dancer struct {
	gorm.Model
	Name string
}

type DancerSlot struct {
	gorm.Model
	EventID uint
	Date    uint64
}

type DancerSlotTalent struct {
	gorm.Model
	DancerSlotID uint
	DancerID     uint
}

type SignUpForm struct {
	gorm.Model
	MessageID string
	ChannelID string
	GuildID   string
	ClubID    uint
	EventID   uint
}

func InitializeDatabase() *gorm.DB {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")

	dsn := "host=" + postgres_host + " user=" + postgres_user + " password=" + postgres_password + " dbname=" + postgres_db + " port=" + postgres_port + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// migrate the schema
	db.AutoMigrate(
		&User{},
		&Club{},
		&Event{},
		&Talent{},
		&VrcdnLink{},
		&ClubModerator{},
		&ClubOwner{},
		&Slot{},
		&Dancer{},
		&DancerSlot{},
		&SlotTalent{},
		&DancerSlotTalent{},
		&SignUpForm{},
	)

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

func GetTalentByName(db *gorm.DB, name string) Talent {
	var talent Talent
	db.First(&talent, "name = ?", name)

	return talent
}

func GetTalent(db *gorm.DB, id uint) (Talent, bool) {
	var talent Talent
	db.First(&talent, id)

	if talent.ID == 0 {
		return talent, false
	} else {
		return talent, true
	}
}
func FindVrcdnByTalentID(db *gorm.DB, id uint) (VrcdnLink, bool) {
	var vrcdn VrcdnLink
	db.First(&vrcdn, "talent_id = ?", id)

	if vrcdn.ID == 0 {
		return vrcdn, false
	} else {
		return vrcdn, true
	}
}

// if link cant be found for given link and DJ ID, create a new link with the given DJ
// note: 1:n relationship between one DJ and multiple links (i think thats how it worked???)

func GetVrcdnByLink(db *gorm.DB, rtspt string, talent Talent) VrcdnLink {
	var vrcdn VrcdnLink
	db.First(&vrcdn, "RTSPT = ?", rtspt, "TalentID = ?", talent.ID)

	if vrcdn.ID == 0 {
		vrcdn = VrcdnLink{RTSPT: rtspt, TalentID: talent.ID}
		db.Create(&vrcdn)
	}
	return vrcdn
}

func GetClubByUserID(db *gorm.DB, id uint) (Club, bool) {
	var club Club
	db.First(&club, "user_id = ?", id)

	if club.ID == 0 {
		return club, false
	} else {
		return club, true
	}
}

func CheckUserIsOwnerOfClub(db *gorm.DB, userID uint, clubID uint) bool {
	var owner ClubOwner
	db.First(&owner, "user_id = ? AND club_id = ?", userID, clubID)

	if owner.ID == 0 {
		return false
	} else {
		return true
	}
}

// please tell me theres a better way to do this, UPDATE:  yea, because this does not f-ing work
func GetClubsOwnedByUserID(db *gorm.DB, userID string) ([]Club, error) {
	var clubs []Club
	err := db.Joins("JOIN club_owners ON club_owners.club_id = clubs.id").
		Where("club_owners.user_id = ?", userID).
		Find(&clubs).Error
	return clubs, err
}

func GetClub(db *gorm.DB, id uint) (Club, bool) {
	var club Club
	db.First(&club, "id = ?", id)

	if club.ID == 0 {
		return club, false
	} else {
		return club, true
	}
}

func GetEvent(db *gorm.DB, id uint) (Event, bool) {
	var event Event
	db.First(&event, "id = ?", id)

	if event.ID == 0 {
		return event, false
	} else {
		return event, true
	}
}

func GetEventsByClubID(db *gorm.DB, id uint) []Event {
	var events []Event
	db.Find(&events, "club_id = ?", id)
	return events
}

func GetUser(db *gorm.DB, id uint) (User, bool) {
	var user User
	db.First(&user, id)

	if user.ID == 0 {
		return user, false
	} else {
		return user, true
	}
}

func GetUserByDiscordID(db *gorm.DB, discordID string) User {
	var user User
	db.First(&user, "discord_id = ?", discordID)

	if user.ID == 0 {
		user = User{DiscordID: discordID}
		db.Create(&user)
	}
	return user
}

func GetClubModeratorByUserID(db *gorm.DB, id uint) (ClubModerator, bool) {
	var moderator ClubModerator
	db.First(&moderator, "user_id = ?", id)

	if moderator.ID == 0 {
		return moderator, false
	} else {
		return moderator, true
	}
}

func GetSlot(db *gorm.DB, event_id uint, date uint64) (Slot, bool) {
	var slot Slot
	db.First(&slot, "event_id = ? AND date = ?", event_id, date)

	if slot.ID == 0 {
		return slot, false
	} else {
		return slot, true
	}
}

func GetSlotsByEventID(db *gorm.DB, eventID uint) []Slot {
	var slots []Slot
	db.Find(&slots, "event_id = ?", eventID)
	return slots
}

func GetSlotTalents(db *gorm.DB, slotID uint) []SlotTalent {
	var slotTalents []SlotTalent
	db.Find(&slotTalents, "slot_id = ?", slotID)
	return slotTalents
}

func GetDancerSlots(db *gorm.DB, eventID uint) []DancerSlot {
	var dancerSlots []DancerSlot
	db.Find(&dancerSlots, "event_id = ?", eventID)
	return dancerSlots
}

func GetDancerSlot(db *gorm.DB, event_id uint, date uint64) (DancerSlot, bool) {
	var dancerSlot DancerSlot
	db.First(&dancerSlot, "event_id = ? AND date = ?", event_id, date)

	if dancerSlot.ID == 0 {
		return dancerSlot, false
	} else {
		return dancerSlot, true
	}
}

func GetDancer(db *gorm.DB, id uint) (Dancer, bool) {
	var dancer Dancer
	db.First(&dancer, id)

	if dancer.ID == 0 {
		return dancer, false
	} else {
		return dancer, true
	}
}

func GetDancerByName(db *gorm.DB, name string) (Dancer, bool) {
	var dancer Dancer
	db.First(&dancer, "name = ?", name)

	if dancer.ID == 0 {
		return dancer, false
	} else {
		return dancer, true
	}
}

func GetDancerSlotTalents(db *gorm.DB, DancerSlotID uint) []DancerSlotTalent {
	var dancerTalents []DancerSlotTalent
	db.Find(&dancerTalents, "dancer_slot_id = ?", DancerSlotID)
	return dancerTalents
}

func GetSignUpForm(db *gorm.DB, eventID uint) (SignUpForm, error) {
	var signUpForm SignUpForm
	db.First(&signUpForm, "event_id = ?", eventID)

	if signUpForm.ID == 0 {
		return signUpForm, gorm.ErrRecordNotFound
	}
	return signUpForm, nil
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

func CreateTalent(db *gorm.DB, name string) Talent {
	talent := Talent{Name: name}
	db.Create(&talent)
	return talent
}

func CreateDancer(db *gorm.DB, name string) Dancer {
	dancer := Dancer{Name: name}
	db.Create(&dancer)
	return dancer
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

func CreateSlot(db *gorm.DB, eventID uint, date uint64) Slot {
	slot := Slot{EventID: eventID, Date: date}
	db.Create(&slot)
	return slot
}

func CreateDancerSlot(db *gorm.DB, eventID uint, date uint64) DancerSlot {
	dancerSlot := DancerSlot{EventID: eventID, Date: date}
	db.Create(&dancerSlot)
	return dancerSlot
}

func CreateDancerSlotTalent(db *gorm.DB, dancerSlotID uint, dancerID uint) DancerSlotTalent {
	dancerSlotTalent := DancerSlotTalent{DancerSlotID: dancerSlotID, DancerID: dancerID}
	db.Create(&dancerSlotTalent)
	return dancerSlotTalent
}

func CreateTalentSlot(db *gorm.DB, slotID uint, talentID uint) SlotTalent {
	slotTalent := SlotTalent{SlotID: slotID, TalentID: talentID}
	db.Create(&slotTalent)
	return slotTalent
}

func CreateSignUpForm(db *gorm.DB, messageID string, channelID string, guildID string, clubID uint, eventID uint) SignUpForm {
	signUpForm := SignUpForm{
		MessageID: messageID,
		ChannelID: channelID,
		GuildID:   guildID,
		ClubID:    clubID,
		EventID:   eventID,
	}
	db.Create(&signUpForm)
	return signUpForm
}

func UpdateSlot(db *gorm.DB, slot Slot) Slot {
	db.Save(&slot)
	return slot
}

func UpdateDancerSlot(db *gorm.DB, dancerSlot DancerSlot) DancerSlot {
	db.Save(&dancerSlot)
	return dancerSlot
}

func DeleteModerator(db *gorm.DB, clubID uint, userID uint) error {
	var moderator ClubModerator
	db.First(&moderator, "club_id = ? AND user_id = ?", clubID, userID)
	if moderator.ID != 0 {
		db.Delete(&moderator)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteEvent(db *gorm.DB, eventID uint) error {
	var event Event
	db.First(&event, eventID)
	if event.ID != 0 {
		db.Delete(&event)
		// really need delete all slots, dancers, djs but cba rn
		var slots []Slot
		db.Find(&slots, "event_id = ?", eventID)
		for _, slot := range slots {
			DeleteSlotTalents(db, slot.ID)
			db.Delete(&slot)
		}
		var dancerSlots []DancerSlot
		db.Find(&dancerSlots, "event_id = ?", eventID)
		for _, dancerSlot := range dancerSlots {
			DeleteDancerSlotTalents(db, dancerSlot.ID)
			db.Delete(&dancerSlot)
		}
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteDancerSlotTalents(db *gorm.DB, slot uint) error {
	var dancerSlotTalents []DancerSlotTalent
	db.Find(&dancerSlotTalents, "dancer_slot_id = ?", slot)

	if len(dancerSlotTalents) > 0 {
		for _, dancerSlotTalent := range dancerSlotTalents {
			db.Delete(&dancerSlotTalent)
		}
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteSlotTalents(db *gorm.DB, slot uint) error {
	var slotTalents []SlotTalent
	db.Find(&slotTalents, "slot_id = ?", slot)

	if len(slotTalents) > 0 {
		for _, slotTalent := range slotTalents {
			db.Delete(&slotTalent)
		}
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteDancerSlot(db *gorm.DB, dancerSlotID uint) error {
	var dancerSlot DancerSlot
	db.First(&dancerSlot, dancerSlotID)
	if dancerSlot.ID != 0 {
		DeleteDancerSlotTalents(db, dancerSlotID)
		db.Delete(&dancerSlot)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteSlot(db *gorm.DB, slotID uint) error {
	var slot Slot
	db.First(&slot, slotID)
	if slot.ID != 0 {
		DeleteSlotTalents(db, slotID)
		db.Delete(&slot)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteClub(db *gorm.DB, clubID uint) error {
	var club Club
	db.First(&club, clubID)
	if club.ID != 0 {
		// delete all events, slots, dancers, djs, moderators and owners
		var events []Event
		db.Find(&events, "club_id = ?", clubID)
		for _, event := range events {
			DeleteEvent(db, event.ID)
		}

		var moderators []ClubModerator
		db.Find(&moderators, "club_id = ?", clubID)
		for _, moderator := range moderators {
			db.Delete(&moderator)
		}

		var owners []ClubOwner
		db.Find(&owners, "club_id = ?", clubID)
		for _, owner := range owners {
			db.Delete(&owner)
		}

		db.Delete(&club)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func DeleteSignUpForm(db *gorm.DB, eventID uint) error {
	var signUpForm SignUpForm
	db.First(&signUpForm, "event_id = ?", eventID)
	if signUpForm.ID != 0 {
		db.Delete(&signUpForm)
		return nil
	}
	return gorm.ErrRecordNotFound
}
