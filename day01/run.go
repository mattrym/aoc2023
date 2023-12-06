package day01

import (
	"bufio"
	"os"
	"strings"
)

type ma map[string]int

const INPUT_FILE_PATH = "day01/input.txt"

func Run() {
	// Read input file (from input.txt)
	inputLines := readLines(INPUT_FILE_PATH)
	stringsToDigits := mapStringsToDigits()

	totalCalibrationValuePart1 := 0
	totalCalibrationValuePart2 := 0

	for _, inputLine := range inputLines {
		calibrationValuePart1 := FindCalibrationValuePart1(inputLine)
		calibrationValuePart2 := FindCalibrationValuePart2(inputLine, stringsToDigits)

		totalCalibrationValuePart1 += calibrationValuePart1
		totalCalibrationValuePart2 += calibrationValuePart2
	}

	// Print result
	println("Total calibration value [PART 1]: ", totalCalibrationValuePart1)
	println("Total calibration value [PART 2]: ", totalCalibrationValuePart2)
}

func FindCalibrationValuePart1(inputLine string) int {
	bytes := []byte(inputLine)
	var result int

	for i := 0; ; i++ {
		if bytes[i] >= 48 && bytes[i] <= 57 {
			result += (int(bytes[i]) - 48) * 10
			break
		}
	}

	for i := len(bytes) - 1; ; i-- {
		if bytes[i] >= 48 && bytes[i] <= 57 {
			result += int(bytes[i]) - 48
			break
		}
	}

	return result
}

func FindCalibrationValuePart2(inputLine string, stringsToDigits ma) int {
	digits := FindDigits(inputLine, stringsToDigits)
	return digits[0]*10 + digits[len(digits)-1]

}

func FindDigits(inputLine string, stringsToDigits ma) []int {
	var digits []int
	var substrIndex int
	var minSubstrIndex int
	var minSubstr string

	for {
		minSubstrIndex = -1

		for key, _ := range stringsToDigits {
			substrIndex = strings.Index(inputLine, key)
			if substrIndex != -1 && (minSubstrIndex == -1 || substrIndex < minSubstrIndex) {
				minSubstrIndex = substrIndex
				minSubstr = key
			}
		}

		if minSubstrIndex != -1 {
			digits = append(digits, stringsToDigits[minSubstr])
			inputLine = inputLine[(minSubstrIndex + 1):]
		} else {
			break
		}
	}

	return digits
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

func mapStringsToDigits() ma {
	return ma{
		"zero":  0,
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
		"0":     0,
		"1":     1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
	}
}
