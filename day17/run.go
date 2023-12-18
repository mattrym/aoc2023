package day17

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"

	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

const INPUT_FILE_PATH = "day17/input.txt"

type Direction int

const (
	NONE  Direction = iota
	UP    Direction = iota
	RIGHT Direction = iota
	DOWN  Direction = iota
	LEFT  Direction = iota
)

type Roadmap struct {
	Grid      [][]int64
	DimX      int
	DimY      int
	Predicate func(Step, Step) bool
}

type Position struct {
	X int
	Y int
}

type Step struct {
	Position  Position
	Direction Direction
	Stride    int
}

func Run() {
	lines := readLines(INPUT_FILE_PATH)
	roadmap, err := parseInput(lines)
	if err != nil {
		panic(err)
	}

	roadmap.Predicate = func(oldStep Step, step Step) bool {
		return step.Stride < 4
	}
	sinkDistancePart1 := roadmap.DijkstraShortestPathLength()
	fmt.Println("Shortest path distance [PART 1]: ", sinkDistancePart1)

	roadmap.Predicate = func(oldStep Step, step Step) bool {
		if oldStep.Stride < 4 {
			return step.Direction == oldStep.Direction
		}
		return step.Stride < 11
	}

	sinkDistancePart2 := roadmap.DijkstraShortestPathLength()
	fmt.Println("Shortest path distance [PART 2]: ", sinkDistancePart2)
}

func (roadmap *Roadmap) DijkstraShortestPathLength() int64 {
	visited := make(map[Step]bool, roadmap.DimY)
	queue := pq.New[Step, int64](pq.MinHeap)

	queue.Put(Step{Position{1, 0}, RIGHT, 1}, roadmap.Grid[0][1])
	queue.Put(Step{Position{0, 1}, DOWN, 1}, roadmap.Grid[1][0])

	for !queue.IsEmpty() {
		item := queue.Get()
		step, distance := item.Value, item.Priority

		if nodeVisited, ok := visited[step]; !ok || !nodeVisited {
			visited[step] = true
		} else {
			continue
		}

		if roadmap.EndingConditionMet(step, step) {
			return distance
		}

		for _, nextStep := range roadmap.AllowedSteps(step) {
			queue.Put(nextStep, roadmap.IncDistance(nextStep, distance))
		}
	}

	return math.MaxInt64
}

func (roadmap *Roadmap) EndingConditionMet(oldStep Step, step Step) bool {
	if step.Position.X != roadmap.DimX-1 ||
		step.Position.Y != roadmap.DimY-1 {
		return false
	}

	if !roadmap.Predicate(oldStep, step) {
		return false
	}

	return true
}

func (roadmap *Roadmap) IncDistance(step Step, distance int64) int64 {
	return distance + roadmap.Grid[step.Position.Y][step.Position.X]
}

func (roadmap *Roadmap) AllowedSteps(step Step) []Step {
	allowedDirections := roadmap.clampDirections(step.Position)
	allowedDirections[oppositeDirection(step.Direction)] = false
	allowedSteps := make([]Step, 0)

	for direction, allowed := range allowedDirections {
		if allowed {
			newStep := Step{
				Position:  Move(step.Position, direction),
				Direction: direction,
				Stride:    1,
			}

			if direction == step.Direction {
				newStep.Stride = step.Stride + 1
			}

			if !roadmap.Predicate(step, newStep) {
				allowed = false
			}

			if allowed {
				allowedSteps = append(allowedSteps, newStep)
			}
		}
	}

	return allowedSteps
}

func (roadmap *Roadmap) clampDirections(position Position) map[Direction]bool {
	allowedDirections := make(map[Direction]bool)

	if position.Y == 0 {
		allowedDirections[UP] = false
	} else {
		allowedDirections[UP] = true
	}

	if position.X == roadmap.DimX-1 {
		allowedDirections[RIGHT] = false
	} else {
		allowedDirections[RIGHT] = true
	}

	if position.Y == roadmap.DimY-1 {
		allowedDirections[DOWN] = false
	} else {
		allowedDirections[DOWN] = true
	}

	if position.X == 0 {
		allowedDirections[LEFT] = false
	} else {
		allowedDirections[LEFT] = true
	}

	return allowedDirections
}

func (roadmap *Roadmap) AllowedDirections(node Step) {
	allowedDirections := make(map[Direction]bool)

	if node.Position.Y > 0 {
		allowedDirections[UP] = true
	}

	if node.Position.X < roadmap.DimX-1 {
		allowedDirections[RIGHT] = true
	}

	if node.Position.Y < roadmap.DimY-1 {
		allowedDirections[DOWN] = true
	}

	if node.Position.X > 0 {
		allowedDirections[LEFT] = true
	}

	allDirections := []Direction{UP, RIGHT, DOWN, LEFT}

	if node.Stride < 4 {
		for _, direction := range allDirections {
			if node.Direction != direction {
				allowedDirections[direction] = false
			}
		}
	}

	if node.Stride >= 10 {
		allowedDirections[node.Direction] = false
	}

	allowedDirections[oppositeDirection(node.Direction)] = false

}

func Move(position Position, direction Direction) Position {
	switch direction {
	case UP:
		position.Y--
	case RIGHT:
		position.X++
	case DOWN:
		position.Y++
	case LEFT:
		position.X--
	}
	return position
}

func parseInput(lines []string) (*Roadmap, error) {
	roadmap := Roadmap{}
	roadmap.DimX = len(lines[0])

	for y, line := range lines {
		if line == "" {
			break
		}

		roadmap.Grid = append(roadmap.Grid, make([]int64, roadmap.DimX))

		for x := 0; x < len(line); x++ {
			heatLoss, err := strconv.Atoi(line[x : x+1])
			if err != nil {
				return nil, err
			}

			roadmap.Grid[y][x] = int64(heatLoss)
		}

		roadmap.DimY++
	}

	return &roadmap, nil
}

func oppositeDirection(direction Direction) Direction {
	switch direction {
	case UP:
		return DOWN
	case RIGHT:
		return LEFT
	case DOWN:
		return UP
	case LEFT:
		return RIGHT
	}
	return 0
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
