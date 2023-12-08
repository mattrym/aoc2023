package day07

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const INPUT_FILE_PATH = "day07/input.txt"
const CARDS_ORDER_PART_1 = "23456789TJQKA"
const CARDS_ORDER_PART_2 = "J23456789TQKA"

const FIVE_OF_A_KIND = 7
const FOUR_OF_A_KIND = 6
const FULL_HOUSE = 5
const THREE_OF_A_KIND = 4
const TWO_PAIRS = 3
const ONE_PAIR = 2
const HIGH_CARD = 1

type Hand struct {
	Cards []byte
	Bid   int
}

type HandsPart1 []Hand
type HandsPart2 []Hand

func (h *HandsPart1) Len() int {
	return len(*h)
}

func (h *HandsPart2) Len() int {
	return len(*h)
}

func (h *HandsPart1) Less(i, j int) bool {
	handValueI := (*h)[i].ValuePart1()
	handValueJ := (*h)[j].ValuePart1()

	if handValueI != handValueJ {
		return handValueI < handValueJ
	}

	for k := 0; k < 5; k++ {
		cardValueI := CardValuePart1((*h)[i].Cards[k])
		cardValueJ := CardValuePart1((*h)[j].Cards[k])

		if cardValueI != cardValueJ {
			return cardValueI < cardValueJ
		}
	}
	return false
}

func (h *HandsPart2) Less(i, j int) bool {
	handValueI := (*h)[i].ValuePart2()
	handValueJ := (*h)[j].ValuePart2()

	if handValueI != handValueJ {
		return handValueI < handValueJ
	}

	for k := 0; k < 5; k++ {
		cardValueI := CardValuePart2((*h)[i].Cards[k])
		cardValueJ := CardValuePart2((*h)[j].Cards[k])

		if cardValueI != cardValueJ {
			return cardValueI < cardValueJ
		}
	}
	return false
}

func (h *HandsPart1) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *HandsPart2) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}
func (hand Hand) ValuePart1() int {
	cardMap := make(map[byte]int)
	for _, card := range hand.Cards {
		cardMap[card]++
	}

	values := make([]int, 0)
	for _, card := range CARDS_ORDER_PART_1 {
		values = append(values, cardMap[byte(card)])
	}

	multiples := make(map[int]int)
	for _, value := range values {
		multiples[value]++
	}

	if multiples[5] == 1 {
		return FIVE_OF_A_KIND
	}
	if multiples[4] == 1 {
		return FOUR_OF_A_KIND
	}
	if multiples[3] == 1 && multiples[2] == 1 {
		return FULL_HOUSE
	}
	if multiples[3] == 1 {
		return THREE_OF_A_KIND
	}
	if multiples[2] == 2 {
		return TWO_PAIRS
	}
	if multiples[2] == 1 {
		return ONE_PAIR
	}
	return HIGH_CARD
}

func (hand Hand) ValuePart2() int {
	cardMap := make(map[byte]int)
	for _, card := range hand.Cards {
		cardMap[card]++
	}

	values := make([]int, 0)
	for _, card := range CARDS_ORDER_PART_2[1:] {
		values = append(values, cardMap[byte(card)])
	}

	multiples := make(map[int]int)
	for _, value := range values {
		multiples[value]++
	}

	if cardMap[byte('J')] == 5 {
		return FIVE_OF_A_KIND
	}
	if cardMap[byte('J')] == 4 {
		return FIVE_OF_A_KIND
	}
	if cardMap[byte('J')] == 3 && multiples[2] == 1 {
		return FIVE_OF_A_KIND
	}
	if cardMap[byte('J')] == 3 {
		return FOUR_OF_A_KIND
	}
	if cardMap[byte('J')] == 2 && multiples[3] == 1 {
		return FIVE_OF_A_KIND
	}
	if cardMap[byte('J')] == 2 && multiples[2] == 1 {
		return FOUR_OF_A_KIND
	}
	if cardMap[byte('J')] == 2 {
		return THREE_OF_A_KIND
	}
	if cardMap[byte('J')] == 1 && multiples[4] == 1 {
		return FIVE_OF_A_KIND
	}
	if cardMap[byte('J')] == 1 && multiples[3] == 1 {
		return FOUR_OF_A_KIND
	}
	if cardMap[byte('J')] == 1 && multiples[2] == 2 {
		return FULL_HOUSE
	}
	if cardMap[byte('J')] == 1 && multiples[2] == 1 {
		return THREE_OF_A_KIND
	}
	if cardMap[byte('J')] == 1 {
		return ONE_PAIR
	}

	if multiples[5] > 0 {
		return FIVE_OF_A_KIND
	}
	if multiples[4] > 0 {
		return FOUR_OF_A_KIND
	}
	if multiples[3] == 1 && multiples[2] == 1 {
		return FULL_HOUSE
	}
	if multiples[3] == 1 {
		return THREE_OF_A_KIND
	}
	if multiples[2] == 2 {
		return TWO_PAIRS
	}
	if multiples[2] == 1 {
		return ONE_PAIR
	}
	return HIGH_CARD
}

func CardValuePart1(card byte) int {
	for i, c := range CARDS_ORDER_PART_1 {
		if byte(c) == card {
			return i
		}
	}
	return -1
}

func CardValuePart2(card byte) int {
	for i, c := range CARDS_ORDER_PART_2 {
		if byte(c) == card {
			return i
		}
	}
	return -1
}

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	hands, err := parseInput(lines)
	if err != nil {
		panic(err)
	}

	var handsPart1 HandsPart1 = make(HandsPart1, len(hands))
	var handsPart2 HandsPart2 = make(HandsPart2, len(hands))

	copy(handsPart1, hands)
	copy(handsPart2, hands)

	sort.Sort(&handsPart1)
	sort.Sort(&handsPart2)

	totalWinningsPart1 := 0
	for rank, hand := range handsPart1 {
		winnings := hand.Bid * (rank + 1)
		totalWinningsPart1 += winnings
	}
	fmt.Println("Total winnings [PART 1]: ", totalWinningsPart1)

	totalWinningsPart2 := 0
	for rank, hand := range handsPart2 {
		winnings := hand.Bid * (rank + 1)
		totalWinningsPart2 += winnings
	}
	fmt.Println("Total winnings [PART 2]: ", totalWinningsPart2)
}

func parseInput(lines []string) ([]Hand, error) {
	var hands []Hand

	for _, line := range lines {
		hand, err := parseHand(line)
		if err != nil {
			return nil, err
		}
		hands = append(hands, *hand)
	}

	return hands, nil
}

func parseHand(line string) (*Hand, error) {
	var hand Hand

	tokens := strings.Split(line, " ")
	if len(tokens) != 2 {
		return nil, errors.New("invalid input")
	}

	cards := []byte(tokens[0])
	if len(cards) != 5 {
		return nil, errors.New("invalid input")
	}

	bid, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, errors.New("invalid input")
	}

	hand.Cards = cards
	hand.Bid = bid
	return &hand, nil
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
