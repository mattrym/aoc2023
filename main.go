package main

import (
	"aoc2023/day01"
	"aoc2023/day02"
	"aoc2023/day03"
	"aoc2023/day04"
	"aoc2023/day05"
	"aoc2023/day06"
	"fmt"
	"os"
	"strconv"
)

// Run the day based
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide a day number")
		return
	}

	dayNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Please provide a valid day number")
		return
	}

	switch dayNumber {
	case 1: // day01
		day01.Run()
	case 2: // day02
		day02.Run()
	case 3: // day03
		day03.Run()
	case 4: // day04
		day04.Run()
	case 5: // day05
		day05.Run()
	case 6: // day06
		day06.Run()
	}
}
