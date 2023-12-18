package main

import (
	"aoc2023/day01"
	"aoc2023/day02"
	"aoc2023/day03"
	"aoc2023/day04"
	"aoc2023/day05"
	"aoc2023/day06"
	"aoc2023/day07"
	"aoc2023/day08"
	"aoc2023/day09"
	"aoc2023/day10"
	"aoc2023/day11"
	"aoc2023/day12"
	"aoc2023/day13"
	"aoc2023/day14"
	"aoc2023/day15"
	"aoc2023/day16"
	"aoc2023/day17"
	"aoc2023/day18"
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
	case 7: // day07
		day07.Run()
	case 8: // day08
		day08.Run()
	case 9: // day09
		day09.Run()
	case 10: // day10
		day10.Run()
	case 11: // day11
		day11.Run()
	case 12: // day12
		day12.Run()
	case 13: // day13
		day13.Run()
	case 14: // day14
		day14.Run()
	case 15: // day15
		day15.Run()
	case 16: // day16
		day16.Run()
	case 17: // day17
		day17.Run()
	case 18: // day18
		day18.Run()
	}
}
