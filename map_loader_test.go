package main

import (
	"testing"
)

// testing file where two stations are with same name - albinoni
func TestLoadNetworkMap_StationIsMissing(t *testing.T) {
	filePath := "network_err1-1.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where two stations are with same name.")
	} else if err.Error() != "station list has two stations with same name: albinoni" {
		t.Fatalf("Test didn't pass. Expected 'station list has two stations with same name' error, got: %v", err)
	}
}

// testing: no handel station in stations list
func TestLoadNetworkMap_StationsWithSameName(t *testing.T) {
	filePath := "network_err1.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where some station in stations list is missing.")
	} else if err.Error() != "station handel does not exist" {
		t.Fatalf("Test didn't pass. Expected 'station handel does not exist' error, got: %v", err)
	}
}

// testing where two or more stations have same coordinates
func TestLoadNetworkMap_StationsWithSameCoordinates(t *testing.T) {
	filePath := "network_err3.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where two or more stations have same coordinates.")
	} else if err.Error() != "two or more stations have same coordinates" {
		t.Fatalf("Test didn't pass. Expected 'two or more stations have same coordinates' error, got: %v", err)
	}
}

// testing file when one station has faulty coordinates: "handel,3"
func TestLoadNetworkMap_InvalidStationFormat(t *testing.T) {
	filePath := "network_err4.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file with invalid format.")
	} else if err.Error() != "station handel does not have correct amount of coordinates" {
		t.Fatalf("Test didn't pass. Expected 'station handel does not have correct amount of coordinates' error, got: %v", err)
	}
}

// testing file where "stations:" row is missing
func TestLoadNetworkMap_NoStatusRow(t *testing.T) {
	filePath := "network_err5.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file with no 'stations:' row.")
	} else if err.Error() != "'stations:' section does not exist" {
		t.Fatalf("Test didn't pass. Expected ''stations:' section does not exist' error, got: %v", err)
	}
}

// testing file where "connections:" row is missing
func TestLoadNetworkMap_NoConnectionRow(t *testing.T) {
	filePath := "network_err6.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file with no 'connections:' row.")
	} else if err.Error() != "'connections:' section does not exist" {
		t.Fatalf("Test didn't pass. Expected ''connections:' section does not exist' error, got: %v", err)
	}
}

// testing file where duplicate connections exist between handel-mozart
func TestLoadNetworkMap_DuplicateConnections(t *testing.T) {
	filePath := "network_err7.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where duplicate connections exist (between handel-mozart).")
	} else if err.Error() != "duplicate connection between handel and mozart" {
		t.Fatalf("Test didn't pass. Expected 'duplicate connection between handel and mozart' error, got: %v", err)
	}
}

// testing empty file
func TestLoadNetworkMap_EmptyFile(t *testing.T) {
	filePath := "network_err8.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Function didn't return error for empty file.")
	} else if err.Error() != "file is empty" {
		t.Fatalf("Test didn't pass. Expected 'file is empty' error, got: %v", err)
	}
}

// testing file where file has more than 10000 stations
func TestLoadNetworkMap_TooLongMap(t *testing.T) {
	filePath := "network_err9.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where map contains more than 10000 stations.")
	} else if err.Error() != "map contains more than 10000 stations" {
		t.Fatalf("Test didn't pass. Expected 'map contains more than 10000 stations' error, got: %v", err)
	}
}

// testing file where duplicate connections in reverse exist between handel-mozart
func TestLoadNetworkMap_DuplicateConnectionsReverse(t *testing.T) {
	filePath := "network_err10.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where duplicate connections exist (between handel-mozart in reverse).")
	} else if err.Error() != "duplicate connection between mozart and handel" {
		t.Fatalf("Test didn't pass. Expected 'duplicate connection between mozart and handel' error, got: %v", err)
	}
}

// testing file when one station has faulty, negative, coordinates: "albinoni,1,-1"
func TestLoadNetworkMap_InvalidStationCoordinates(t *testing.T) {
	filePath := "network_err11.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file with negative station coordinate.")
	} else if err.Error() != "station albinoni has invalid coordinate -1" {
		t.Fatalf("Test didn't pass. Expected 'station albinoni has invalid coordinate -1' error, got: %v", err)
	}
}

// testing file when one row in "connections:" section is faulty (only one station, no connections to another station)
func TestLoadNetworkMap_InvalidConnectionRow(t *testing.T) {
	filePath := "network_err12.map"
	_, err := LoadNetworkMap(filePath)
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where connections section has fault in 'handel' row: incorrect amount of stations in row.")
	} else if err.Error() != "connections section has fault in 'handel' row: incorrect amount of stations in row" {
		t.Fatalf("Test didn't pass. Expected 'connections section has fault in 'handel' row: incorrect amount of stations in row' error, got: %v", err)
	}
}
