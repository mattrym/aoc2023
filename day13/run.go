package day13

import (
	"bufio"
	"os"
)

const INPUT_FILE_PATH = "day13/input.txt"

func Run() {
	patterns := parsePatterns(readLines(INPUT_FILE_PATH))
	sumOfNotesPart1, sumOfNotesPart2 := 0, 0

	for _, pattern := range patterns {
		abovePart1 := pattern.FindHorizontalSymmetry()
		leftPart1 := pattern.FindVerticalSymmetry()

		abovePart2 := pattern.FindHorizontalSymmetryWithSmudge()
		leftPart2 := pattern.FindVerticalSymmetryWithSmudge()

		sumOfNotesPart1 += abovePart1*100 + leftPart1
		sumOfNotesPart2 += abovePart2*100 + leftPart2
	}

	println("Sum of notes [PART 1]: ", sumOfNotesPart1)
	println("Sum of notes [PART 2]: ", sumOfNotesPart2)
}

type Pattern struct {
	Bytes [][]byte
	DimX  int
	DimY  int
}

func (p *Pattern) Print() {
	for y := 0; y < p.DimY; y++ {
		for x := 0; x < p.DimX; x++ {
			print(string(p.Bytes[y][x]))
		}
		println()
	}
}

func (p *Pattern) Get(x, y int) byte {
	return p.Bytes[y][x]
}

func (p *Pattern) FindHorizontalSymmetry() int {
	for y := 0; y < p.DimY-1; y++ {
		symmetrical := true
		for x := 0; x < p.DimX; x++ {
			if p.Bytes[y][x] != p.Bytes[y+1][x] {
				symmetrical = false
				break
			}
		}

		if symmetrical {
			for i := 0; y-i >= 0 && y+1+i < p.DimY; i++ {
				for x := 0; x < p.DimX; x++ {
					if p.Bytes[y-i][x] != p.Bytes[y+1+i][x] {
						symmetrical = false
						break
					}
				}
			}
		}

		if symmetrical {
			return y + 1
		}
	}
	return 0
}

func (p *Pattern) FindVerticalSymmetry() int {
	for x := 0; x < p.DimX-1; x++ {
		symmetrical := true
		for y := 0; y < p.DimY; y++ {
			if p.Bytes[y][x] != p.Bytes[y][x+1] {
				symmetrical = false
				break
			}
		}

		if symmetrical {
			for i := 0; x-i >= 0 && x+1+i < p.DimX; i++ {
				for y := 0; y < p.DimY; y++ {
					if p.Bytes[y][x-i] != p.Bytes[y][x+1+i] {
						symmetrical = false
						break
					}
				}
			}
		}

		if symmetrical {
			return x + 1
		}
	}
	return 0
}

func (p *Pattern) FindHorizontalSymmetryWithSmudge() int {
	for y := 0; y < p.DimY-1; y++ {
		symmetrical := true
		smudgeX, smudgeY := -1, -1
		hasSmudge := true

		for x := 0; x < p.DimX; x++ {
			if p.Bytes[y][x] != p.Bytes[y+1][x] {
				if hasSmudge {
					smudgeX, smudgeY = x, y
					hasSmudge = false
				} else {
					if y == smudgeY+1 && x == smudgeX {
						continue
					} else {
						symmetrical = false
						break
					}
				}
			}
		}

		if symmetrical {
			for i := 0; y-i >= 0 && y+1+i < p.DimY; i++ {
				for x := 0; x < p.DimX; x++ {
					if p.Bytes[y-i][x] != p.Bytes[y+1+i][x] {
						if hasSmudge {
							smudgeX, smudgeY = x, y
							hasSmudge = false
						} else {
							if y-i == smudgeY && x == smudgeX {
								continue
							} else {
								symmetrical = false
								break
							}
						}
					}
				}
			}
		}

		if symmetrical && !hasSmudge {
			return y + 1
		}
	}
	return 0
}

func (p *Pattern) FindVerticalSymmetryWithSmudge() int {
	for x := 0; x < p.DimX-1; x++ {
		symmetrical := true
		smudgeX, smudgeY := -1, -1
		hasSmudge := true

		for y := 0; y < p.DimY; y++ {
			if p.Bytes[y][x] != p.Bytes[y][x+1] {
				if hasSmudge {
					smudgeX, smudgeY = x, y
					hasSmudge = false
				} else {
					if x == smudgeX+1 && y == smudgeY {
						continue
					} else {
						symmetrical = false
						break
					}
				}
			}
		}

		if symmetrical {
			for i := 0; x-i >= 0 && x+1+i < p.DimX; i++ {
				for y := 0; y < p.DimY; y++ {
					if p.Bytes[y][x-i] != p.Bytes[y][x+1+i] {
						if hasSmudge {
							smudgeX, smudgeY = x, y
							hasSmudge = false
						} else {
							if x-i == smudgeX && y == smudgeY {
								continue
							} else {
								symmetrical = false
								break
							}
						}
					}
				}
			}
		}

		if symmetrical && !hasSmudge {
			return x + 1
		}
	}
	return 0
}

func parsePatterns(lines []string) []Pattern {
	var patterns []Pattern
	var pattern Pattern

	for _, line := range lines {
		if line == "" {
			patterns = append(patterns, pattern)
			pattern = Pattern{}
			continue
		}

		pattern.Bytes = append(pattern.Bytes, []byte(line))
		pattern.DimX = len(line)
		pattern.DimY++
	}

	if pattern.DimY > 0 {
		patterns = append(patterns, pattern)
	}

	return patterns
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
