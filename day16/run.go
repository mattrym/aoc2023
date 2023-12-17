package day16

import (
	"bufio"
	"fmt"
	"os"
)

const INPUT_FILE_PATH = "day16/input.txt"

type Direction int

const (
	UP    Direction = 1
	RIGHT Direction = 2
	DOWN  Direction = 4
	LEFT  Direction = 8
)

type Contraption struct {
	Grid  [][]byte
	Beams [][]Direction
	DimX  int
	DimY  int
}

type Position struct {
	X        int
	Y        int
	EntryDir Direction
}

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	contraption := parseInput(lines)

	energizedTilesFromTopLeft := contraption.EnergizedTiles(Position{0, 0, RIGHT})
	fmt.Println("Number of energized tiles from top left [PART 1]: ", energizedTilesFromTopLeft)

	maxEnergizedTiles := contraption.FindMaxEnergizedTiles()
	fmt.Println("Max number of energized tiles [PART 2]: ", maxEnergizedTiles)

}

func (contraption *Contraption) FindMaxEnergizedTiles() int {
	maxEnergizedTiles := 0

	for y := 0; y < contraption.DimY; y++ {
		energizedTiles := contraption.EnergizedTiles(Position{0, y, RIGHT})
		if energizedTiles > maxEnergizedTiles {
			maxEnergizedTiles = energizedTiles
		}

		energizedTiles = contraption.EnergizedTiles(Position{contraption.DimX - 1, y, LEFT})
		if energizedTiles > maxEnergizedTiles {
			maxEnergizedTiles = energizedTiles
		}
	}

	for x := 0; x < contraption.DimX; x++ {
		energizedTiles := contraption.EnergizedTiles(Position{x, 0, DOWN})
		if energizedTiles > maxEnergizedTiles {
			maxEnergizedTiles = energizedTiles
		}

		energizedTiles = contraption.EnergizedTiles(Position{x, contraption.DimY - 1, UP})
		if energizedTiles > maxEnergizedTiles {
			maxEnergizedTiles = energizedTiles
		}
	}

	return maxEnergizedTiles
}

func (contraption *Contraption) EnergizedTiles(initialPosition Position) int {
	contraption.Beams = make([][]Direction, contraption.DimY)
	for y := 0; y < contraption.DimY; y++ {
		contraption.Beams[y] = make([]Direction, contraption.DimX)
	}

	oldPositions := make(map[Position]bool)
	analyzedPositions := []Position{initialPosition}

	for len(analyzedPositions) > 0 {
		analyzedPosition := analyzedPositions[0]
		analyzedPositions = analyzedPositions[1:]

		nextPositions := contraption.EnergizeTile(analyzedPosition)
		oldPositions[analyzedPosition] = true

		for _, nextPosition := range nextPositions {
			if !oldPositions[nextPosition] {
				analyzedPositions = append(analyzedPositions, nextPosition)
			}
		}
	}

	numEnergizedTiles := 0
	for y := 0; y < contraption.DimY; y++ {
		for x := 0; x < contraption.DimX; x++ {
			if contraption.Beams[y][x] != 0 {
				numEnergizedTiles++
			}
		}
	}

	return numEnergizedTiles
}

func (contraption *Contraption) EnergizeTile(position Position) []Position {
	tile := contraption.Grid[position.Y][position.X]
	beamDirections := DetermineBeamDirections(tile, position.EntryDir)
	nextAnalyzedPositions := []Position{}

	oldBeams := make([][]Direction, contraption.DimY)
	for y := 0; y < contraption.DimY; y++ {
		oldBeams[y] = make([]Direction, contraption.DimX)
		copy(oldBeams[y], contraption.Beams[y])
	}

	for _, beamDirection := range beamDirections {
		contraption.Beams[position.Y][position.X] |= beamDirection
		neighborPosition, ok := contraption.FindNeighbor(position, beamDirection)

		if ok {
			nextAnalyzedPositions = append(nextAnalyzedPositions, neighborPosition)
		}
	}

	return nextAnalyzedPositions
}

func (contraption *Contraption) FindNeighbor(position Position, beamDirection Direction) (Position, bool) {
	switch beamDirection {
	case UP:
		if position.Y == 0 {
			return Position{}, false
		}
		return Position{position.X, position.Y - 1, UP}, true
	case RIGHT:
		if position.X == contraption.DimX-1 {
			return Position{}, false
		}
		return Position{position.X + 1, position.Y, RIGHT}, true
	case DOWN:
		if position.Y == contraption.DimY-1 {
			return Position{}, false
		}
		return Position{position.X, position.Y + 1, DOWN}, true
	case LEFT:
		if position.X == 0 {
			return Position{}, false
		}
		return Position{position.X - 1, position.Y, LEFT}, true
	}
	return Position{}, false
}

func (contraption *Contraption) PrintTiles() {
	for y := 0; y < contraption.DimY; y++ {
		fmt.Println(string(contraption.Grid[y]))
	}
	fmt.Println()
}

func (contraption *Contraption) PrintEnergizedTiles() {
	for y := 0; y < contraption.DimY; y++ {
		for x := 0; x < contraption.DimX; x++ {
			if contraption.Beams[y][x] != 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func DetermineBeamDirections(tile byte, entryDir Direction) []Direction {
	switch tile {
	case '.':
		return []Direction{entryDir}
	case '/':
		return slashDirections(entryDir)
	case '\\':
		return bashslashDirections(entryDir)
	case '|':
		return verticalDirections(entryDir)
	case '-':
		return horizontalDirections(entryDir)
	}
	return []Direction{}
}

func slashDirections(entryDir Direction) []Direction {
	switch entryDir {
	case UP:
		return []Direction{RIGHT}
	case RIGHT:
		return []Direction{UP}
	case DOWN:
		return []Direction{LEFT}
	case LEFT:
		return []Direction{DOWN}
	}
	return []Direction{}
}

func bashslashDirections(entryDir Direction) []Direction {
	switch entryDir {
	case UP:
		return []Direction{LEFT}
	case RIGHT:
		return []Direction{DOWN}
	case DOWN:
		return []Direction{RIGHT}
	case LEFT:
		return []Direction{UP}
	}
	return []Direction{}
}

func verticalDirections(entryDir Direction) []Direction {
	switch entryDir {
	case UP:
		return []Direction{UP}
	case RIGHT:
		return []Direction{UP, DOWN}
	case DOWN:
		return []Direction{DOWN}
	case LEFT:
		return []Direction{UP, DOWN}
	}
	return []Direction{}
}

func horizontalDirections(entryDir Direction) []Direction {
	switch entryDir {
	case UP:
		return []Direction{LEFT, RIGHT}
	case RIGHT:
		return []Direction{RIGHT}
	case DOWN:
		return []Direction{LEFT, RIGHT}
	case LEFT:
		return []Direction{LEFT}
	}
	return []Direction{}
}

func parseInput(lines []string) *Contraption {
	contraption := Contraption{}

	for _, line := range lines {
		if line == "" {
			break
		}

		contraption.Grid = append(contraption.Grid, []byte(line))
		contraption.DimX = len(line)
		contraption.DimY++
	}

	return &contraption
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
