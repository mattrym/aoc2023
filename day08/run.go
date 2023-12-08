package day08

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

const INPUT_FILE_PATH = "day08/input.txt"
const NETWORK_REGEX = `^(\w+) = \((\w+), (\w+)\)$`

type Navigation struct {
	LeftRights string
	Network    map[string]Node
}

type Node struct {
	Left  string
	Right string
}

func Run() {
	inputLines := readLines(INPUT_FILE_PATH)
	navigation, err := parseInput(inputLines)
	if err != nil {
		panic(err)
	}

	stepsPart1 := findWayOutPart1(navigation)
	fmt.Println("Steps to exit [PART 1]: ", stepsPart1)

	stepsPart2 := findWayOutPart2(navigation)
	fmt.Println("Steps to exit [PART 2]: ", stepsPart2)
}

func findWayOutPart2(navigation *Navigation) int {
	locations := make(map[string]int)

	for location := range navigation.Network {
		if location[2] == 'A' {
			locations[location] = findWayOut(navigation, location,
				func(location string) bool {
					return location[2] == 'Z'
				})
		}
	}

	stepsToGoOutOfAll := 1
	for _, stepsToGoOut := range locations {
		stepsToGoOutOfAll = leastCommonMultiple(stepsToGoOut, stepsToGoOutOfAll)
	}

	return stepsToGoOutOfAll
}

func greatestCommonDivisor(a int, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func leastCommonMultiple(a int, b int) int {
	return a * b / greatestCommonDivisor(a, b)
}

func findWayOutPart1(navigation *Navigation) int {
	return findWayOut(navigation, "AAA", func(location string) bool {
		return location == "ZZZ"
	})
}

func findWayOut(navigation *Navigation, start string, ended func(string) bool) int {
	location := start
	steps := 0

	for {
		for steps = 0; ; steps++ {
			if ended(location) {
				return steps
			}

			index := steps % len(navigation.LeftRights)
			switch navigation.LeftRights[index] {
			case 'L':
				location = navigation.Network[location].Left
			case 'R':
				location = navigation.Network[location].Right
			}
		}
	}
}

func parseInput(lines []string) (*Navigation, error) {
	var navigation Navigation
	navigation.Network = make(map[string]Node)
	navigation.LeftRights = lines[0]

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		lineRegex := regexp.MustCompile(NETWORK_REGEX)
		matches := lineRegex.FindStringSubmatch(line)

		if len(matches) != 4 {
			return nil, errors.New("invalid input")
		}

		navigation.Network[matches[1]] = Node{Left: matches[2], Right: matches[3]}
	}

	return &navigation, nil
}

func readLines(path string) []string {
	// Read input file (line by line) from input.txt
	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	var lines []string

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return lines
}
