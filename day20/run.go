package day20

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/maps"
)

const INPUT_FILE_PATH = "day20/input.txt"
const ITERATIONS = 1000

const (
	HIGH = true
	LOW  = false
)

type Module interface {
	Propagate(Pulse) []Pulse
	GetName() string
	GetOutputs() []string
}

type FlipFlop struct {
	Name    string
	Outputs []string
	State   bool
}

type Conjunction struct {
	Name    string
	Outputs []string
	State   map[string]bool
}

type Broadcaster struct {
	Name    string
	Outputs []string
}

type Sink struct {
	Name    string
	Signals map[bool]int
}

type Output struct {
	Name string
}

type Machinery struct {
	Modules map[string]Module
}

type Pulse struct {
	From  string
	To    string
	Value bool
}

type State struct {
	FlipFlops    map[string]bool
	Conjunctions map[string]map[string]bool
	HighPulses   map[string]int
	LowPulses    map[string]int
}

func Run() {
	inputLines := readLines(INPUT_FILE_PATH)

	machinery, err := parseInput(inputLines)
	if err != nil {
		panic(err)
	}

	fmt.Println("Multiplication [PART 1]: ", machinery.PushButtonNTimesAndMultiplyPulseCounts(ITERATIONS))

	machinery, err = parseInput(inputLines)
	if err != nil {
		panic(err)
	}

	fmt.Println("Multiplication [PART 2]: ", machinery.SplitIntoSubmachineriesAndFindCommonPeriod())
}

func (machinery *Machinery) PushButtonNTimesAndMultiplyPulseCounts(n int) int {
	allHighPulses, allLowPulses := 0, 0

	for i := 0; i < n; i++ {
		state := machinery.PushButton()
		for _, highPulses := range state.HighPulses {
			allHighPulses += highPulses
		}
		for _, lowPulses := range state.LowPulses {
			allLowPulses += lowPulses
		}

	}

	return allHighPulses * allLowPulses
}

func (machinery *Machinery) SplitIntoSubmachineriesAndFindCommonPeriod() int {
	broadcastOutput := machinery.Modules["broadcaster"].GetOutputs()
	sink := machinery.FindSink()

	machineries := make(map[string]*Machinery)
	periods := make(map[string]int)

	for _, bcOutput := range broadcastOutput {
		newMachinery := machinery.ConstructSubmachinery(bcOutput, sink)
		machineries[bcOutput] = newMachinery
	}

	for bcOutput, machinery := range machineries {
		for iterations := 0; ; iterations++ {
			state := machinery.PushButton()
			if state.HighPulses["sink"] == 1 {
				periods[bcOutput] = iterations
				break
			}
		}
	}

	periodsProduct := 1
	for _, period := range periods {
		periodsProduct *= period
	}
	return periodsProduct
}

func (machinery *Machinery) FindSubgraphNodesFromSource(source, sink string) []string {
	subgraphNodes := make(map[string]bool)
	nodes := []string{source}

	for len(nodes) > 0 {
		newNodes := []string{}
		// add all nodes that are connected to the current nodes
		for _, node := range nodes {
			subgraphNodes[node] = true
		}

		// find all nodes that are connected to the current nodes
		for _, node := range nodes {
			for _, output := range machinery.Modules[node].GetOutputs() {
				// if we reached the sink, we don't need to go further
				if output == sink {
					break
				}

				// if the node is not already in the subgraph, add it to the new nodes
				if _, ok := subgraphNodes[output]; !ok {
					newNodes = append(newNodes, output)
				}
			}
		}

		nodes = newNodes
	}

	return maps.Keys(subgraphNodes)
}

func (machinery *Machinery) ConstructSubmachinery(source, sink string) *Machinery {
	// find all nodes that belong to the subgraph, i.e. they all have a path from source to sink
	subgraphNodes := machinery.FindSubgraphNodesFromSource(source, sink)
	newMachinery := Machinery{Modules: make(map[string]Module)}

	// copy all nodes that belong to the subgraph, we don't need deep copy
	// because all subgraph nodes are not connected to the rest of the machinery
	for _, moduleName := range subgraphNodes {
		newMachinery.Modules[moduleName] = machinery.Modules[moduleName]
	}

	// create a slack source node and connect it to the source
	newMachinery.Modules["broadcaster"] = &Broadcaster{
		Name:    "broadcaster",
		Outputs: []string{source},
	}

	// find the last node in the subgraph
	lastOutput := sink
	for _, module := range newMachinery.Modules {
		for _, output := range module.GetOutputs() {
			if output == sink {
				lastOutput = module.GetName()
			}
		}
	}

	// connect the last node to the sink
	for _, module := range newMachinery.Modules {
		for i, output := range module.GetOutputs() {
			if output == lastOutput {
				module.(*Conjunction).Outputs[i] = "sink"
			}
		}
	}

	// create a sink node and connect it to the last node
	delete(newMachinery.Modules, lastOutput)
	newMachinery.Modules["sink"] = &Sink{
		Name:    "sink",
		Signals: make(map[bool]int),
	}

	return &newMachinery
}

func (machinery *Machinery) FindSink() string {
	for _, module := range machinery.Modules {
		for _, output := range module.GetOutputs() {
			if _, ok := machinery.Modules[output]; !ok {
				return output
			}
		}
	}

	return ""
}

func (pulse Pulse) String() string {
	pulseValue := "LOW"
	if pulse.Value {
		pulseValue = "HIGH"
	}

	return fmt.Sprintf("%s -%s-> %s", pulse.From, pulseValue, pulse.To)
}

func (s State) Equals(other State) bool {
	for name, value := range s.FlipFlops {
		if other.FlipFlops[name] != value {
			return false
		}
	}

	for name, value := range s.Conjunctions {
		for output, state := range value {
			if other.Conjunctions[name][output] != state {
				return false
			}
		}
	}

	return true
}

func (m *Machinery) PushButton() State {
	queue := []Pulse{{"button", "broadcaster", LOW}}
	highPulses, lowPulses := make(map[string]int), make(map[string]int)

	for len(queue) > 0 {
		pulse := queue[0]
		queue = queue[1:]

		outputs := m.Propagate(pulse)
		queue = append(queue, outputs...)

		if pulse.Value == HIGH {
			highPulses[pulse.To]++
		} else {
			lowPulses[pulse.To]++
		}
	}

	state := State{
		HighPulses:   highPulses,
		LowPulses:    lowPulses,
		FlipFlops:    make(map[string]bool),
		Conjunctions: make(map[string]map[string]bool),
	}

	for _, module := range m.Modules {
		if flipFlop, ok := module.(*FlipFlop); ok {
			state.FlipFlops[flipFlop.Name] = flipFlop.State
		}
		if conjunction, ok := module.(*Conjunction); ok {
			state.Conjunctions[conjunction.Name] = conjunction.State
		}
	}

	return state
}

func (m *Machinery) Propagate(pulse Pulse) []Pulse {
	module, ok := m.Modules[pulse.To]
	if !ok {
		return []Pulse{}
	}

	return module.Propagate(pulse)
}

func (module *Broadcaster) Propagate(pulse Pulse) []Pulse {
	result := make([]Pulse, len(module.Outputs))
	for i, output := range module.Outputs {
		result[i] = Pulse{module.Name, output, pulse.Value}
	}

	return result
}

func (module *FlipFlop) Propagate(pulse Pulse) []Pulse {
	if pulse.Value == HIGH {
		return []Pulse{}
	}

	module.State = !module.State
	result := make([]Pulse, len(module.Outputs))
	for i, output := range module.Outputs {
		result[i] = Pulse{module.Name, output, module.State}
	}

	return result
}

func (module *Conjunction) Propagate(pulse Pulse) []Pulse {
	result := make([]Pulse, len(module.Outputs))
	module.State[pulse.From] = pulse.Value

	outputPulse := LOW
	for _, stateValue := range module.State {
		if stateValue == LOW {
			outputPulse = HIGH
		}
	}

	for i, output := range module.Outputs {
		result[i] = Pulse{module.Name, output, outputPulse}
	}

	return result
}

func (module *Sink) Propagate(pulse Pulse) []Pulse {
	module.Signals[pulse.Value]++
	return []Pulse{}
}

func (module Broadcaster) GetName() string {
	return module.Name
}

func (module FlipFlop) GetName() string {
	return module.Name
}

func (module Conjunction) GetName() string {
	return module.Name
}

func (module Broadcaster) GetOutputs() []string {
	return module.Outputs
}

func (module FlipFlop) GetOutputs() []string {
	return module.Outputs
}

func (module Conjunction) GetOutputs() []string {
	return module.Outputs
}

func (module Sink) GetOutputs() []string {
	return []string{}
}

func (module Sink) GetName() string {
	return module.Name
}

func parseInput(lines []string) (*Machinery, error) {
	machinery := Machinery{
		Modules: make(map[string]Module),
	}

	conjunctions := make(map[string]*Conjunction)

	for _, line := range lines {
		if line == "" {
			break
		}

		parts := strings.Split(line, " -> ")
		rawOutputs := strings.Split(parts[1], ", ")
		outputs := []string{}

		for i, rawOutput := range rawOutputs {
			output := strings.Trim(rawOutput, " ")
			if output != "" {
				outputs = append(outputs, rawOutputs[i])
			}
		}

		switch parts[0][0] {
		case '&':
			name := parts[0][1:]
			module := Conjunction{
				Name:    name,
				Outputs: outputs,
				State:   make(map[string]bool),
			}

			machinery.Modules[name] = &module
			conjunctions[name] = &module
		case '%':
			name := parts[0][1:]
			machinery.Modules[name] = &FlipFlop{
				Name:    name,
				Outputs: outputs,
				State:   LOW,
			}
		default:
			name := parts[0]
			machinery.Modules[name] = &Broadcaster{
				Name:    name,
				Outputs: outputs,
			}
		}
	}

	for _, module := range machinery.Modules {
		for _, output := range module.GetOutputs() {
			if _, ok := conjunctions[output]; ok {
				conjunctions[output].State[module.GetName()] = LOW
			}
		}
	}

	return &machinery, nil
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
