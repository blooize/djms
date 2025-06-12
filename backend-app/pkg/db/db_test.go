package db

import (
	"testing"

	"gorm.io/gorm"
)

var connection *gorm.DB

func TestScheme(t *testing.T) {
	connection := InitializeDatabase()
	if connection == nil {
		t.Fatal("Failed to initialize database connection")
	}
}

func TestFindDJByName(t *testing.T) {
	connection = InitializeDatabase()
	dj := FindDJByName(connection, "bloo")

	if dj.ID == 0 {
		t.Errorf("Expected DJ to be found, but got nothing back")
	} else {
		t.Logf("Found DJ: {ID: %d, Name: %s}", dj.ID, dj.Name)
	}
	connection.Delete(&dj)
}

func TestFindDJByID(t *testing.T) {
	connection = InitializeDatabase()
	foo := DJ{Name: "foo"}
	connection.Create(&foo)

	dj, check := FindDJByID(connection, foo.ID)

	if dj.ID == 0 && !check {
		t.Errorf("Expected DJ to be found with ID: %d, but got nothing back", foo.ID)
	} else {
		t.Logf("Created DJ: {ID: %d, Name: %s}", dj.ID, dj.Name)
	}
	connection.Delete(&dj)

	dj, check = FindDJByID(connection, dj.ID)
	if dj.ID != 0 || check {
		t.Errorf("Expected no DJ to be found with ID 999999, but got: {ID: %d, Name: %s}", dj.ID, dj.Name)
	} else {
		t.Logf("No DJ found with ID 999999 as expected")
	}

}

func TestFindVrcdnByDJID(t *testing.T) {
	connection = InitializeDatabase()

	dj := DJ{Name: "bloo"}
	connection.Create(&dj)
	link := VrcdnLink{RTSPT: "foo", DJID: dj.ID}
	connection.Create(&link)

	vrcdn, check := FindVrcdnByDJID(connection, link.DJID)

	if !check && vrcdn.ID == 0 {
		t.Errorf("No vrcdn Found with ID %d as expected but got %d back", link.ID, vrcdn.ID)
	} else {
		t.Logf("Expected vrcdn with ID %d", vrcdn.ID)
	}
	connection.Delete(&vrcdn)
	connection.Delete(&dj)

	vrcdn, check = FindVrcdnByDJID(connection, dj.ID)
	if vrcdn.ID != 0 || check {
		t.Errorf("Expected no vrcdn to be found with DJ ID %d, but got: {ID: %d, RTSPT: %s}", dj.ID, vrcdn.ID, vrcdn.RTSPT)
	} else {
		t.Logf("No vrcdn found with DJ ID %d as expected", dj.ID)
	}
}

func TestFindVrcdnByLink(t *testing.T) {
	connection = InitializeDatabase()

	dj := DJ{Name: "bloo"}
	connection.Create(&dj)
	link := VrcdnLink{RTSPT: "foo", DJID: dj.ID}
	connection.Create(&link)

	vrcdn := FindVrcdnByLink(connection, link.RTSPT, dj)
	if vrcdn.ID == 0 {
		t.Errorf("Expected vrcdn to be found with RTSPT %s, but got nothing back", link.RTSPT)
	} else {
		t.Logf("Found vrcdn: {ID: %d, RTSPT: %s, DJID: %d}", vrcdn.ID, vrcdn.RTSPT, vrcdn.DJID)
	}
	connection.Delete(&vrcdn)
	connection.Delete(&dj)

}

func TestFindClubByID(t *testing.T) {
	connection = InitializeDatabase()

	place := Club{Name: "foo"}
	connection.Create(&place)

	club, check := FindClubByID(connection, place.ID)

	if !check && club.ID == 0 {
		t.Errorf("No Club Found with ID %d as expected", place.ID)
	} else {
		t.Logf("Expected Club with ID %d, and found: {ID: %d, Name: %s}", place.ID, club.ID, club.Name)
	}
	connection.Delete(&club)

	club, check = FindClubByID(connection, place.ID)
	if club.ID != 0 && check {
		t.Errorf("Expected no Club to be found with ID %d, but got: {ID: %d, Name: %s}", place.ID, club.ID, club.Name)
	} else {
		t.Logf("No Club found with ID %d as expected", place.ID)
	}
}

func TestFindEventByID(t *testing.T) {
	connection = InitializeDatabase()

	show := Event{Name: "foo"}
	connection.Create(&show)

	event, check := FindEventByID(connection, show.ID)

	if !check && event.ID == 0 {
		t.Errorf("No Event Found with ID %d", show.ID)
	} else {
		t.Logf("Expected Event with ID %d, and found: {ID: %d, Name: %s}", show.ID, event.ID, event.Name)
	}

	connection.Delete(&event)
	event, check = FindEventByID(connection, show.ID)
	if event.ID != 0 && check {
		t.Errorf("Expected no Event to be found with ID %d, but got: {ID: %d, Name: %s}", show.ID, event.ID, event.Name)
	} else {
		t.Logf("No Event found with ID %d as expected", show.ID)
	}
}

func TestFindUserByID(t *testing.T) {
	connection = InitializeDatabase()
	user := User{Username: "foo"}
	connection.Create(&user)

	foundUser, check := FindUserByID(connection, user.ID)

	if !check && foundUser.ID == 0 {
		t.Errorf("No User Found with ID %d", user.ID)
	} else {
		t.Logf("Expected User with ID %d, and found: {ID: %d, Username: %s}", user.ID, foundUser.ID, foundUser.Username)
	}
	connection.Delete(&foundUser)

	foundUser, check = FindUserByID(connection, user.ID)
	if foundUser.ID != 0 && check {
		t.Errorf("Expected no User to be found with ID %d, but got: {ID: %d, Username: %s}", user.ID, foundUser.ID, foundUser.Username)
	} else {
		t.Logf("No User found with ID %d as expected", user.ID)
	}
}

func TestFindClubModeratorByUserID(t *testing.T) {
	connection = InitializeDatabase()

	user := User{Username: "foo"}
	connection.Create(&user)

	moderator := ClubModerator{UserID: user.ID, ClubID: 1}
	connection.Create(&moderator)

	foundModerator, check := FindClubModeratorByUserID(connection, user.ID)

	if !check && foundModerator.ID == 0 {
		t.Errorf("No Club Moderator Found with User ID %d", user.ID)
	} else {
		t.Logf("Expected Club Moderator with User ID %d, and found: {ID: %d, UserID: %d}", user.ID, foundModerator.ID, foundModerator.UserID)
	}
	connection.Delete(&foundModerator)
	connection.Delete(&user)

	foundModerator, check = FindClubModeratorByUserID(connection, user.ID)
	if foundModerator.ID != 0 && check {
		t.Errorf("Expected no Club Moderator to be found with User ID %d, but got: {ID: %d, UserID: %d}", user.ID, foundModerator.ID, foundModerator.UserID)
	} else {
		t.Logf("No Club Moderator found with User ID %d as expected", user.ID)
	}
}

func TestEventDJsByEventID(t *testing.T) {
	connection = InitializeDatabase()

	event := Event{Name: "foo"}
	connection.Create(&event)
	dj := DJ{Name: "bloo"}
	connection.Create(&dj)
	eventDJ := EventDJ{EventID: event.ID, DJID: dj.ID}
	connection.Create(&eventDJ)

	foundEventDJs, check := FindEventDJsByEventID(connection, event.ID)
	if !check && len(foundEventDJs) == 0 {
		t.Errorf("No Event DJs found for Event ID %d", event.ID)
	} else {
		t.Logf("Found Event DJs for Event ID %d: %v", event.ID, foundEventDJs)
	}

	connection.Delete(&eventDJ)
	connection.Delete(&dj)
	connection.Delete(&event)

	foundEventDJs, check = FindEventDJsByEventID(connection, event.ID)
	if check && len(foundEventDJs) != 0 {
		t.Errorf("Expected no Event DJs to be found for Event ID %d, but got: %v", event.ID, foundEventDJs)
	} else {
		t.Logf("No Event DJs found for Event ID %d as expected", event.ID)
	}
}

func TestFindEventDJByDJID(t *testing.T) {
	connection = InitializeDatabase()

	dj := DJ{Name: "bloo"}
	connection.Create(&dj)
	event := Event{Name: "foo"}
	connection.Create(&event)
	eventDJ := EventDJ{EventID: event.ID, DJID: dj.ID}
	connection.Create(&eventDJ)

	foundEventDJ, check := FindEventDJByDJID(connection, dj.ID)
	if !check && foundEventDJ.ID == 0 {
		t.Errorf("No Event DJ found for DJ ID %d", dj.ID)
	} else {
		t.Logf("Found Event DJ for DJ ID %d: {ID: %d, EventID: %d}", dj.ID, foundEventDJ.ID, foundEventDJ.EventID)
	}

	connection.Delete(&eventDJ)
	connection.Delete(&dj)
	connection.Delete(&event)

	foundEventDJ, check = FindEventDJByDJID(connection, dj.ID)
	if foundEventDJ.ID != 0 && check {
		t.Errorf("Expected no Event DJ to be found for DJ ID %d, but got: {ID: %d, EventID: %d}", dj.ID, foundEventDJ.ID, foundEventDJ.EventID)
	} else {
		t.Logf("No Event DJ found for DJ ID %d as expected", dj.ID)
	}
}
