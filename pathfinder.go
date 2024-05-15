package main

import (
    "fmt"
    "strings"
)



// FindPaths uses BFS to find a minimal set of distinct paths and allows their reuse for multiple trains.
func FindPaths(stations map[string]Station, connections map[string][]string, start, end string, numTrains int) ([]string, error) {
    if _, exists := stations[start]; !exists {
        return nil, fmt.Errorf("Error: start station '%s' does not exist", start)
    }
    if _, exists := stations[end]; !exists {
        return nil, fmt.Errorf("Error: end station '%s' does not exist", end)
    }
    if start == end {
        return nil, fmt.Errorf("Error: start and end station are the same")
    }

    type Path struct {
        route   []string
        visited map[string]bool
    }

    queue := []Path{{route: []string{start}, visited: map[string]bool{start: true}}}
    var pathsFound []string

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        lastStation := current.route[len(current.route)-1]

        if lastStation == end {
            foundPath := strings.Join(current.route, " -> ")
            pathsFound = append(pathsFound, foundPath)
            if len(pathsFound) >= numTrains { // Stop if we have as many paths as trains
                break
            }
            continue
        }

        for _, neighbor := range connections[lastStation] {
            if !current.visited[neighbor] {
                newVisited := make(map[string]bool)
                for k, v := range current.visited {
                    newVisited[k] = v
                }
                newVisited[neighbor] = true
                newPath := Path{
                    route: append([]string(nil), append(current.route, neighbor)...),
                    visited: newVisited,
                }
                queue = append(queue, newPath)
            }
        }
    }

    // If we found fewer distinct paths than needed, reuse the found paths
    if len(pathsFound) < numTrains {
        extraPathsNeeded := numTrains - len(pathsFound)
        for i := 0; i < extraPathsNeeded; i++ {
            pathsFound = append(pathsFound, pathsFound[i%len(pathsFound)])
        }
    }

    return pathsFound, nil
}
