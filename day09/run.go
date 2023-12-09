package day09

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const INPUT_FILE_PATH = "day09/input.txt"

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	sequences := parseInput(lines)

	sumOfExtrapolatedStarts := 0
	sumOfExtrapolatedEnds := 0

	for _, sequence := range sequences {
		extrapolatedSequence := extrapolateSequence(sequence)
		sumOfExtrapolatedStarts += extrapolatedSequence[0]
		sumOfExtrapolatedEnds += extrapolatedSequence[len(extrapolatedSequence)-1]
	}

	fmt.Println("Sum of extrapolated values [PART 1]: ", sumOfExtrapolatedEnds)
	fmt.Println("Sum of extrapolated values [PART 2]: ", sumOfExtrapolatedStarts)
}

func extrapolateSequence(sequence []int) []int {
	sequences := [][]int{}
	sequences = append(sequences, sequence)
	var newSequenceHasOnlyZeros bool

	for !newSequenceHasOnlyZeros {
		lastSequence := sequences[len(sequences)-1]
		newSequence := make([]int, len(lastSequence)-1)
		newSequenceHasOnlyZeros = true

		for i := 0; i < len(lastSequence)-1; i++ {
			newSequence[i] = lastSequence[i+1] - lastSequence[i]
			if newSequence[i] != 0 {
				newSequenceHasOnlyZeros = false
			}
		}

		if newSequenceHasOnlyZeros {
			newSequence = append(newSequence, 0)
		}
		sequences = append(sequences, newSequence)
	}

	for i := len(sequences) - 2; i >= 0; i-- {
		newSequence := make([]int, len(sequences[i])+2)
		copy(newSequence[1:], sequences[i])

		startValue := sequences[i][0] - sequences[i+1][0]
		endValue := sequences[i][len(sequences[i])-1] + sequences[i+1][len(sequences[i+1])-1]

		newSequence[0] = startValue
		newSequence[len(newSequence)-1] = endValue
		sequences[i] = newSequence
	}

	return sequences[0]
}

func parseInput(lines []string) [][]int {
	sequences := [][]int{}

	for _, line := range lines {
		tokens := strings.Split(line, " ")
		sequence := make([]int, len(tokens))

		for j, token := range tokens {
			number, err := strconv.Atoi(token)
			if err != nil {
				panic(err)
			}
			sequence[j] = number
		}

		sequences = append(sequences, sequence)
	}

	return sequences
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
