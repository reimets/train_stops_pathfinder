// go run . network7.map small large 9
// go run . network6.map jungle desert 10
// go run . network11.map bond_square space_port 4
// go run . network8.map beethoven part 9
// go run . network10.map beginning terminus 20
// go run . network5.map two four 4
// go run . network3.map waterloo st_pancras 2
// go run . network2.map waterloo st_pancras 4

// // tests with faulty maps
// go run . network_err1.map beethoven part 9
// go run . network_err1-1.map beethoven part 9
// go run . network_err2.map beethoven part 9
// go run . network_err3.map beethoven part 9
// go run . network_err4.map beethoven part 9
// go run . network_err5.map beethoven part 9
// go run . network_err6.map beethoven part 9
// go run . network_err7.map beethoven part 9
// go run . network_err8.map beethoven part 9
// go run . network_err9.map beethoven part 9
// go run . network_err10.map beethoven part 9
// go run . network_err11.map beethoven part 9
// go run . network_err12.map beethoven part 9

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
	if network.links[start][end] || network.links[end][start] {
		return fmt.Errorf("duplicate connection between %s and %s", start, end)
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
	if len(routes) == 0 {
		return nil, errors.New("no routes found from start to end")
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
	used := make(map[int]bool) // Hoidke marsruutide indekseid, mida juba kasutati kombodes

	for i, route := range routes {
		if used[i] {
			continue
		}

		var occupiedStations []string
		occupiedStations = append(occupiedStations, route[1:len(route)-1]...)
		combo := [][]string{route}
		used[i] = true

		for j, otherRoute := range routes {
			if i == j {
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
				used[j] = true
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

	stationsSectionFound := false
	connectionsSectionFound := false
	isEmpty := true

	// first control if "stations:" and "connections:" rows exists //
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, "#")[0]
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		isEmpty = false

		if line == "stations:" {
			stationsSectionFound = true
		} else if line == "connections:" {
			connectionsSectionFound = true
		}
	}

	if isEmpty {
		return nil, fmt.Errorf("file is empty")
	}

	if !stationsSectionFound {
		return nil, fmt.Errorf("'stations:' section does not exist")
	}

	if !connectionsSectionFound {
		return nil, fmt.Errorf("'connections:' section does not exist")
	}

	// If the sections exist, we do a second pass to process the contents of the file
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	readStations := false
	readLinks := false

	stationsCount := 0
	coordinates := make(map[string]string)
	stations := make(map[string]string)
	links := make(map[string]bool)

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
		}
		if line == "connections:" {
			readStations = false
			readLinks = true
			continue
		}

		if readStations {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("station %s does not have correct amount of coordinates", parts[0])
			}
			name := strings.TrimSpace(parts[0])

			xCoord := strings.TrimSpace(parts[1])
			yCoord := strings.TrimSpace(parts[2])
			// Check if the coordinates are numeric and not negative
			if x, err := strconv.Atoi(xCoord); err != nil || x < 0 {
				return nil, fmt.Errorf("station %s has invalid coordinate %s", name, xCoord)
			}
			if y, err := strconv.Atoi(yCoord); err != nil || y < 0 {
				return nil, fmt.Errorf("station %s has invalid coordinate %s", name, yCoord)
			}

			// Check if the station name is unique
			if _, exists := stations[name]; exists {
				return nil, fmt.Errorf("station list has two stations with same name: %s", name)
			}
			stations[name] = line

			// Check if coordinates are unique
			coord := fmt.Sprintf("%s,%s", xCoord, yCoord)
			if _, exists := coordinates[coord]; exists {
				return nil, fmt.Errorf("two or more stations have same coordinates")
			}
			coordinates[coord] = name

			network.AddLocation(name)
			stationsCount++

		} else if readLinks {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("connections section has fault in '%s' row: incorrect amount of stations in row", parts[0])
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			// check for doplicate connections
			linkKey1 := fmt.Sprintf("%s-%s", from, to)
			linkKey2 := fmt.Sprintf("%s-%s", to, from)
			if links[linkKey1] || links[linkKey2] {
				return nil, fmt.Errorf("duplicate connection between %s and %s", from, to)
			}
			links[linkKey1] = true
			links[linkKey2] = true

			err := network.AddLink(from, to)
			if err != nil {
				return nil, err
			}
		}
	}

	if stationsCount > 10000 {
		return nil, fmt.Errorf("map contains more than 10000 stations")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return network, nil
}

func main() {
	startTime := time.Now()
	if len(os.Args) != 5 {
		log.Fatal("\033[41m ! Error ! \033[0m \nUsage: go run . <network map file> <start station> <end station> <number of trains>")
	}
	fileName := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	trainCount, err := strconv.Atoi(os.Args[4])
	if err != nil || trainCount <= 0 {
		log.Fatal("\033[41m ! Error ! \033[0m Number of trains must be a valid positive integer")
	}

	network, err := LoadNetworkMap(fileName)
	if err != nil {
		log.Fatal("\033[41m ! Error ! \033[0m Error loading network map:", err)
	}

	routes, err := network.ExplorePaths(startStation, endStation)
	if err != nil {
		log.Fatal("\033[41m ! Error ! \033[0m Error exploring paths:", err)
	}

	validRoutes := ValidateRoutes(routes)
	optimalCombos := SelectOptimalCombos(validRoutes)
	bestPlan := AllocateTrains(trainCount, optimalCombos)
	DisplaySchedule(bestPlan, trainCount)

	// those two lines (just different versions) are the ones that cause the number to be bigger when you add "wc -l" in the terminal
	// fmt.Println("\nProgram executed in: ", time.Since(startTime))
	// fmt.Println("\n\033[100m Program executed in: \033[0m ", time.Since(startTime))

	// 
	// fmt.Fprintln(os.Stderr, "\nProgram executed in: ", time.Since(startTime))
	fmt.Fprintln(os.Stderr, "\n\033[100m Program executed in: \033[0m", time.Since(startTime))
}
