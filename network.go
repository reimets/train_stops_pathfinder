package main

import (
	"fmt"
)

// RailNetwork represents the railway network
type RailNetwork struct {
	stations map[string]*Location
	links    map[string]map[string]bool
}

// Location represents a station in the network
type Location struct {
	name string
}

// NewRailNetwork initializes a new railway network
func NewRailNetwork() *RailNetwork {
	return &RailNetwork{
		stations: make(map[string]*Location),
		links:    make(map[string]map[string]bool),
	}
}

// AddLocation adds a new station to the network
func (network *RailNetwork) AddLocation(name string) {
	if _, exists := network.stations[name]; !exists {
		network.stations[name] = &Location{name: name}
		network.links[name] = make(map[string]bool)
	}
}

// AddLink adds a bidirectional track between two stations
func (network *RailNetwork) AddLink(start, end string) error {
	if _, exists := network.stations[start]; !exists {
		return fmt.Errorf("station %s does not exist", start)
	}
	if _, exists := network.stations[end]; !exists {
		return fmt.Errorf("station %s does not exist", end)
	}
	if network.links[start][end] || network.links[end][start] {
		return fmt.Errorf("duplicate connection between %s and %s", start, end)
	}
	network.links[start][end] = true
	network.links[end][start] = true
	return nil
}
