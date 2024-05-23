//go run . network7.map small large 9
//go run . network6.map jungle desert 10
//go run . network11.map bond_square space_port 4
//go run . network8.map beethoven part 9
// // faulty result // go run . network10.map beginning  terminus 20
//go run . network5.map two four 4
//go run . network1-1.map waterloo st_pancras 2

// tests with faulty maps
//go run . network1.map beethoven part 9
//go run . network6.map jungle desert 10

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	if network.links[start][end] {
		return fmt.Errorf("Error: duplicate connection between %s and %s", start, end)
	}
	network.links[start][end] = true
	network.links[end][start] = true
	return nil
}

// ExplorePaths uses BFS to find all routes from source to destination
func (network *RailNetwork) ExplorePaths(source, destination string) ([][]string, error) {
	if source == destination {
		return nil, errors.New("source and destination stations are the same")
	}
	if _, exists := network.stations[source]; !exists {
		return nil, fmt.Errorf("source station %s does not exist", source)
	}
	if _, exists := network.stations[destination]; !exists {
		return nil, fmt.Errorf("destination station %s does not exist", destination)
	}

	queue := [][]string{{source}}
	var routes [][]string

	for len(queue) > 0 {
		route := queue[0]
		queue = queue[1:]
		current := route[len(route)-1]

		if current == destination {
			routes = append(routes, route)
			continue
		}

		for neighbor := range network.links[current] {
			if !contains(route, neighbor) {
				newRoute := append([]string{}, route...)
				newRoute = append(newRoute, neighbor)
				queue = append(queue, newRoute)
			}
		}
	}

	return routes, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateRoutes filters and returns valid route combinations without overlaps
func ValidateRoutes(routes [][]string) [][][]string {
	var validCombos [][][]string

	for _, route := range routes {
		var occupiedStations []string
		occupiedStations = append(occupiedStations, route[1:len(route)-1]...)
		combo := [][]string{route}

		for _, otherRoute := range routes {
			if &route == &otherRoute {
				continue
			}
			valid := true
			for _, station := range otherRoute[1 : len(otherRoute)-1] {
				if contains(occupiedStations, station) {
					valid = false
					break
				}
			}
			if valid {
				combo = append(combo, otherRoute)
				occupiedStations = append(occupiedStations, otherRoute[1:len(otherRoute)-1]...)
			}
		}

		validCombos = append(validCombos, combo)
	}

	for _, routes := range validCombos {
		sort.SliceStable(routes, func(i, j int) bool {
			return len(routes[i]) < len(routes[j])
		})
	}

	return validCombos
}

// SelectOptimalCombos selects the best combinations of routes
func SelectOptimalCombos(combos [][][]string) [][][]string {
	var maxRoutes int
	for _, combo := range combos {
		if len(combo) > maxRoutes {
			maxRoutes = len(combo)
		}
	}

	var optimalCombos [][][]string
	for i := 1; i <= maxRoutes; i++ {
		for _, combo := range combos {
			if len(combo) >= i {
				optimalCombos = append(optimalCombos, combo[:i])
			}
		}
	}

	sort.SliceStable(optimalCombos, func(a, b int) bool {
		return len(optimalCombos[a][0]) < len(optimalCombos[b][0])
	})

	return optimalCombos
}

// AllocateTrains determines the best routes for the trains to minimize turns
func AllocateTrains(trainCount int, optimalCombos [][][]string) routePlan {
	plans := make([]routePlan, len(optimalCombos))
	for i, combo := range optimalCombos {
		for _, route := range combo {
			plans[i].lengths = append(plans[i].lengths, len(route)-2)
			plans[i].trainDistribution = append(plans[i].trainDistribution, len(route)-2)
		}
		plans[i].routes = combo
	}

	for i := range optimalCombos {
		trainsLeft := trainCount
		for trainsLeft > 0 {
			shortest := plans[i].trainDistribution[0]
			shortestIndex := 0
			for j, trains := range plans[i].trainDistribution {
				if trains < shortest {
					shortest = trains
					shortestIndex = j
				}
			}
			plans[i].trainDistribution[shortestIndex]++
			trainsLeft--
		}
		plans[i].totalTurns = plans[i].trainDistribution[0]
	}

	minTurns := plans[0].totalTurns
	bestPlan := plans[0]
	for _, plan := range plans {
		if plan.totalTurns < minTurns {
			minTurns = plan.totalTurns
			bestPlan = plan
		}
		for j, length := range plan.lengths {
			plan.trainDistribution[j] -= length
		}
	}

	return bestPlan
}

// DisplaySchedule prints the train movements per turn
func DisplaySchedule(plan routePlan, trainCount int) {
	trains := make([]trainStatus, trainCount)
	schedule := make([][]trainStatus, plan.totalTurns)

	for i := 0; i < trainCount; i++ {
		trains[i].id = i + 1
		trains[i].position = 1
	}

	trainIndex := 0
	for turn := 0; turn < plan.totalTurns; turn++ {
		if turn > 0 {
			for _, train := range schedule[turn-1] {
				if train.position != plan.lengths[train.route]+1 {
					train.position++
					train.location = plan.routes[train.route][train.position]
					schedule[turn] = append(schedule[turn], train)
				}
			}
		}

		for routeIdx, trainsLeft := range plan.trainDistribution {
			if trainsLeft > 0 {
				trains[trainIndex].route = routeIdx
				trains[trainIndex].location = plan.routes[routeIdx][trains[trainIndex].position]
				schedule[turn] = append(schedule[turn], trains[trainIndex])
				plan.trainDistribution[routeIdx]--
				trainIndex++
			}
		}
	}

	for _, turn := range schedule {
		for _, train := range turn {
			fmt.Printf("T%d-%s ", train.id, train.location)
		}
		fmt.Println()
	}
}

// trainStatus holds the train's current status
type trainStatus struct {
	id       int
	route    int
	position int
	location string
}

// routePlan holds the planned routes for the trains
type routePlan struct {
	routes            [][]string
	lengths           []int
	trainDistribution []int
	totalTurns        int
}

// LoadNetworkMap reads and constructs the railway network from the file
func LoadNetworkMap(filename string) (*RailNetwork, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	network := NewRailNetwork()
	readStations := false
	readLinks := false
	
	stationsCount := 0
	coordinates := make(map[string]string)
	stations := make(map[string]string)

	isEmpty := true
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, "#")[0]
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		isEmpty = false

		if line == "stations:" {
			readStations = true
			readLinks = false
			continue
		} else if line == "connections:" {
			readStations = false
			readLinks = true
			continue
		}

		if readStations {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("Error loading network map: station %s does not have correct amount of coordinates", parts[0])
			}
			name := strings.TrimSpace(parts[0])
			xCoord := strings.TrimSpace(parts[1])
			yCoord := strings.TrimSpace(parts[2])

			// Kontroll, kas koordinaadid on numbrilised ja mitte negatiivsed
			if x, err := strconv.Atoi(xCoord); err != nil || x < 0 {
				return nil, fmt.Errorf("Error loading network map: station %s has invalid coordinate %s", name, xCoord)
			}
			if y, err := strconv.Atoi(yCoord); err != nil || y < 0 {
				return nil, fmt.Errorf("Error loading network map: station %s has invalid coordinate %s", name, yCoord)
			}

			// Kontroll, kas jaama nimi on kordumatu
			if _, exists := stations[name]; exists {
				return nil, fmt.Errorf("Error loading network map: station list has two stations with same name: %s", name)
			}
			stations[name] = line

			// Kontroll, kas koordinaadid on kordumatud
			coord := fmt.Sprintf("%s,%s", xCoord, yCoord)
			if _, exists := coordinates[coord]; exists {
				return nil, fmt.Errorf("Error loading network map: two or more stations have same coordinates")
			}
			coordinates[coord] = name

			network.AddLocation(name)
			stationsCount++
		} else if readLinks {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				continue
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			err := network.AddLink(from, to)
			if err != nil {
				return nil, err
			}
		}
	}

	if isEmpty {
		return nil, fmt.Errorf("file is empty")
	}

	if !readStations {
		return nil, fmt.Errorf("Error loading network map: \"stations:\" section does not exist")
	}

	if !readLinks {
		return nil, fmt.Errorf("Error loading network map: \"connections:\" section does not exist")
	}

	if stationsCount > 10000 {
		return nil, fmt.Errorf("Error: map contains more than 10000 stations")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return network, nil
}

func main() {
	startTime := time.Now()
	if len(os.Args) != 5 {
		log.Fatal("Usage: go run . <network map file> <start station> <end station> <number of trains>")
	}
	fileName := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	trainCount, err := strconv.Atoi(os.Args[4])
	if err != nil || trainCount <= 0 {
		log.Fatal("Error: Number of trains must be a valid positive integer")
	}

	network, err := LoadNetworkMap(fileName)
	if err != nil {
		log.Fatal("\033[41mError:\033[0m", err)
	}

	routes, err := network.ExplorePaths(startStation, endStation)
	if err != nil {
		log.Fatal("Error exploring paths:", err)
	}
	if len(routes) == 0 {
		log.Fatal("Error: No routes found from start to end")
	}

	validRoutes := ValidateRoutes(routes)
	optimalCombos := SelectOptimalCombos(validRoutes)
	bestPlan := AllocateTrains(trainCount, optimalCombos)
	DisplaySchedule(bestPlan, trainCount)

	fmt.Println("\nProgram executed in:", time.Since(startTime))
}
