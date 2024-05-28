# Railway Network Simulation - README.md

This project simulates railway network operations, including loading network maps, exploring paths between stations, and optimizing train schedules.

# Usage

Run the program using the following command:

go run . <network map file> <start station> <end station> <number of trains>
<network map file>: Path to the network map file.
<start station>: Name of the starting station.
<end station>: Name of the ending station.
<number of trains>: Number of trains to be allocated.

# Examples

Running with Different Network Maps

go run . network7.map small large 9
go run . network6.map jungle desert 10
go run . network11.map bond_square space_port 4
go run . network8.map beethoven part 9
go run . network10.map beginning terminus 20
go run . network5.map two four 4
go run . network3.map waterloo st_pancras 2
go run . network2.map waterloo st_pancras 4

# Tests with Faulty Maps

go run . network_err1.map beethoven part 9
go run . network_err2.map beethoven part 9
go run . network_err12.map beethoven part 9

# Command Description

Counting the Output Lines
You can count the number of lines in the output of the program using the wc -l command. This is useful for checking how many lines of output the program generates. For example:


go run . network10.map beginning terminus 20 | wc -l

This command does the following:

go run . network10.map beginning terminus 20: Runs the program with the specified network map file, starting station, ending station, and number of trains.
|: Pipes the output of the program to the next command.
wc -l: Counts the number of lines in the output.
This can be useful for debugging or analyzing the program's output size.

# Error Handling
The program handles various errors such as:

Incorrect number of arguments
Invalid number of trains
Errors in loading network maps
No routes found between stations
Duplicate or invalid station names and coordinates
Error messages are printed with a red background for visibility.