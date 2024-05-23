package main

import (
	"errors"
	"fmt"
	"sort"
)

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

// Helper function to check if a slice contains an item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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
