package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

// Test the ExplorePaths function
func TestExplorePaths_SameInputStations(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"part": {name: "part"},
		},
		links: map[string]map[string]bool{
			"part": {},
		},
	}

	_, err := network.ExplorePaths("part", "part")
	if err == nil {
		t.Fatalf("Test didn't pass. Expected error when source and destination stations are the same, got nil")
	} else if err.Error() != "source and destination stations are the same" {
		t.Fatalf("Test didn't pass. Expected 'source and destination stations are the same', got %s", err.Error())
	}
}

// Test when the source station does not exist
func TestExplorePaths_InputSourceStationDoesNotExists(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"part": {name: "part"},
		},
		links: map[string]map[string]bool{
			"part": {},
		},
	}

	_, err := network.ExplorePaths("beethoveni", "part")
	if err == nil {
		t.Fatalf("Test didn't pass. Expected error when source station does not exist, got nil")
	} else if err.Error() != "source station beethoveni does not exist" {
		t.Fatalf("Test didn't pass. Expected 'source station beethoveni does not exist', got %s", err.Error())
	}
}

// Test when the destination station does not exist
func TestExplorePaths_InputDestinationStationDoesNotExists(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"beethoven": {name: "beethoven"},
		},
		links: map[string]map[string]bool{
			"beethoven": {},
		},
	}

	_, err := network.ExplorePaths("beethoven", "parts")
	if err == nil {
		t.Fatalf("Test didn't pass. Expected error when destination station does not exist, got nil")
	} else if err.Error() != "destination station parts does not exist" {
		t.Fatalf("Test didn't pass. Expected 'destination station parts does not exist', got %s", err.Error())
	}
}

// Test when no routes are found between source and destination
func TestExplorePaths_NoRoutesBetweenSourceAndDestination(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"beethoven": {name: "beethoven"},
			"part":      {name: "part"},
		},
		links: map[string]map[string]bool{
			"beethoven": {},
			"part":      {},
		},
	}

	_, err := network.ExplorePaths("beethoven", "part")
	if err == nil {
		t.Fatalf("Test didn't pass. Expected an error for file where no routes found from start to end.")
	} else if err.Error() != "no routes found from start to end" {
		t.Fatalf("Test didn't pass. Expected 'no routes found from start to end' error, got: %v", err)
	}
}

// Test if source and destination stations are connected directly
func TestExplorePaths_DirectRoute(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"beethoven": {name: "beethoven"},
			"part":      {name: "part"},
		},
		links: map[string]map[string]bool{
			"beethoven": {"part": true},
			"part":      {"beethoven": true},
		},
	}

	routes, err := network.ExplorePaths("beethoven", "part")
	if err != nil {
		t.Fatalf("Test didn't pass. Unexpected error: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("Test didn't pass. Expected 1 route, got %d", len(routes))
	}
	expectedRoute := []string{"beethoven", "part"}
	for i, station := range expectedRoute {
		if routes[0][i] != station {
			t.Fatalf("Test didn't pass. Expected route %v, got %v", expectedRoute, routes[0])
		}
	}
}

// Test if there are multiple routes between the source and destination stations
func TestExplorePaths_MultipleRoutes(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"beethoven": {name: "beethoven"},
			"part":      {name: "part"},
			"mozart":    {name: "mozart"},
		},
		links: map[string]map[string]bool{
			"beethoven": {"part": true, "mozart": true},
			"part":      {"beethoven": true},
			"mozart":    {"beethoven": true, "part": true},
		},
	}

	routes, err := network.ExplorePaths("beethoven", "part")
	if err != nil {
		t.Fatalf("Test didn't pass. Unexpected error: %v", err)
	}
	if len(routes) != 2 {
		t.Fatalf("Test didn't pass. Expected 2 routes, got %d", len(routes))
	}
}

// Test if the route includes intermediate stations
func TestExplorePaths_RouteWithIntermediateStations(t *testing.T) {
	network := RailNetwork{
		stations: map[string]*Location{
			"beethoven": {name: "beethoven"},
			"part":      {name: "part"},
			"mozart":    {name: "mozart"},
		},
		links: map[string]map[string]bool{
			"beethoven": {"mozart": true},
			"mozart":    {"beethoven": true, "part": true},
			"part":      {"mozart": true},
		},
	}

	routes, err := network.ExplorePaths("beethoven", "part")
	if err != nil {
		t.Fatalf("Test didn't pass. Unexpected error: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("Test didn't pass. Expected 1 route, got %d", len(routes))
	}
	expectedRoute := []string{"beethoven", "mozart", "part"}
	for i, station := range expectedRoute {
		if routes[0][i] != station {
			t.Fatalf("Test didn't pass. Expected route %v, got %v", expectedRoute, routes[0])
		}
	}
}

// Test the Contains function
func TestContains(t *testing.T) {
	// Test: element in slice
	slice := []string{"beethoven", "mozart", "bach"}
	item := "mozart"
	result := contains(slice, item)
	if !result {
		t.Fatalf("Test didn't pass. Expected true, got %v for item %v in slice %v", result, item, slice)
	}

	// Test: element is not in slice
	item = "verdi"
	result = contains(slice, item)
	if result {
		t.Fatalf("Test didn't pass. Expected false, got %v for item %v in slice %v", result, item, slice)
	}

	// Test: empty slice
	slice = []string{}
	item = "beethoven"
	result = contains(slice, item)
	if result {
		t.Fatalf("Test didn't pass. Expected false, got %v for item %v in slice %v", result, item, slice)
	}

	// Test: A slice contains the same element multiple times
	slice = []string{"beethoven", "mozart", "mozart", "bach"}
	item = "mozart"
	result = contains(slice, item)
	if !result {
		t.Fatalf("Test didn't pass. Expected true, got %v for item %v in slice %v", result, item, slice)
	}
}

// TestValidateRoutes tests ValidateRoutes function
func TestValidateRoutes(t *testing.T) {
	// Test: one route
	routes := [][]string{
		{"start", "A", "end"},
	}
	expected := [][][]string{
		{{"start", "A", "end"}},
	}
	result := ValidateRoutes(routes)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Test 1 failed. Expected %v, got %v", expected, result)
	}

	// Test Multiple routes that do not overlap
	routes = [][]string{
		{"start", "A", "end"},
		{"start", "B", "end"},
	}
	expected = [][][]string{
		{{"start", "A", "end"}, {"start", "B", "end"}},
	}
	result = ValidateRoutes(routes)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Test 2 failed. Expected %v, got %v", expected, result)
	}

	// Test: Multiple routes that overlap
	routes = [][]string{
		{"start", "A", "end"},
		{"start", "A", "B", "end"},
		{"start", "C", "end"},
	}
	expected = [][][]string{
		{{"start", "A", "end"}, {"start", "C", "end"}},
		{{"start", "C", "end"}, {"start", "A", "B", "end"}},
	}
	result = ValidateRoutes(routes)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Test 3 failed. Expected %v, got %v", expected, result)
	}
}

// TestDisplaySchedule tests DisplaySchedule function
func TestDisplaySchedule(t *testing.T) {
	tests := []struct {
		name        string
		plan        routePlan
		trainCount  int
		expectedOut string
	}{
		{
			name: "One train one route",
			plan: routePlan{
				lengths:           []int{1},
				trainDistribution: []int{1},
				routes:            [][]string{{"start", "A", "end"}},
				totalTurns:        2,
			},
			trainCount:  1,
			expectedOut: "T1-A \nT1-end \n",
		},
		{
			name: "Multiple trains multiple routes",
			plan: routePlan{
				lengths:           []int{1, 1},
				trainDistribution: []int{1, 1},
				routes:            [][]string{{"start", "A", "end"}, {"start", "B", "end"}},
				totalTurns:        2,
			},
			trainCount:  2,
			expectedOut: "T1-A T2-B \nT1-end T2-end \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture the output
			r, w, _ := os.Pipe()
			old := os.Stdout
			os.Stdout = w

			var buf bytes.Buffer
			done := make(chan struct{})
			go func() {
				io.Copy(&buf, r)
				close(done)
			}()

			DisplaySchedule(tt.plan, tt.trainCount)

			w.Close()
			os.Stdout = old
			<-done

			got := buf.String()
			if got != tt.expectedOut {
				t.Errorf("DisplaySchedule() = %v, want %v", got, tt.expectedOut)
			}
		})
	}
}
