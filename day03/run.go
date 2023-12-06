package day03

import (
	"bufio"
	"fmt"
	"os"
)

const INPUT_FILE_PATH = "day03/input.txt"

type Asterisk struct {
	X               int
	Y               int
	AdjacentNumbers []int
}

type Coords struct {
	X int
	Y int
}

type Number struct {
	X   int
	Y   int
	Len int
	Val int
}

type Bytemap struct {
	Bytes   [][]byte
	RowSize int
	ColSize int
}

type Asterisks map[Coords][]int

func Run() {
	bytemap := readBytemap(INPUT_FILE_PATH)
	numbers := FindAllNumbers(bytemap)
	asterisks := FindAllAsterisks(bytemap)
	var sum1, sum2 int

	// Part 1
	for _, coords := range *numbers {
		if IsPartNumber(bytemap, coords) {
			sum1 += coords.Val
		}
	}
	fmt.Println("Sum of part numbers [PART 1]: ", sum1)

	// Part 2
	for _, coords := range *numbers {
		adjacent, asteriskCoords := IsAdjacentToAsterisk(bytemap, coords)
		if adjacent {
			asterisks[asteriskCoords] = append(asterisks[asteriskCoords], coords.Val)
		}
	}

	for _, asterisk := range asterisks {
		if len(asterisk) == 2 {
			sum2 += asterisk[0] * asterisk[1]
		}
	}

	fmt.Println("Sum of gear ratios [PART 2]: ", sum2)
}

func findAsterisk(asterisks []Asterisk, coords Coords) *Asterisk {
	for _, asterisk := range asterisks {
		if asterisk.X == coords.X && asterisk.Y == coords.Y {
			return &asterisk
		}
	}
	return nil
}

func IsAdjacentToAsterisk(bytemap *Bytemap, coords Number) (bool, Coords) {
	// previous row
	if coords.X != 0 {
		startOffset, endOffset := findStartOffset(bytemap, coords), findEndOffset(bytemap, coords)

		for columnIndex := startOffset; columnIndex < endOffset; columnIndex++ {
			if isAsterisk(bytemap.Bytes[coords.X-1][columnIndex]) {
				return true, Coords{X: coords.X - 1, Y: columnIndex}
			}
		}
	}

	// current row
	if coords.Y != 0 {
		if isAsterisk(bytemap.Bytes[coords.X][coords.Y-1]) {
			return true, Coords{X: coords.X, Y: coords.Y - 1}
		}
	}

	if coords.Y+coords.Len < bytemap.ColSize {
		if isAsterisk(bytemap.Bytes[coords.X][coords.Y+coords.Len]) {
			return true, Coords{X: coords.X, Y: coords.Y + coords.Len}
		}
	}

	// next row
	if coords.X != bytemap.ColSize-1 {
		startOffset, endOffset := findStartOffset(bytemap, coords), findEndOffset(bytemap, coords)

		for columnIndex := startOffset; columnIndex < endOffset; columnIndex++ {
			if isAsterisk(bytemap.Bytes[coords.X+1][columnIndex]) {
				return true, Coords{X: coords.X + 1, Y: columnIndex}
			}
		}
	}

	return false, Coords{}

}

func IsPartNumber(bytemap *Bytemap, coords Number) bool {
	// previous row
	if coords.X != 0 {
		startOffset, endOffset := findStartOffset(bytemap, coords), findEndOffset(bytemap, coords)

		for columnIndex := startOffset; columnIndex < endOffset; columnIndex++ {
			if isSymbolChar(bytemap.Bytes[coords.X-1][columnIndex]) {
				return true
			}
		}
	}

	// current row
	if coords.Y != 0 {
		if isSymbolChar(bytemap.Bytes[coords.X][coords.Y-1]) {
			return true
		}
	}

	if coords.Y+coords.Len < bytemap.ColSize {
		if isSymbolChar(bytemap.Bytes[coords.X][coords.Y+coords.Len]) {
			return true
		}
	}

	// next row
	if coords.X != bytemap.ColSize-1 {
		startOffset, endOffset := findStartOffset(bytemap, coords), findEndOffset(bytemap, coords)

		for columnIndex := startOffset; columnIndex < endOffset; columnIndex++ {
			if isSymbolChar(bytemap.Bytes[coords.X+1][columnIndex]) {
				return true
			}
		}
	}

	return false

}

func findStartOffset(bytemap *Bytemap, coords Number) int {
	if coords.Y == 0 {
		return 0
	}
	return coords.Y - 1
}

func findEndOffset(bytemap *Bytemap, coords Number) int {
	if coords.Y+coords.Len == bytemap.RowSize {
		return coords.Y + coords.Len
	}
	return coords.Y + coords.Len + 1
}

func FindAllAsterisks(bytemap *Bytemap) Asterisks {
	result := make(Asterisks)

	for x := 0; x < bytemap.ColSize; x++ {
		for y := 0; y < bytemap.RowSize; y++ {
			if bytemap.Bytes[x][y] == 42 {
				result[Coords{X: x, Y: y}] = []int{}
			}
		}
	}

	return result
}

func FindAllNumbers(bytemap *Bytemap) *[]Number {
	var result []Number

	for x := 0; x < bytemap.ColSize; x++ {
		for y := 0; y < bytemap.RowSize; {
			z, len := 0, 0

			if isDigitChar(bytemap.Bytes[x][y]) {
				len++
				for z = y + 1; z < bytemap.RowSize; z++ {
					if isDigitChar(bytemap.Bytes[x][z]) {
						len++
					} else {
						break
					}
				}
			}

			if len > 0 {
				coords := Number{X: x, Y: y, Len: len}
				coords.Val = findNumberValue(bytemap, coords)
				result = append(result, coords)

				y = z
			} else {
				y++
			}
		}
	}

	return &result
}

func findNumberValue(bytemap *Bytemap, coords Number) int {
	value := 0
	for i := 0; i < coords.Len; i++ {
		value = value*10 + charToDigit(bytemap.Bytes[coords.X][coords.Y+i])
	}

	return value
}

func isSymbolChar(char byte) bool {
	return !isDigitChar(char) && char != 46
}

func isDigitChar(char byte) bool {
	return char >= 48 && char <= 57
}

func charToDigit(char byte) int {
	return int(char) - 48
}

func isAsterisk(char byte) bool {
	return int(char) == 42
}

func readBytemap(path string) *Bytemap {
	// Read input file (line by line) from input.txt
	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	var result [][]byte
	var rowSize int

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		row := []byte(scanner.Text())
		if len(row) > rowSize {
			rowSize = len(row)
		}
		result = append(result, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return &Bytemap{
		Bytes:   result,
		RowSize: rowSize,
		ColSize: len(result),
	}
}
