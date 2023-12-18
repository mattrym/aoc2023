package day18

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const INPUT_FILE_PATH = "day18/input.txt"

type Direction int

const (
	NONE  Direction = 0
	UP    Direction = 1
	RIGHT Direction = 2
	DOWN  Direction = 4
	LEFT  Direction = 8
)

type Edge struct {
	Direction Direction
	Length    int
}

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	edges, err := parseInput(lines, parseLinePart1)
	if err != nil {
		panic(err)
	}

	areaSize := CalculateArea(edges)
	fmt.Println("Area size [PART 1]: ", areaSize)

	edges, err = parseInput(lines, parseLinePart2)
	if err != nil {
		panic(err)
	}

	areaSize = CalculateArea(edges)
	fmt.Println("Area size [PART 2]: ", areaSize)
}

func CalculateArea(edges []Edge) int {
	_, offsetX, _, offsetY := AccessMaxDims(edges)
	positionX, positionY := offsetX, offsetY
	area := 0
	edgeLengths := 0

	// use shoelace formula to calculate area
	for _, edge := range edges {
		switch edge.Direction {
		case RIGHT:
			area -= positionY * edge.Length
			positionX += edge.Length
		case LEFT:
			area += positionY * edge.Length
			positionX -= edge.Length
		case DOWN:
			positionY += edge.Length
		case UP:
			positionY -= edge.Length
		}
		edgeLengths += edge.Length
	}

	return area + edgeLengths/2 + 1
}

func AccessMaxDims(edges []Edge) (int, int, int, int) {
	maxDimX, maxDimY := 0, 0
	minDimX, minDimY := 0, 0
	dimX, dimY := 0, 0

	for _, edge := range edges {
		switch edge.Direction {
		case RIGHT:
			dimX += edge.Length
		case LEFT:
			dimX -= edge.Length
		case DOWN:
			dimY += edge.Length
		case UP:
			dimY -= edge.Length
		}

		if dimX > maxDimX {
			maxDimX = dimX
		}
		if dimX < minDimX {
			minDimX = dimX
		}
		if dimY < minDimY {
			minDimY = dimY
		}
		if dimY > maxDimY {
			maxDimY = dimY
		}
	}

	return maxDimX - minDimX + 1, -minDimX, maxDimY - minDimY + 1, -minDimY
}

func parseInput(lines []string, lineParser func(string) (*Edge, error)) ([]Edge, error) {
	var edges []Edge

	for _, line := range lines {
		edge, error := lineParser(line)
		if error != nil {
			return nil, error
		}

		edges = append(edges, *edge)
	}

	return edges, nil
}

func parseLinePart1(line string) (*Edge, error) {
	var edge Edge

	tokens := strings.Split(line, " ")
	if len(tokens) != 3 {
		return nil, fmt.Errorf("Invalid line: %s", line)
	}

	switch tokens[0] {
	case "U":
		edge.Direction = UP
	case "R":
		edge.Direction = RIGHT
	case "D":
		edge.Direction = DOWN
	case "L":
		edge.Direction = LEFT
	default:
		return nil, fmt.Errorf("Invalid direction: %s", tokens[0])
	}

	length, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, err
	}
	edge.Length = length

	return &edge, nil
}

func parseLinePart2(line string) (*Edge, error) {
	var edge Edge

	tokens := strings.Split(line, " ")
	if len(tokens) != 3 {
		return nil, fmt.Errorf("Invalid line: %s", line)
	}

	analyzedString := strings.TrimPrefix(tokens[2], "(#")
	analyzedString = strings.TrimSuffix(analyzedString, ")")

	number, err := strconv.ParseUint(analyzedString[:5], 16, 32)
	if err != nil {
		return nil, err
	}
	edge.Length = int(number)

	switch analyzedString[5] {
	case '0':
		edge.Direction = RIGHT
	case '1':
		edge.Direction = DOWN
	case '2':
		edge.Direction = LEFT
	case '3':
		edge.Direction = UP
	default:
		return nil, fmt.Errorf("Invalid direction: %s", tokens[0])
	}

	return &edge, nil
}

func readLines(path string) []string {
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
