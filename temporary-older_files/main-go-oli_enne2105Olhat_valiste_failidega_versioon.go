package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Error: Incorrect number of command line arguments")
		os.Exit(1)
	}

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of trains must be a valid positive integer")
		os.Exit(1)
	}

	stations, connections, err := LoadNetworkMap(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	paths, err := FindPaths(stations, connections, startStation, endStation, numTrains)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	err = Simulate(paths, numTrains)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func LoadNetworkMap(filePath string) (map[string]Station, map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stations := make(map[string]Station)
	connections := make(map[string][]string)
	parsingStations := true // Toggle to switch between parsing stations and connections

	isEmpty := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		isEmpty = false

		if line == "stations:" {
			parsingStations = true
			continue
		} else if line == "connections:" {
			parsingStations = false
			continue
		}

		if parsingStations {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, nil, fmt.Errorf("invalid station format")
			}
			name := strings.TrimSpace(parts[0])
			stations[name] = Station{Name: name}
		} else {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("invalid connection format")
			}
			src := strings.TrimSpace(parts[0])
			dst := strings.TrimSpace(parts[1])
			connections[src] = append(connections[src], dst)
			connections[dst] = append(connections[dst], src)
		}
	}

	if isEmpty {
		return nil, nil, fmt.Errorf("file is empty")
	}

	return stations, connections, nil
}
