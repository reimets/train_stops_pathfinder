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
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

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
