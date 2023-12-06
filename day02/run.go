package day02

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Game struct {
	Id       int
	CubeSets []CubeSet
}

type CubeSet struct {
	RedCubes   int
	GreenCubes int
	BlueCubes  int
}

const INPUT_FILE_PATH = "day02/input.txt"
const MAX_RED_CUBES = 12
const MAX_GREEN_CUBES = 13
const MAX_BLUE_CUBES = 14

func Run() {
	// Read input file (from input.txt)
	var sumOfPossibleGameIds, sumOfPossibleGamePowers int
	inputLines := readLines(INPUT_FILE_PATH)
	games := make([]*Game, len(inputLines))

	for _, inputLine := range inputLines {
		game, err := ParseGameString(inputLine)
		if err != nil {
			panic(err)
		}

		games[game.Id-1] = game
	}

	for _, game := range games {
		minCubeSet := MinGameCubeSet(game)
		minCubeSetPower := CubeSetPower(minCubeSet)

		if IsGamePossible(game) {
			sumOfPossibleGameIds += game.Id
		}
		sumOfPossibleGamePowers += minCubeSetPower
	}

	// Print result
	println("Sum of possible game IDs 	 [PART1]: ", sumOfPossibleGameIds)
	println("Sum of possible game powers [PART2]: ", sumOfPossibleGamePowers)
}

func CubeSetPower(cubeSet *CubeSet) int {
	return cubeSet.RedCubes * cubeSet.GreenCubes * cubeSet.BlueCubes
}

func MinGameCubeSet(game *Game) *CubeSet {
	minCubeSet := CubeSet{
		RedCubes:   0,
		GreenCubes: 0,
		BlueCubes:  0,
	}

	for _, cubeSet := range game.CubeSets {
		if cubeSet.RedCubes > minCubeSet.RedCubes {
			minCubeSet.RedCubes = cubeSet.RedCubes
		}
		if cubeSet.GreenCubes > minCubeSet.GreenCubes {
			minCubeSet.GreenCubes = cubeSet.GreenCubes
		}
		if cubeSet.BlueCubes > minCubeSet.BlueCubes {
			minCubeSet.BlueCubes = cubeSet.BlueCubes
		}
	}

	fmt.Println(game)
	fmt.Println(minCubeSet)

	return &minCubeSet
}

func IsGamePossible(game *Game) bool {
	for _, cubeSet := range game.CubeSets {
		if cubeSet.RedCubes > MAX_RED_CUBES ||
			cubeSet.GreenCubes > MAX_GREEN_CUBES ||
			cubeSet.BlueCubes > MAX_BLUE_CUBES {
			return false
		}
	}

	return true
}

func ParseGameString(gameString string) (*Game, error) {
	prefixRegex := regexp.MustCompile(`^Game (\d+): `)
	matches := prefixRegex.FindStringSubmatch(gameString)

	var cubeSets []CubeSet
	var prefix string

	if len(matches) != 2 {
		return nil, errors.New("invalid game string")
	}

	prefix = matches[0]
	gameId, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, errors.New("invalid game string")
	}

	parsedString, _ := strings.CutPrefix(gameString, prefix)
	cubeSetsStrings := strings.Split(parsedString, "; ")

	for _, cubeSetString := range cubeSetsStrings {
		cubeStrings := strings.Split(cubeSetString, ", ")
		cubeSet := CubeSet{}

		for _, cubeString := range cubeStrings {
			countAndColor := strings.Split(cubeString, " ")
			if len(countAndColor) != 2 {
				return nil, errors.New("invalid game string")
			}

			count, err := strconv.Atoi(countAndColor[0])
			color := countAndColor[1]
			if err != nil {
				return nil, errors.New("invalid game string")
			}

			switch color {
			case "red":
				cubeSet.RedCubes = count
			case "green":
				cubeSet.GreenCubes = count
			case "blue":
				cubeSet.BlueCubes = count
			default:
				return nil, errors.New("invalid game string")
			}
		}

		cubeSets = append(cubeSets, cubeSet)
	}

	return &Game{
		Id:       gameId,
		CubeSets: cubeSets,
	}, nil
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
