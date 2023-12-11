package day11

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const INPUT_FILE_PATH = "day11/input.txt"

type Point struct {
	X int
	Y int
}

type Image struct {
	Bytes     [][]byte
	EmptyRows []int
	EmptyCols []int
	DimX      int
	DimY      int
}

func (image *Image) Print() {
	for i := 0; i < image.DimY; i++ {
		fmt.Println(string(image.Bytes[i]))
	}
}

func (image *Image) FindGalaxies() []Point {
	galaxies := []Point{}

	for i := 0; i < image.DimY; i++ {
		for j := 0; j < image.DimX; j++ {
			if image.Bytes[i][j] == '#' {
				galaxies = append(galaxies, Point{j, i})
			}
		}
	}

	return galaxies
}

func (image *Image) Distance(p1 Point, p2 Point, emptyCoef int) int {
	x1, x2 := p1.X, p2.X
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	y1, y2 := p1.Y, p2.Y
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	xDiff := 0
	for i := x1 + 1; i <= x2; i++ {
		if slices.Contains(image.EmptyCols, i) {
			xDiff += emptyCoef
		} else {
			xDiff += 1
		}
	}

	yDiff := 0
	for i := y1 + 1; i <= y2; i++ {
		if slices.Contains(image.EmptyRows, i) {
			yDiff += emptyCoef
		} else {
			yDiff += 1
		}
	}

	return xDiff + yDiff
}

func (image *Image) Expand(*Image) *Image {
	emptyRows, emptyCols := []int{}, []int{}

	for i := 0; i < image.DimY; i++ {
		isEmpty := true
		for j := 0; j < image.DimX; j++ {
			if image.Bytes[i][j] == '#' {
				isEmpty = false
				continue
			}
		}

		if isEmpty {
			emptyRows = append(emptyRows, i)
		}
	}

	for j := 0; j < image.DimX; j++ {
		isEmpty := true
		for i := 0; i < image.DimY; i++ {
			if image.Bytes[i][j] == '#' {
				isEmpty = false
				continue
			}
		}

		if isEmpty {
			emptyCols = append(emptyCols, j)
		}
	}

	image.EmptyRows = emptyRows
	image.EmptyCols = emptyCols
	return image
}

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	image := parseInput(lines)

	expandedImage := image.Expand(image)
	galaxies := expandedImage.FindGalaxies()

	sumDistancesPart1 := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			sumDistancesPart1 += image.Distance(galaxies[i], galaxies[j], 2)
		}
	}
	fmt.Println("Sum of distances [PART 1]: ", sumDistancesPart1)

	sumDistancesPart2 := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			sumDistancesPart2 += image.Distance(galaxies[i], galaxies[j], 1000000)
		}
	}
	fmt.Println("Sum of distances [PART 2]: ", sumDistancesPart2)
}

func parseInput(lines []string) *Image {
	var image Image

	image.DimX = len(lines[0])
	image.DimY = len(lines)

	for _, line := range lines {
		image.Bytes = append(image.Bytes, []byte(line))
	}

	return &image
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
