package day04

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const INPUT_FILE_PATH = "day04/input.txt"

type Card struct {
	Id      int
	Winning []int
	Present []int
}

func Run() {
	inputLines := readLines(INPUT_FILE_PATH)
	cards := parseInputLines(inputLines)
	cardCopyCount := make(map[int]int)

	totalScore := 0
	totalCards := 0

	for _, card := range cards {
		cardCopyCount[card.Id] += 1
		score := calculateScore(card)
		intersectionSize := intersectionSize(card)

		totalScore += score
		for i := 1; i <= intersectionSize; i++ {
			cardCopyCount[card.Id + i] += cardCopyCount[card.Id]
		}
	}

	for _, copyCount := range cardCopyCount {
		totalCards += copyCount
	}

	fmt.Println("Total score [PART 1]: ", totalScore)
	fmt.Println("Total cards [PART 2]: ", totalCards)
}

func intersectionSize(card *Card) int {
	intersection := make(map[int]int)
	size := 0

	for _, number := range card.Present {
		intersection[number] += 1
	}
	for _, number := range card.Winning {
		intersection[number] += 1
	}

	for _, occurrence := range intersection {
		if occurrence == 2 {
			size += 1
		}
	}

	return size
}

func calculateScore(card *Card) int {
	occurrences := make(map[int]int)
	score := 0

	for _, number := range card.Present {
		occurrences[number] += 1
	}
	for _, number := range card.Winning {
		occurrences[number] += 1
	}

	for _, occurrence := range occurrences {
		if occurrence == 2 {
			if score == 0 {
				score = 1
			} else {
				score *= 2
			}
		}
	}

	return score
}

func parseInputLines(inputLines []string) []*Card {
	var cards []*Card

	for _, inputLine := range inputLines {
		card, err := parseCardLine(inputLine)
		if err != nil {
			panic(err)
		}

		cards = append(cards, card)
	}

	return cards
}

func parseCardLine(line string) (*Card, error) {
	prefixRegex := regexp.MustCompile(`^Card\s+(\d+): `)
	matches := prefixRegex.FindStringSubmatch(line)
	var card Card

	if len(matches) != 2 {
		return nil, errors.New("invalid card string: wrong prefix")
	}

	prefix := matches[0]
	cardId, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, errors.New("invalid card string: non-integer card ID")
	}
	card.Id = cardId

	parsedString, _ := strings.CutPrefix(line, prefix)
	numberSetStrings := strings.Split(parsedString, " | ")

	if len(numberSetStrings) != 2 {
		return nil, errors.New("invalid card string: no numbers separator")
	}

	for _, numberString := range strings.Split(numberSetStrings[0], " ") {
		number, err := strconv.Atoi(numberString)
		if err == nil {
			card.Winning = append(card.Winning, number)
		}
	}

	for _, numberString := range strings.Split(numberSetStrings[1], " ") {
		number, err := strconv.Atoi(numberString)
		if err == nil {
			card.Present = append(card.Present, number)
		}
	}

	return &card, nil
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
