package day15

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const INPUT_FILE_PATH = "day15/input.txt"

type Lense struct {
	Label       string
	FocalLength int
}
type HashMap map[int][]Lense

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	instructions := parseInput(lines)

	totalHash := 0
	for _, instruction := range instructions {
		totalHash += Hash(instruction)
	}
	fmt.Println("Total hash [PART 1]: ", totalHash)

	hashMap := InitializeHashMap(instructions)
	fmt.Println("Focus power [PART 2]: ", hashMap.FocusPower())
}

func (hashMap HashMap) FocusPower() int {
	focusPower := 0

	for boxIndex := 0; boxIndex < 256; boxIndex++ {
		numLenses := len(hashMap[boxIndex])
		for lenseIndex := 0; lenseIndex < numLenses; lenseIndex++ {
			focalLength := hashMap[boxIndex][lenseIndex].FocalLength
			focusPower += (boxIndex + 1) * (lenseIndex + 1) * focalLength
		}
	}

	return focusPower
}

func (hashMap HashMap) Add(label string, focalLength int) {
	boxIndex := Hash(label)

	for lenseIndex := 0; lenseIndex < len(hashMap[boxIndex]); lenseIndex++ {
		if hashMap[boxIndex][lenseIndex].Label == label {
			hashMap[boxIndex][lenseIndex].FocalLength = focalLength
			return
		}
	}

	hashMap[boxIndex] = append(hashMap[boxIndex], Lense{
		Label:       label,
		FocalLength: focalLength,
	})
}

func (hashMap HashMap) Remove(label string) {
	boxIndex := Hash(label)

	for lenseIndex := 0; lenseIndex < len(hashMap[boxIndex]); lenseIndex++ {
		if hashMap[boxIndex][lenseIndex].Label == label {
			hashMap[boxIndex] = append(hashMap[boxIndex][:lenseIndex], hashMap[boxIndex][lenseIndex+1:]...)
			return
		}
	}
}

func InitializeHashMap(instructions []string) *HashMap {
	hashMap := make(HashMap)
	for i := 0; i < 256; i++ {
		hashMap[i] = []Lense{}
	}

	for _, instruction := range instructions {
		if strings.Contains(instruction, "=") {
			label, focalLength, ok := parseEqualSignInstruction(instruction)
			if !ok {
				continue
			}

			hashMap.Add(label, focalLength)

		}
		if strings.Contains(instruction, "-") {
			label, ok := parseDashInstruction(instruction)
			if !ok {
				continue
			}

			hashMap.Remove(label)
		}
	}

	return &hashMap
}

func parseEqualSignInstruction(instruction string) (string, int, bool) {
	tokens := strings.Split(instruction, "=")
	if len(tokens) != 2 {
		return "", 0, false
	}

	focalLength, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", 0, false
	}

	return tokens[0], focalLength, true
}

func parseDashInstruction(instruction string) (string, bool) {
	tokens := strings.Split(instruction, "-")
	if len(tokens[0]) < 1 {
		return "", false
	}

	return tokens[0], true
}

func Hash(instruction string) int {
	bytes := []byte(instruction)
	hash := 0

	for _, b := range bytes {
		hash = (hash + int(b)) * 17 % 256
	}

	return hash
}

func parseInput(lines []string) []string {
	var result []string

	for _, line := range lines {
		instructions := strings.Split(line, ",")
		for _, instruction := range instructions {
			if instruction != "" {
				result = append(result, instruction)
			}
		}
	}

	return result
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
