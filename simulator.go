package main

import (
    "fmt"
    "strings"
)

func Simulate(paths []string, numTrains int) error {
    if len(paths) == 0 {
        return fmt.Errorf("no paths available for simulation")
    }

    // Preparing data structures
    trainPaths := make([][]string, numTrains)
    trainPositions := make([]int, numTrains)
    trainActive := make([]bool, numTrains)
    stationOccupancy := make(map[string]int) // Map to track station occupancy

    for i := range trainPaths {
        pathIndex := i % len(paths) // Assign paths in a round-robin fashion
        trainPaths[i] = strings.Split(paths[pathIndex], " -> ")
        stationOccupancy[trainPaths[i][0]]++ // Initially mark the start station as occupied
    }

    allFinished := false
    for !allFinished {
        allFinished = true
        movements := []string{}

        for i := 0; i < numTrains; i++ {
            if trainPositions[i] >= len(trainPaths[i])-1 {
                continue // Skip if train has finished its path
            }

            currentStation := trainPaths[i][trainPositions[i]]
            nextStation := trainPaths[i][trainPositions[i]+1]

            if stationOccupancy[nextStation] == 0 || nextStation == trainPaths[i][len(trainPaths[i])-1] {
                // Move train if next station is free or it is the final station
                trainPositions[i]++
                stationOccupancy[currentStation]--   // Free current station
                stationOccupancy[nextStation]++     // Occupy next station
                movements = append(movements, fmt.Sprintf("T%d-%s", i+1, nextStation))
                if trainPositions[i] < len(trainPaths[i])-1 {
                    allFinished = false // Continue simulation if any train hasn't finished
                }
            } else {
                // Keep train active for next check if it couldn't move
                trainActive[i] = true
            }
        }

        if len(movements) > 0 {
            fmt.Println(strings.Join(movements, " "))
        }
    }

    return nil
}
