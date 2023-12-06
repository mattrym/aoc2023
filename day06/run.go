package day06

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	Time     int
	Distance int
}

const INPUT_FILE_PATH = "day06/input.txt"

func Run() {
	races, err := readInputPart1(readLines(INPUT_FILE_PATH))
	if err != nil {
		panic(err)
	}

	product := 1
	for _, race := range races {
		numBetterTimes := CountBetterTimes(race)
		product *= numBetterTimes
	}

	fmt.Println("Product of better times [PART 1]: ", product)

	race, err := readInputPart2(readLines(INPUT_FILE_PATH))
	if err != nil {
		panic(err)
	}

	fmt.Println("Better times for race [PART 2]: ", CountBetterTimes(*race))
}

func CountBetterTimes(race Race) int {
	count := 0

	for time := 0; time <= race.Time; time++ {
		if IsBetterTime(race, time) {
			count++
		}
	}

	return count
}

func IsBetterTime(race Race, timeToHold int) bool {
	return Distance(timeToHold, race.Time)-race.Distance > 0
}

func Distance(timeToHold int, timeTotal int) int {
	return (timeTotal - timeToHold) * timeToHold
}

func readInputPart1(lines []string) ([]Race, error) {
	if len(lines) < 2 {
		return nil, errors.New("Invalid input: too few lines")
	}

	timesString, _ := strings.CutPrefix(lines[0], "Time:")
	timeStrings := strings.Fields(timesString)
	fmt.Println(timeStrings)

	distancesString, _ := strings.CutPrefix(lines[1], "Distance:")
	distanceStrings := strings.Fields(distancesString)
	fmt.Println(distanceStrings)

	if len(timeStrings) != len(distanceStrings) {
		return nil, errors.New("Invalid input: unequal number of times and distances")
	}

	var result []Race
	for i := 0; i < len(timeStrings); i++ {
		time, err := strconv.Atoi(timeStrings[i])
		if err != nil {
			return nil, fmt.Errorf("Invalid input: invalid time#%d: %s", i, timeStrings[i])
		}

		distance, err := strconv.Atoi(distanceStrings[i])
		if err != nil {
			return nil, fmt.Errorf("Invalid input: invalid distance #%d: %s", i, distanceStrings[i])
		}

		race := Race{Time: time, Distance: distance}
		result = append(result, race)
	}

	return result, nil
}

func readInputPart2(lines []string) (*Race, error) {
	if len(lines) < 2 {
		return nil, errors.New("Invalid input: too few lines")
	}

	timesString, _ := strings.CutPrefix(lines[0], "Time:")
	timeStrings := strings.Fields(timesString)
	fmt.Println(timeStrings)

	distancesString, _ := strings.CutPrefix(lines[1], "Distance:")
	distanceStrings := strings.Fields(distancesString)
	fmt.Println(distanceStrings)

	if len(timeStrings) != len(distanceStrings) {
		return nil, errors.New("Invalid input: unequal number of times and distances")
	}

	timeString := strings.Join(timeStrings, "")
	distanceString := strings.Join(distanceStrings, "")

	time, err := strconv.Atoi(timeString)
	if err != nil {
		return nil, fmt.Errorf("Invalid input: invalid time: %s", timeString)
	}

	distance, err := strconv.Atoi(distanceString)
	if err != nil {
		return nil, fmt.Errorf("Invalid input: invalid distance: %s", distanceString)
	}

	return &Race{Time: time, Distance: distance}, nil
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
