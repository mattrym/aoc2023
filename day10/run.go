package day10

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
)

const INPUT_FILE_PATH = "day10/input.txt"

const INVALID = 0
const UP = 1
const DOWN = -1
const LEFT = -2
const RIGHT = 2

type Direction int
type Tile byte

type Point struct {
	X int
	Y int
}

type Loop []Point

type ByteMap struct {
	Bytes [][]byte
	DimX  int
	DimY  int
}

func (byteMap *ByteMap) Print() {
	for i := 0; i < byteMap.DimY; i++ {
		fmt.Println(string(byteMap.Bytes[i]))
	}
}

func (byteMap *ByteMap) PrintLoopAndInsidePoints(loop *Loop, sweepedPoints []Point) {
	red := color.New(color.FgRed)
	blue := color.New(color.FgBlue)

	for i := 0; i < byteMap.DimY; i++ {
		for j := 0; j < byteMap.DimX; j++ {
			if slices.Contains(*loop, Point{j, i}) {
				red.Print(string(byteMap.Bytes[i][j]))
			} else if slices.Contains(sweepedPoints, Point{j, i}) {
				blue.Print(string(byteMap.Bytes[i][j]))
			} else {
				fmt.Print(string(byteMap.Bytes[i][j]))
			}
		}
		fmt.Println()
	}
}

func (byteMap *ByteMap) Get(X int, Y int) Tile {
	return Tile(byteMap.Bytes[Y][X])
}

func (byteMap *ByteMap) GetPoint(point Point) Tile {
	return Tile(byteMap.Bytes[point.Y][point.X])
}

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	byteMap := parseInput(lines)

	loop, err := DiscoverLoop(byteMap)
	if err != nil {
		panic(err)
	}

	loopDistance := (len(*loop) + 1) / 2
	fmt.Println("Loop [PART 1]: ", loopDistance)

	sweepedPoints := byteMap.SweepAndFindInnerPoints(loop)
	fmt.Println("Sweeped points [PART 2]: ", len(sweepedPoints))

	byteMap.PrintLoopAndInsidePoints(loop, sweepedPoints)
}

func DiscoverLoop(byteMap *ByteMap) (*Loop, error) {
	startPoint, err := byteMap.StartPoint()
	var loop Loop
	if err != nil {
		return nil, err
	}

	loop = append(loop, *startPoint)
	currentPoint := *startPoint
	currentDirection := FindNextDirectionAfterStart(byteMap, startPoint)

	for {
		nextPoint := byteMap.MakeMove(currentPoint, currentDirection)
		nextDirection := byteMap.NextDirection(nextPoint, currentDirection)

		if byteMap.GetPoint(nextPoint) == 'S' {
			return &loop, nil
		}

		loop = append(loop, nextPoint)
		currentPoint = nextPoint
		currentDirection = nextDirection
	}
}

func (byteMap *ByteMap) SweepAndFindInnerPoints(loop *Loop) []Point {
	downConnectors := []Tile{'|', '7', 'F'}
	sweptPoints := []Point{}

	// if starting point is a down connector, we need to include it
	// simply check if it's possible to go down from the starting point
	startPoint := (*loop)[0]
	if (startPoint.Y + 1) < byteMap.DimY {
		pointBelow := Point{startPoint.X, startPoint.Y + 1}
		acceptedConnectors := []Tile{'|', 'L', '7'}
		byteBelow := byteMap.GetPoint(pointBelow)

		if slices.Contains(acceptedConnectors, byteBelow) {
			downConnectors = append(downConnectors, 'S')
		}
	}

	// iterate over all lines, at the beginning we're outside the loop
	for i := 0; i < byteMap.DimY; i++ {
		insideLoop := false
		for j := 0; j < byteMap.DimX; j++ {
			currPoint := Point{j, i}
			currTile := byteMap.GetPoint(currPoint)

			// if we're encountering a pipe going down, we're entering or exiting the loop
			if slices.Contains(*loop, currPoint) && slices.Contains(downConnectors, currTile) {
				insideLoop = !insideLoop
			}

			// if we're in the loop and current point is not a part of the loop
			// we need to add it to the swept points, that constitute the inner points
			if insideLoop && !slices.Contains(*loop, currPoint) {
				sweptPoints = append(sweptPoints, currPoint)
			}
		}
	}

	return sweptPoints
}

func (byteMap *ByteMap) StartPoint() (*Point, error) {
	for i := 0; i < byteMap.DimY; i++ {
		for j := 0; j < byteMap.DimX; j++ {
			if byteMap.Bytes[i][j] == 'S' {
				return &Point{j, i}, nil
			}
		}
	}
	return nil, errors.New("No start point found")
}

func FindNextDirectionAfterStart(byteMap *ByteMap, startPoint *Point) Direction {
	// point to left
	if (startPoint.X - 1) >= 0 {
		pointToLeft := Point{startPoint.X - 1, startPoint.Y}
		acceptedConnectors := []Tile{'-', 'L', 'F'}
		byteAtLeft := byteMap.GetPoint(pointToLeft)

		if slices.Contains(acceptedConnectors, byteAtLeft) {
			return LEFT
		}
	}

	// point to right
	if (startPoint.X + 1) < byteMap.DimX {
		pointToRight := Point{startPoint.X + 1, startPoint.Y}
		acceptedConnectors := []Tile{'-', 'J', '7'}
		byteAtRight := byteMap.GetPoint(pointToRight)

		if slices.Contains(acceptedConnectors, byteAtRight) {
			return RIGHT
		}
	}

	// point above
	if (startPoint.Y - 1) >= 0 {
		pointAbove := Point{startPoint.X, startPoint.Y - 1}
		acceptedConnectors := []Tile{'|', 'J', 'F'}
		byteAbove := byteMap.GetPoint(pointAbove)

		if slices.Contains(acceptedConnectors, byteAbove) {
			return UP
		}
	}

	// point below
	if (startPoint.Y + 1) < byteMap.DimY {
		pointBelow := Point{startPoint.X, startPoint.Y + 1}
		acceptedConnectors := []Tile{'|', 'L', '7'}
		byteBelow := byteMap.GetPoint(pointBelow)

		if slices.Contains(acceptedConnectors, byteBelow) {
			return DOWN
		}
	}

	return INVALID
}

func (byteMap *ByteMap) MakeMove(from Point, direction Direction) Point {
	switch direction {
	case UP:
		return Point{from.X, from.Y - 1}
	case DOWN:
		return Point{from.X, from.Y + 1}
	case LEFT:
		return Point{from.X - 1, from.Y}
	case RIGHT:
		return Point{from.X + 1, from.Y}
	default:
		return Point{-1, -1}
	}
}

func (byteMap *ByteMap) NextDirection(currPoint Point, prevDirection Direction) Direction {
	tile := byteMap.GetPoint(currPoint)
	oppositePrevDirection := OppositeDirection(prevDirection)
	possibleDirections := PossibleDirections(tile)

	for direction := range possibleDirections {
		if direction != oppositePrevDirection {
			return direction
		}
	}
	return INVALID
}

func PossibleDirections(tile Tile) map[Direction]bool {
	switch tile {
	case '|':
		return map[Direction]bool{UP: true, DOWN: true}
	case '-':
		return map[Direction]bool{LEFT: true, RIGHT: true}
	case 'J':
		return map[Direction]bool{UP: true, LEFT: true}
	case 'L':
		return map[Direction]bool{UP: true, RIGHT: true}
	case '7':
		return map[Direction]bool{DOWN: true, LEFT: true}
	case 'F':
		return map[Direction]bool{DOWN: true, RIGHT: true}
	default:
		return map[Direction]bool{}
	}
}

func OppositeDirection(direction Direction) Direction {
	switch direction {
	case UP:
		return DOWN
	case DOWN:
		return UP
	case LEFT:
		return RIGHT
	case RIGHT:
		return LEFT
	default:
		return INVALID
	}
}

func parseInput(lines []string) *ByteMap {
	bytes := make([][]byte, len(lines))
	dimY, dimX := len(lines), len(lines[0])

	for i, line := range lines {
		bytes[i] = []byte(line)
	}

	return &ByteMap{bytes, dimX, dimY}
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
