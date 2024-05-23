package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadNetworkMap reads and constructs the railway network from the file
func LoadNetworkMap(filename string) (*RailNetwork, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	network := NewRailNetwork()

	isEmpty, stationsSectionFound, connectionsSectionFound := checkSections(scanner)
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

	stationsCount, err := processStationsAndConnections(scanner, network)
	if err != nil {
		return nil, err
	}
	if stationsCount > 10000 {
		return nil, fmt.Errorf("map contains more than 10000 stations")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return network, nil
}

func checkSections(scanner *bufio.Scanner) (bool, bool, bool) {
	isEmpty := true
	stationsSectionFound := false
	connectionsSectionFound := false

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

	return isEmpty, stationsSectionFound, connectionsSectionFound
}

func processStationsAndConnections(scanner *bufio.Scanner, network *RailNetwork) (int, error) {
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
			if err := processStation(line, stations, coordinates, network); err != nil {
				return 0, err
			}
			stationsCount++
		} else if readLinks {
			if err := processLink(line, links, network); err != nil {
				return 0, err
			}
		}
	}

	return stationsCount, nil
}

func processStation(line string, stations, coordinates map[string]string, network *RailNetwork) error {
	parts := strings.Split(line, ",")
	if len(parts) != 3 {
		return fmt.Errorf("station %s does not have correct amount of coordinates", parts[0])
	}
	name := strings.TrimSpace(parts[0])

	xCoord := strings.TrimSpace(parts[1])
	yCoord := strings.TrimSpace(parts[2])
	// Check if the coordinates are numeric and not negative
	if x, err := strconv.Atoi(xCoord); err != nil || x < 0 {
		return fmt.Errorf("station %s has invalid coordinate %s", name, xCoord)
	}
	if y, err := strconv.Atoi(yCoord); err != nil || y < 0 {
		return fmt.Errorf("station %s has invalid coordinate %s", name, yCoord)
	}

	// Check if the station name is unique
	if _, exists := stations[name]; exists {
		return fmt.Errorf("station list has two stations with same name: %s", name)
	}
	stations[name] = line

	// Check if coordinates are unique
	coord := fmt.Sprintf("%s,%s", xCoord, yCoord)
	if _, exists := coordinates[coord]; exists {
		return fmt.Errorf("two or more stations have same coordinates")
	}
	coordinates[coord] = name

	network.AddLocation(name)
	return nil
}

func processLink(line string, links map[string]bool, network *RailNetwork) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("connections section has fault in '%s' row: incorrect amount of stations in row", parts[0])
	}
	from := strings.TrimSpace(parts[0])
	to := strings.TrimSpace(parts[1])
	// check for duplicate connections
	linkKey1 := fmt.Sprintf("%s-%s", from, to)
	linkKey2 := fmt.Sprintf("%s-%s", to, from)
	if links[linkKey1] || links[linkKey2] {
		return fmt.Errorf("duplicate connection between %s and %s", from, to)
	}
	links[linkKey1] = true
	links[linkKey2] = true

	return network.AddLink(from, to)
}
