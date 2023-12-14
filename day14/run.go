package day14

import (
	"bufio"
	"fmt"
	"os"
)

const INPUT_FILE_PATH = "day14/input.txt"
const NUM_OF_CYCLES = 1000000000

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	platform := parseInput(lines)

	northPlatform := platform.Tilt(NORTH)
	fmt.Println("Load [PART 1]: ", northPlatform.LoadOnNorthBeams())

	platformAfterCycles := RunNCycles(platform, NUM_OF_CYCLES)
	println("Load [PART 2]: ", platformAfterCycles.LoadOnNorthBeams())
}

func RunNCycles(platform *Platform, n int) *Platform {
	platformCache := []*Platform{platform}
	startPeriod, endPeriod := -1, -1
	periodFound := false

	for it := 0; it < n && !periodFound; it++ {
		newPlatform := platformCache[it].RunCycle()
		for i := 0; i < it+1 && !periodFound; i++ {
			if newPlatform.Compare(platformCache[i]) {
				startPeriod, endPeriod = i, it+1
				periodFound = true
			}
		}
		if !periodFound {
			platformCache = append(platformCache, newPlatform)
		}
	}

	if !periodFound {
		return platformCache[n]
	}

	periodicPlatforms := platformCache[startPeriod:endPeriod]
	platformIndex := (n - startPeriod) % (endPeriod - startPeriod)
	return periodicPlatforms[platformIndex]
}

type Platform struct {
	Bytes [][]byte
	DimX  int
	DimY  int
}

func (platform *Platform) Print() {
	for y := 0; y < platform.DimY; y++ {
		for x := 0; x < platform.DimX; x++ {
			print(string(platform.Bytes[y][x]))
		}
		println()
	}
	println()
}

func (p1 *Platform) Compare(p2 *Platform) bool {
	if p1.DimX != p2.DimX || p1.DimY != p2.DimY {
		return false
	}

	for y := 0; y < p1.DimY; y++ {
		for x := 0; x < p1.DimX; x++ {
			if p1.Bytes[y][x] != p2.Bytes[y][x] {
				return false
			}
		}
	}
	return true
}

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

func (platform *Platform) RunCycle() *Platform {
	directions := []Direction{NORTH, WEST, SOUTH, EAST}
	for _, direction := range directions {
		platform = platform.Tilt(direction)
	}
	return platform
}

func (platform *Platform) Tilt(direction Direction) *Platform {
	switch direction {
	case NORTH:
		return platform.doTilt(true, true)
	case EAST:
		return platform.doTilt(false, false)
	case SOUTH:
		return platform.doTilt(true, false)
	case WEST:
		return platform.doTilt(false, true)
	}

	return platform
}

func (platform *Platform) doTilt(vertical bool, reverse bool) *Platform {
	var firstDim, secondDim int
	var value func(int, int) byte
	var setValue func(int, int, byte)

	bytes := make([][]byte, platform.DimY)
	for i := 0; i < platform.DimY; i++ {
		bytes[i] = make([]byte, platform.DimX)
	}

	if vertical {
		firstDim, secondDim = platform.DimX, platform.DimY
		value = func(i, j int) byte { return platform.Bytes[j][i] }
		setValue = func(i, j int, value byte) { bytes[j][i] = value }

	} else {
		firstDim, secondDim = platform.DimY, platform.DimX
		value = func(i, j int) byte { return platform.Bytes[i][j] }
		setValue = func(i, j int, value byte) { bytes[i][j] = value }
	}

	for i := 0; i < firstDim; i++ {
		lane := make([]byte, platform.DimX)
		for j := 0; j < secondDim; j++ {
			lane[j] = value(i, j)
		}

		tiltedLane := TiltLane(lane, reverse)
		for j := 0; j < secondDim; j++ {
			setValue(i, j, tiltedLane[j])
		}
	}

	return &Platform{bytes, platform.DimX, platform.DimY}

}

func TiltLane(lane []byte, reverse bool) []byte {
	tiltedLane := make([]byte, len(lane))
	roundRocks := 0

	for it := 0; it < len(lane); it++ {
		index := it
		if reverse {
			index = len(lane) - 1 - it
		}

		switch lane[index] {
		case '.':
			tiltedLane[index] = '.'
		case 'O':
			tiltedLane[index] = '.'
			roundRocks++
		case '#':
			tiltedLane[index] = '#'
			startIndex := index - 1
			if reverse {
				startIndex = index + 1
			}

			FillRoundRocks(tiltedLane, startIndex, roundRocks, reverse)
			roundRocks = 0
		}
	}

	if roundRocks > 0 {
		startIndex := len(lane) - 1
		if reverse {
			startIndex = 0
		}

		FillRoundRocks(tiltedLane, startIndex, roundRocks, reverse)
	}

	return tiltedLane
}

func FillRoundRocks(lane []byte, index, rocks int, reverse bool) {
	for offset := 0; offset < rocks; offset++ {
		if reverse {
			lane[index+offset] = 'O'
		} else {
			lane[index-offset] = 'O'
		}
	}
}

func (platform *Platform) LoadOnNorthBeams() int {
	load := 0

	for y := 0; y < platform.DimY; y++ {
		for x := 0; x < platform.DimX; x++ {
			if platform.Bytes[y][x] == 'O' {
				load += platform.DimY - y
			}
		}
	}

	return load
}

func parseInput(lines []string) *Platform {
	bytes := make([][]byte, len(lines))
	dimY, dimX := len(lines), len(lines[0])

	for i, line := range lines {
		bytes[i] = []byte(line)
	}

	return &Platform{bytes, dimX, dimY}
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
