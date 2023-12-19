package day19

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Category int

const (
	X Category = iota
	M
	A
	S
)

type Part map[Category]int

type Condition struct {
	Category Category
	Operator string
	Number   int
	Effect   string
}

type Workflow struct {
	Name       string
	Conditions []Condition
	LastEffect string
}

type System map[string]*Workflow

const MIN_FROM = 1
const MAX_TO = 4000

type Interval struct {
	From int
	To   int
}

type Range map[Category]Interval
type Results map[string][]Range

const INPUT_FILE_PATH = "day19/input.txt"

func Run() {
	inputLines := readLines(INPUT_FILE_PATH)
	system, parts, err := parseInput(inputLines)
	if err != nil {
		panic(err)
	}

	sumOfRatings := 0
	for _, part := range parts {
		if system.EvaluatePart(part) {
			sumOfRatings += part.TotalRating()
		}
	}

	fmt.Println("Sum of ratings [PART 1]: ", sumOfRatings)

	cachedRanges := make(Results)
	acceptedRanges := system.EvaluateRange(cachedRanges, "in")

	totalCombinations := 0
	for _, acceptedRange := range acceptedRanges {
		totalCombinations += acceptedRange.Combinations()
	}

	fmt.Println("Total combinations [PART 2]: ", totalCombinations)
}

func (range_ Range) String() string {
	if range_.IsEmpty() {
		return "∅"
	}

	parts := []string{}

	for _, category := range []Category{X, M, A, S} {
		char := categoryChar(category)
		from := range_[category].From
		to := range_[category].To

		part := fmt.Sprintf("%s: [%d,%d]", char, from, to)
		parts = append(parts, part)
	}

	return "{" + strings.Join(parts, ", ") + "}"
}

func categoryChar(category Category) string {
	switch category {
	case X:
		return "x"
	case M:
		return "m"
	case A:
		return "a"
	case S:
		return "s"
	}
	return ""
}

func (range_ Range) Combinations() int {
	combinations := 1
	for _, category := range []Category{X, M, A, S} {
		if range_[category].From == 0 {
			return 0
		}
		combinations *= range_[category].To - range_[category].From + 1
	}
	return combinations
}

func (range_ Range) IsEmpty() bool {
	for _, category := range []Category{X, M, A, S} {
		_, ok := range_[category]
		if !ok {
			return true
		}

		if range_[category].From == 0 || range_[category].To == 0 {
			return true
		}
	}
	return false
}

func (part Part) TotalRating() int {
	totalRating := 0
	for _, value := range part {
		totalRating += value
	}
	return totalRating
}

func (system System) EvaluateRange(cachedRanges Results, workflowName string) []Range {
	if workflowName == "ACCEPT" {
		return []Range{FullRange()}
	}
	if workflowName == "REJECT" {
		return []Range{}
	}
	if cachedRange, ok := cachedRanges[workflowName]; ok {
		return cachedRange
	}

	workflow := system[workflowName]
	acceptedRanges := make([]Range, 0)
	analyzedRange := FullRange()

	for _, condition := range workflow.Conditions {
		trueConditionRange, falseConditionRange := condition.SplitRange(analyzedRange)
		nextWorkflowRanges := system.EvaluateRange(cachedRanges, condition.Effect)
		intersectedRanges := IntersectMany(trueConditionRange, nextWorkflowRanges)

		acceptedRanges = append(acceptedRanges, intersectedRanges...)
		analyzedRange = falseConditionRange
	}

	nextWorkflowRanges := system.EvaluateRange(cachedRanges, workflow.LastEffect)
	intersectedRange := IntersectMany(analyzedRange, nextWorkflowRanges)
	acceptedRanges = append(acceptedRanges, intersectedRange...)
	newAcceptedRanges := []Range{}

	for _, acceptedRange := range acceptedRanges {
		if !acceptedRange.IsEmpty() {
			newAcceptedRanges = append(newAcceptedRanges, acceptedRange)
		}
	}

	cachedRanges[workflowName] = newAcceptedRanges
	return acceptedRanges
}

func FullRange() Range {
	partRange := make(Range)
	for _, category := range []Category{X, M, A, S} {
		partRange[category] = Interval{MIN_FROM, MAX_TO}
	}
	return partRange
}

func EmptyRange() Range {
	return map[Category]Interval{}
}

func (condition *Condition) SplitRange(range_ Range) (Range, Range) {
	fullPartRange := FullRange()
	trueCondRange := maps.Clone(fullPartRange)
	falseCondRange := maps.Clone(fullPartRange)

	switch condition.Operator {
	case "<":
		trueCondRange[condition.Category] = Interval{MIN_FROM, condition.Number - 1}
		falseCondRange[condition.Category] = Interval{condition.Number, MAX_TO}
	case ">":
		trueCondRange[condition.Category] = Interval{condition.Number + 1, MAX_TO}
		falseCondRange[condition.Category] = Interval{MIN_FROM, condition.Number}
	}

	return Intersect(range_, trueCondRange), Intersect(range_, falseCondRange)
}

func IntersectMany(pr Range, prs []Range) []Range {
	result := []Range{}

	for _, pr2 := range prs {
		resultPr := Intersect(pr, pr2)
		result = append(result, resultPr)
	}

	return result
}

func Intersect(pr1 Range, pr2 Range) Range {
	if pr1.IsEmpty() || pr2.IsEmpty() {
		return EmptyRange()
	}

	pr := make(Range)
	var ok bool

	for _, category := range []Category{X, M, A, S} {
		ok, pr[category] = Intersection(pr1[category], pr2[category])
		if !ok {
			return map[Category]Interval{}
		}
	}

	return pr
}

func Intersection(r1 Interval, r2 Interval) (bool, Interval) {
	if r1.From > r2.From {
		r1, r2 = r2, r1
	}

	if r1.To <= r2.From {
		return false, Interval{}
	}

	if r1.To <= r2.To {
		return true, Interval{r2.From, r1.To}
	}

	return true, Interval{r2.From, r2.To}
}

func (system System) EvaluatePart(part Part) bool {
	workflowName := "in"

	for {
		workflow := system[workflowName]
		workflowName = workflow.EvaluatePart(part)
		if workflowName == "ACCEPT" {
			return true
		}
		if workflowName == "REJECT" {
			return false
		}
	}
}

func (workflow *Workflow) EvaluatePart(part Part) string {
	for _, condition := range workflow.Conditions {
		effect := condition.EvaluatePart(part)
		if effect != "" {
			return effect
		}
	}
	return workflow.LastEffect
}

func (condition *Condition) EvaluatePart(part Part) string {
	partValue := part[condition.Category]
	switch condition.Operator {
	case "<":
		if partValue < condition.Number {
			return condition.Effect
		}
	case ">":
		if partValue > condition.Number {
			return condition.Effect
		}
	}
	return ""
}

func parseInput(lines []string) (*System, []Part, error) {
	system := make(System)
	parts := []Part{}
	i := 0

	for i = 0; i < len(lines); i++ {
		if lines[i] == "" {
			break
		}

		workflow, err := parseWorkflow(lines[i])
		if err != nil {
			return nil, nil, err
		}
		system[workflow.Name] = workflow
	}

	for ; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}

		part, err := parsePart(lines[i])
		if err != nil {
			return nil, nil, err
		}

		parts = append(parts, *part)
	}

	return &system, parts, nil
}

func parseWorkflow(line string) (*Workflow, error) {
	parts := strings.Split(line, "{")
	if len(parts) != 2 {
		return nil, errors.New("invalid workflow")
	}
	clausesString := strings.TrimSuffix(parts[1], "}")
	clauses := strings.Split(clausesString, ",")

	workflow := Workflow{}
	workflow.Name = parts[0]
	workflow.LastEffect = clauses[len(clauses)-1]

	switch workflow.LastEffect {
	case "A":
		workflow.LastEffect = "ACCEPT"
	case "R":
		workflow.LastEffect = "REJECT"
	}

	for _, clause := range clauses[:len(clauses)-1] {
		conditionRegex := regexp.MustCompile(`^([xmas])([<>])(\d+):(\w+)$`)
		conditionParts := conditionRegex.FindStringSubmatch(clause)
		if len(conditionParts) != 5 {
			return nil, errors.New("invalid condition: " + clause)
		}

		category, err := stringToCategory(conditionParts[1])
		if err != nil {
			return nil, errors.New("invalid condition category: " + clause)
		}

		number, err := strconv.Atoi(conditionParts[3])
		if err != nil {
			return nil, errors.New("invalid condition number: " + clause)
		}

		effect := conditionParts[4]
		switch effect {
		case "A":
			effect = "ACCEPT"
		case "R":
			effect = "REJECT"
		}

		condition := Condition{
			Category: category,
			Operator: conditionParts[2][:1],
			Number:   number,
			Effect:   effect,
		}
		workflow.Conditions = append(workflow.Conditions, condition)
	}

	return &workflow, nil
}

func parsePart(line string) (*Part, error) {
	line = strings.TrimPrefix(line, "{")
	line = strings.TrimSuffix(line, "}")
	categoryStrings := strings.Split(line, ",")

	part := Part{}
	for _, categoryString := range categoryStrings {
		categoryRegex := regexp.MustCompile(`^([xmas])=(\d+)$`)
		categoryParts := categoryRegex.FindStringSubmatch(categoryString)
		if len(categoryParts) != 3 {
			return nil, errors.New("invalid category: " + categoryString)
		}

		category, err := stringToCategory(categoryParts[1])
		if err != nil {
			return nil, errors.New("invalid category: " + categoryString)
		}

		number, err := strconv.Atoi(categoryParts[2])
		if err != nil {
			return nil, errors.New("invalid category number: " + categoryString)
		}

		part[category] = number
	}

	return &part, nil
}

func stringToCategory(categoryString string) (Category, error) {
	switch categoryString {
	case "x":
		return X, nil
	case "m":
		return M, nil
	case "a":
		return A, nil
	case "s":
		return S, nil
	default:
		return -1, errors.New("invalid category: " + categoryString)
	}
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
