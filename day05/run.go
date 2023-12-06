package day05

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Range struct {
	Start  int
	Length int
}

type Ranges []Range

func (r1 Range) Subtract(r2 Range) []Range {
	if r1.Start > r2.Start {
		r1, r2 = r2, r1
	}

	if r1.Start+r1.Length <= r2.Start {
		return []Range{r1}
	}

	if r1.Start+r1.Length <= r2.Start+r2.Length {
		if (r2.Start - r1.Start) == 0 {
			return []Range{}
		}
		return []Range{{Start: r1.Start, Length: r2.Start - r1.Start}}
	}

	ranges := []Range{}
	if (r2.Start - r1.Start) != 0 {
		ranges = append(ranges, Range{Start: r1.Start, Length: r2.Start - r1.Start})
	}
	if (r1.Start + r1.Length - r2.Start - r2.Length) != 0 {
		ranges = append(ranges, Range{
			Start:  r2.Start + r2.Length,
			Length: r1.Start + r1.Length - r2.Start - r2.Length,
		})
	}
	return ranges
}

func (ranges Ranges) Deduplicate() Ranges {
	allRangesMap := make(map[Range]bool)
	rangesDeduped := []Range{}

	for _, range_ := range ranges {
		if _, ok := allRangesMap[range_]; !ok {
			rangesDeduped = append(rangesDeduped, range_)
			allRangesMap[range_] = true
		}
	}

	return rangesDeduped
}

func Intersection(r1 Range, r2 Range) (bool, Range) {
	if r1.Start > r2.Start {
		r1, r2 = r2, r1
	}

	if r1.Start+r1.Length <= r2.Start {
		return false, Range{}
	}

	if r1.Start+r1.Length <= r2.Start+r2.Length {
		return true, Range{Start: r2.Start, Length: r1.Start + r1.Length - r2.Start}
	}

	return true, Range{Start: r2.Start, Length: r2.Length}
}

type Mapping struct {
	Diff        int
	SourceRange Range
}

type CategoryMap struct {
	FromCategory string
	ToCategory   string
	Mappings     []Mapping
}

type Almanac struct {
	CategoryMaps map[string]*CategoryMap
}

const INPUT_FILE_PATH = "day05/input_test.txt"

func Run() {
	inputLines := readLines(INPUT_FILE_PATH)
	almanac, err := parseInput(inputLines)
	if err != nil {
		panic(err)
	}

	seeds, err := parseSeedsAsSingleNumbers(inputLines[0])
	if err != nil {
		panic(err)
	}

	seedRanges, err := parseSeedsAsRanges(inputLines[0])
	if err != nil {
		panic(err)
	}

	var startTime time.Time
	var elapsedTime time.Duration

	startTime = time.Now()
	result1 := almanac.MinimalLocationFromSeeds(seeds)
	elapsedTime = time.Since(startTime)

	fmt.Println("Minimal converted seed [PART 1]: ", result1, ";\tElapsed time: ", elapsedTime)

	startTime = time.Now()
	result2 := almanac.OptimalMinimalLocationFromSeedRanges(seedRanges)
	elapsedTime = time.Since(startTime)

	fmt.Println("Minimal converted seed [PART 2]: ", result2, ";\tElapsed time: ", elapsedTime)
}

func (almanac *Almanac) MinimalLocationFromSeeds(seeds []int) int {
	var locations []int
	var minLocation int

	for _, seed := range seeds {
		locations = append(locations, almanac.ConvertSeedFromInt(seed))
	}
	minLocation = math.MaxInt

	for _, seed := range locations {
		if seed < minLocation {
			minLocation = seed
		}
	}

	return minLocation
}

func (almanac *Almanac) BruteForceMinimalLocationFromSeedRanges(seedRanges []Range) int {
	minLocation := math.MaxInt

	for _, seedRange := range seedRanges {
		for i := 0; i < seedRange.Length; i++ {
			seed := almanac.ConvertSeedFromInt(seedRange.Start + i)
			if seed < minLocation {
				minLocation = seed
			}
		}
	}

	return minLocation
}

func (almanac *Almanac) OptimalMinimalLocationFromSeedRanges(seedRanges []Range) int {
	var ranges Ranges

	for _, seedRange := range seedRanges {
		newRanges := almanac.ConvertSeedFromRange(seedRange)
		ranges = append(ranges, newRanges...)
	}

	minRangeStart := math.MaxInt
	for _, range_ := range ranges {
		if range_.Start < minRangeStart {
			minRangeStart = range_.Start
		}
	}
	return minRangeStart

}

func (almanac *Almanac) ConvertSeedFromInt(initNumber int) int {
	category := "seed"
	number := initNumber

	for {
		categoryMap := almanac.CategoryMaps[category]
		if categoryMap == nil {
			return number
		}

		category = categoryMap.ToCategory
		number = categoryMap.ConvertInt(number)
	}
}

func (almanac *Almanac) ConvertSeedFromRange(initRange Range) []Range {
	category := "seed"
	result := []Range{initRange}

	for {
		categoryMap := almanac.CategoryMaps[category]
		if categoryMap == nil {
			return result
		}

		var newRanges []Range
		category = categoryMap.ToCategory
		for _, range_ := range result {
			convertedRanges := categoryMap.ConvertRange(range_)
			newRanges = append(newRanges, convertedRanges...)
		}
		result = newRanges
	}
}

func (categoryMap *CategoryMap) ConvertRange(_range Range) []Range {
	result := Ranges{}
	ranges := Ranges{_range}

	for len(ranges) > 0 {
		newResult := Ranges{}
		range_ := ranges[0]
		ranges = ranges[1:]

		for _, mapping := range categoryMap.Mappings {
			if ok, destRange := Intersection(mapping.SourceRange, range_); ok {
				newResult = append(newResult, Range{Start: destRange.Start + mapping.Diff, Length: destRange.Length})
				ranges = append(ranges, range_.Subtract(destRange)...)
			}
		}

		if len(newResult) > 0 {
			result = append(result, newResult...)
			result = result.Deduplicate()
			ranges = ranges.Deduplicate()
		} else {
			result = append(result, range_)
		}
	}

	return result
}

func (categoryMap *CategoryMap) ConvertInt(source int) int {
	for _, mapping := range categoryMap.Mappings {
		if ok, dest := mapping.Convert(source); ok {
			return dest
		}
	}

	return source
}

func (mapping *Mapping) Convert(source int) (bool, int) {
	if source >= mapping.SourceRange.Start && source < mapping.SourceRange.Start+mapping.SourceRange.Length {
		return true, source + mapping.Diff
	}
	return false, 0
}

func parseInput(lines []string) (*Almanac, error) {
	result := Almanac{}
	result.CategoryMaps = make(map[string]*CategoryMap)

	for lineIndex := 1; lineIndex < len(lines); {
		if lines[lineIndex] == "" {
			lineIndex++
			continue
		}

		if isMapPrefix(lines[lineIndex]) {
			categoryMap, newIndex, err := parseMap(&lines, lineIndex)
			if err != nil {
				return nil, err
			}

			result.CategoryMaps[categoryMap.FromCategory] = categoryMap
			lineIndex = newIndex + 1
			continue
		}
	}

	return &result, nil
}

func isMapPrefix(line string) bool {
	prefixRegex := regexp.MustCompile(`^(\w+)\-to\-(\w+)\s+map:\s*$`)
	return prefixRegex.MatchString(line)
}

func parseSeedsAsSingleNumbers(line string) ([]int, error) {
	line = strings.TrimPrefix(line, "seeds: ")
	numbersAsString := strings.Split(line, " ")
	var result []int

	for _, numberAsString := range numbersAsString {
		number, err := strconv.Atoi(numberAsString)
		if err != nil {
			return nil, fmt.Errorf("invalid seed number: %s", numberAsString)
		}

		result = append(result, number)
	}

	return result, nil
}

func parseSeedsAsRanges(line string) ([]Range, error) {
	line = strings.TrimPrefix(line, "seeds: ")
	numbersAsString := strings.Split(line, " ")
	var result []Range

	for i := 0; i < len(numbersAsString); i += 2 {
		start, err := strconv.Atoi(numbersAsString[i])
		if err != nil {
			return nil, fmt.Errorf("invalid seed range start: %s", numbersAsString[i])
		}

		length, err := strconv.Atoi(numbersAsString[i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid seed range length: %s", numbersAsString[i+1])
		}

		result = append(result, Range{Start: start, Length: length})
	}

	return result, nil
}

func parseMap(lines *[]string, initIndex int) (*CategoryMap, int, error) {
	prefixRegex := regexp.MustCompile(`^(\w+)\-to\-(\w+)\s+map:\s*$`)
	matches := prefixRegex.FindStringSubmatch((*lines)[initIndex])
	result := CategoryMap{}

	if len(matches) != 3 {
		return nil, initIndex, fmt.Errorf("invalid map line (mapping prefix): %s", (*lines)[initIndex])
	}

	result.FromCategory = matches[1]
	result.ToCategory = matches[2]
	index := initIndex + 1

	for ; index < len(*lines); index++ {
		mapping := Mapping{}
		line := (*lines)[index]

		if line == "" {
			break
		}

		numbersAsString := strings.Split(line, " ")
		if len(numbersAsString) != 3 {
			return nil, index, fmt.Errorf("invalid map line: %s (number count)", line)
		}

		destFrom, err := strconv.Atoi(numbersAsString[0])
		if err != nil {
			return nil, index, fmt.Errorf("invalid map line: %s (destination start)", line)
		}

		sourceFrom, err := strconv.Atoi(numbersAsString[1])
		if err != nil {
			return nil, index, fmt.Errorf("invalid map line: %s (source start)", line)
		}

		length, err := strconv.Atoi(numbersAsString[2])
		if err != nil {
			return nil, index, fmt.Errorf("invalid map line: %s (range length)", line)
		}

		mapping.Diff = destFrom - sourceFrom
		mapping.SourceRange = Range{Start: sourceFrom, Length: length}
		result.Mappings = append(result.Mappings, mapping)
	}

	return &result, index, nil
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
