package main

import (
	"testing"
)

// TestAddLocation tests AddLocation function
func TestAddLocation(t *testing.T) {
	network := &RailNetwork{
		stations: make(map[string]*Location),
		links:    make(map[string]map[string]bool),
	}

	// Test that the station is added to the network if it does not exist
	network.AddLocation("beethoven")
	if _, exists := network.stations["beethoven"]; !exists {
		t.Fatalf("AddLocation failed to add new station 'beethoven'")
	}
	if _, exists := network.links["beethoven"]; !exists {
		t.Fatalf("AddLocation failed to initialize links for new station 'beethoven'")
	}

	// Test that the station is not added again if it already exists
	initialStationsCount := len(network.stations)
	initialLinksCount := len(network.links)
	network.AddLocation("beethoven")
	if len(network.stations) != initialStationsCount {
		t.Fatalf("AddLocation should not add station 'beethoven' again")
	}
	if len(network.links) != initialLinksCount {
		t.Fatalf("AddLocation should not initialize links for station 'beethoven' again")
	}
}

// TestAddLink tests AddLink function
func TestAddLink(t *testing.T) {
	network := &RailNetwork{
		stations: map[string]*Location{
			"beethoven": {name: "beethoven"},
			"mozart":    {name: "mozart"},
		},
		links: map[string]map[string]bool{
			"beethoven": {},
			"mozart":    {},
		},
	}

	// Test that the connection is added successfully
	err := network.AddLink("beethoven", "mozart")
	if err != nil {
		t.Fatalf("AddLink failed to add link: %v", err)
	}
	if !network.links["beethoven"]["mozart"] || !network.links["mozart"]["beethoven"] {
		t.Fatalf("AddLink did not create a bidirectional link between beethoven and mozart")
	}

	// Test that the source station does not exist
	err = network.AddLink("bach", "mozart")
	if err == nil {
		t.Fatalf("Expected an error for non-existent station 'bach'")
	}
	if err.Error() != "station bach does not exist" {
		t.Fatalf("Expected error message 'station bach does not exist', got: %v", err)
	}

	// Test that the destination station does not exist
	err = network.AddLink("beethoven", "bach")
	if err == nil {
		t.Fatalf("Expected an error for non-existent station 'bach'")
	}
	if err.Error() != "station bach does not exist" {
		t.Fatalf("Expected error message 'station bach does not exist', got: %v", err)
	}

	// Test that the connection already exists
	err = network.AddLink("beethoven", "mozart")
	if err == nil {
		t.Fatalf("Expected an error for duplicate connection between 'beethoven' and 'mozart'")
	}
	if err.Error() != "duplicate connection between beethoven and mozart" {
		t.Fatalf("Expected error message 'duplicate connection between beethoven and mozart', got: %v", err)
	}
}
