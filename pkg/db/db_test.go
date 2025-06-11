package db

import (
	"testing"

	"gorm.io/gorm"
)

var connection *gorm.DB

func TestScheme(t *testing.T) {
	connection := InitializeTestDatabase()
	if connection == nil {
		t.Fatal("Failed to initialize database connection")
	}
}

func TestFindDJByName(t *testing.T) {
	connection = InitializeTestDatabase()
	dj := FindDJByName(connection, "bloo")

	if dj.ID == 0 {
		t.Errorf("Expected DJ to be found, but got nothing back")
	} else {
		t.Logf("Found DJ: {ID: %d, Name: %s}", dj.ID, dj.Name)
	}
}

func TestFindDJByID(t *testing.T) {
	connection = InitializeTestDatabase()
	foo := DJ{Name: "foo"}
	connection.Create(&foo)

	dj, check := FindDJByID(connection, foo.ID)

	if dj.ID == 0 && !check {
		t.Errorf("Expected DJ to be found with ID: %d, but got nothing back", foo.ID)
	} else {
		t.Logf("Created DJ: {ID: %d, Name: %s}", dj.ID, dj.Name)
	}
	connection.Delete(&dj)

}

func TestFindVrcdnByDJID(t *testing.T) {
	connection = InitializeTestDatabase()

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
}

func TestFindVrcdnByLink(t *testing.T) {
	connection = InitializeTestDatabase()

	dj := DJ{Name: "bloo"}
	connection.Create(&dj)
	link := VrcdnLink{RTSPT: "foo", DJID: dj.ID}
	connection.Create(&link)

	vrcdn := FindVrcdnByLink(connection, link.RTSPT, dj)

	if vrcdn.ID == 0 {
		t.Errorf("Expected vrcdn %s to be found, but got nothing back", link.RTSPT)
	} else {
		t.Logf("Found vrcdn: {ID: %d, RTSPT: %s}", vrcdn.ID, vrcdn.RTSPT)
	}

	connection.Delete(&vrcdn)

}

func TestFindClubByID(t *testing.T) {
	connection = InitializeTestDatabase()

	place := Club{Name: "foo"}
	connection.Create(&place)

	club, check := FindClubByID(connection, place.ID)

	if !check && club.ID == 0 {
		t.Errorf("No Club Found with ID %d as expected", place.ID)
	} else {
		t.Logf("Expected Club with ID %d, and found: {ID: %d, Name: %s}", place.ID, club.ID, club.Name)
	}
	connection.Delete(&club)
}

func TestFindEventByID(t *testing.T) {
	connection = InitializeTestDatabase()

	show := Event{Name: "foo"}
	connection.Create(&show)

	event, check := FindEventByID(connection, show.ID)

	if !check && event.ID == 0 {
		t.Errorf("No Event Found with ID %d", show.ID)
	} else {
		t.Logf("Expected Event with ID %d, and found: {ID: %d, Name: %s}", show.ID, event.ID, event.Name)
	}

}
