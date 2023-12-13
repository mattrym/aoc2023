package day12

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const INPUT_FILE_PATH = "day12/input.txt"
const UNFOLD_MULTIPLIER = 5

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	springRows := parseInput(lines)

	sumOfCountOfValidVariants := 0
	for _, springRow := range springRows {
		sumOfCountOfValidVariants += SolveBruteForcefullyWithHeuristics(springRow)
	}

	fmt.Println("Sum of count of valid variants [PART 1]: ", sumOfCountOfValidVariants)

	sumOfCountOfUnfoldedValidVariants := 0
	for _, springRow := range springRows {
		dynamicProgramming := DynamicProgramming{
			Pattern: unfoldPattern(springRow.Pattern, UNFOLD_MULTIPLIER),
			Groups:  unfoldGroups(springRow.Groups, UNFOLD_MULTIPLIER),
		}
		sumOfCountOfUnfoldedValidVariants += dynamicProgramming.Solve()
	}

	fmt.Println("Sum of count of valid variants [PART 2]: ", sumOfCountOfUnfoldedValidVariants)

}

// PART 2

func unfoldPattern(pattern string, coef int) string {
	patternDups := []string{}
	for i := 0; i < coef; i++ {
		patternDups = append(patternDups, pattern)
	}
	return strings.Join(patternDups, "?") + "."
}

func unfoldGroups(groups []int, coef int) []int {
	var groupsDups []int
	for i := 0; i < coef; i++ {
		groupsDups = append(groupsDups, groups...)
	}
	return groupsDups
}

type DynamicProgramming struct {
	Memory  [][]int
	Pattern string
	Groups  []int
}

func (dp *DynamicProgramming) Solve() int {
	dp.Memory = [][]int{}
	for i := 0; i < len(dp.Pattern); i++ {
		dp.Memory = append(dp.Memory, make([]int, len(dp.Groups)))
		for j := 0; j < len(dp.Groups); j++ {
			dp.Memory[i][j] = -1
		}
	}

	return dp.Dispatch(0, 0)
}

func (dp *DynamicProgramming) GetFromMemory(patternPos, groupPos int) (int, bool) {
	if dp.Memory[patternPos][groupPos] == -1 {
		return 0, false
	}

	return dp.Memory[patternPos][groupPos], true
}

func (dp *DynamicProgramming) Dispatch(patternPos, groupPos int) int {
	if patternPos == len(dp.Pattern) {
		// we have no more pattern to match against the groups
		// check if we have any groups left
		if groupPos == len(dp.Groups) {
			return 1
		}
		return 0
	}

	if groupPos == len(dp.Groups) {
		// we have no more groups to match against the pattern
		// check if we have only dots left in the pattern
		for i := patternPos; i < len(dp.Pattern); i++ {
			if dp.Pattern[i] == '#' {
				return 0
			}
		}
		return 1
	}

	if result, ok := dp.GetFromMemory(patternPos, groupPos); ok {
		return result
	}

	switch dp.Pattern[patternPos] {
	case '.':
		return dp.DispatchDot(patternPos, groupPos)
	case '#':
		return dp.DispatchHash(patternPos, groupPos)
	case '?':
		return dp.DispatchQuestionMark(patternPos, groupPos)
	default:
		panic("Unknown pattern sign")
	}
}

func (dp *DynamicProgramming) DispatchDot(patternPos, groupPos int) int {
	dp.Memory[patternPos][groupPos] = dp.Dispatch(patternPos+1, groupPos)
	return dp.Memory[patternPos][groupPos]
}

func (dp *DynamicProgramming) DispatchQuestionMark(patternPos, groupPos int) int {
	dp.Memory[patternPos][groupPos] = dp.DispatchDot(patternPos, groupPos) +
		dp.DispatchHash(patternPos, groupPos)
	return dp.Memory[patternPos][groupPos]
}

func (dp *DynamicProgramming) DispatchHash(patternPos, groupPos int) int {
	if !groupMatchesPattern(dp.Pattern, dp.Groups[groupPos], patternPos) {
		dp.Memory[patternPos][groupPos] = 0
		return dp.Memory[patternPos][groupPos]
	}

	nextPatternPos := patternPos + dp.Groups[groupPos]
	if nextPatternPos == len(dp.Pattern) {
		if groupPos == len(dp.Groups) {
			dp.Memory[patternPos][groupPos] = 1
		} else {
			dp.Memory[patternPos][groupPos] = 0
		}
		return dp.Memory[patternPos][groupPos]
	}

	dp.Memory[patternPos][groupPos] = dp.Dispatch(nextPatternPos+1, groupPos+1)
	return dp.Memory[patternPos][groupPos]
}

func groupMatchesPattern(pattern string, groupSize, patternPos int) bool {
	hashCount := 0

	for pos := patternPos; pos < len(pattern) && hashCount <= groupSize; pos++ {
		sign := pattern[pos]
		if sign == '?' {
			if hashCount < groupSize {
				hashCount++
			} else {
				break
			}
		}
		if sign == '#' {
			hashCount++
		}
		if sign == '.' {
			break
		}
	}

	return hashCount == groupSize
}

type Springs struct {
	Pattern string
	Groups  []int
}

func SolveBruteForcefullyWithHeuristics(springs Springs) int {
	springVariants := []string{springs.Pattern}
	result := []string{}

	for len(springVariants) > 0 {
		springVariant := springVariants[0]
		springVariants = springVariants[1:]

		isVariantComplete, isVariantValid := areSpringsValid(springVariant, springs.Groups)

		if !isVariantComplete {
			springVariants = append(springVariants, strings.Replace(springVariant, "?", ".", 1))
			springVariants = append(springVariants, strings.Replace(springVariant, "?", "#", 1))
		} else {
			if isVariantValid {
				result = append(result, springVariant)
			}
		}

	}

	return len(result)
}

func areSpringsValid(springs string, actualGroupSizes []int) (bool, bool) {
	groupSizes := []int{}
	activeGroup := false

	for _, spring := range springs {
		if spring == '.' {
			if activeGroup {
				activeGroup = false
			}
		}
		if spring == '#' {
			if !activeGroup {
				groupSizes = append(groupSizes, 1)
				activeGroup = true
			} else {
				groupSizes[len(groupSizes)-1]++
			}
		}
		if spring == '?' {
			return false, false
		}

		if !activeGroup {
			for index, groupSize := range groupSizes {
				if index >= len(actualGroupSizes) {
					return true, false
				}
				if groupSize != actualGroupSizes[index] {
					return true, false
				}
			}
		}
	}

	if len(groupSizes) != len(actualGroupSizes) {
		return true, false
	}

	for index, groupSize := range groupSizes {
		if groupSize != actualGroupSizes[index] {
			return true, false
		}
	}

	return true, true
}

func parseInput(lines []string) []Springs {
	springRows := []Springs{}

	for _, line := range lines {
		words := strings.Split(line, " ")
		groupSizesAsStrings := strings.Split(words[1], ",")
		groupSizes := []int{}

		for _, groupSizeAsString := range groupSizesAsStrings {
			groupSize, _ := strconv.Atoi(groupSizeAsString)
			groupSizes = append(groupSizes, groupSize)
		}

		springRow := Springs{
			Pattern: words[0],
			Groups:  groupSizes,
		}

		springRows = append(springRows, springRow)
	}

	return springRows
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
