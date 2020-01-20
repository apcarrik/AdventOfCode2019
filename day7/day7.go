package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func CountDigits(i int) (count int) {
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

func GetOpcodeAndParameterMode(opcode int) (int, *[]int) {
	parameterMode := []int{0, 0, 0}
	if CountDigits(opcode) > 2 {
		var opcodeArr []int
		for opcode != 0 {
			opcodeArr = append(opcodeArr, opcode%10)
			opcode /= 10
		}
		opcode = opcodeArr[0] + opcodeArr[1]*10
		for j := 2; j < len(opcodeArr); j++ {
			parameterMode[j-2] = opcodeArr[j]
		}
	}
	return opcode, &parameterMode
}

func GetParameters(number_of_params int, instructionPtr int, parameterModePtr *[]int, intcodesPtr *[]int) *[]int {
	parameterMode := *parameterModePtr
	intcodes := *intcodesPtr
	var parameters []int
	for i := 0; i < number_of_params; i++ {
		if parameterMode[i] == 1 {
			parameters = append(parameters, intcodes[instructionPtr+i+1])
		} else {
			parameters = append(parameters, intcodes[intcodes[instructionPtr+i+1]])
		}

	}
	return &parameters
}

func runOperations(orig_intcodes *[]int, c1 chan int, c2 chan int, cfinished chan int, idx int) *[]int {
	// Clone original intcodes
	var intcodes []int
	for _, i := range *orig_intcodes {
		intcodes = append(intcodes, i)
	}

	// Do operation
	for i := 0; i < len(intcodes); {
		opcode, parameterModePtr := GetOpcodeAndParameterMode(intcodes[i])
		parameterMode := *parameterModePtr
		// fmt.Printf("i: %d, opcode: %d, parameterMode: %v, ", i, opcode, parameterMode)
		switch opcode {
		case 1: // addition
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			intcodes[intcodes[i+3]] = parameters[0] + parameters[1]
			i += 4
		case 2: // multiplication
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			intcodes[intcodes[i+3]] = parameters[0] * parameters[1]
			i += 4
		case 3: // get input
			// fmt.Printf("Enter a number: ")
			parametersPtr := GetParameters(1, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			in := <-c1
			fmt.Printf("location %d has %d\n", intcodes[i+1], intcodes[parameters[0]])
			intcodes[parameters[0]] = in
			fmt.Printf("location %d has %d\n", intcodes[i+1], intcodes[parameters[0]])
			i += 2
		case 4: // output value
			parametersPtr := GetParameters(1, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			c2 <- parameters[0]
			// fmt.Printf("%d\n", parameters[0])
			i += 2
		case 5: // jump if true
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			if parameters[0] != 0 {
				i = parameters[1]
			} else {
				i += 3
			}
		case 6: // jump if false
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			if parameters[0] == 0 {
				i = parameters[1]
			} else {
				i += 3
			}
		case 7: // less than
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			if parameters[0] < parameters[1] {
				intcodes[intcodes[i+3]] = 1
			} else {
				intcodes[intcodes[i+3]] = 0
			}
			i += 4
		case 8: // equals
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			if parameters[0] == parameters[1] {
				intcodes[intcodes[i+3]] = 1
			} else {
				intcodes[intcodes[i+3]] = 0
			}
			i += 4
		case 99:
			// fmt.Println("exiting program")
			cfinished <- idx
			return &intcodes
		default:
			panic(fmt.Sprintf("Opcode %s is not recognized", opcode))
		}
	}
	panic("Never recieved program end opcode")
	return &intcodes
}

func parseInput(input *[]byte) *[]int {
	intcodeStrings := strings.Split(string(*input), ",")
	var intcodes []int
	for _, a := range intcodeStrings {
		j, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		intcodes = append(intcodes, j)
	}
	return &intcodes
}

func createAmplifier(intcodes *[]int, cin, cfinished chan int, phaseSetting, idx int) chan int {
	var cout chan int = make(chan int)
	go runOperations(intcodes, cin, cout, cfinished, idx)
	cin <- phaseSetting
	return cout
}

func createController(numAmps int, cfirstamp, clastamp, ccontroller, cfinished chan int) {
	cfirstamp <- 0
	var lastAmpOut int
	for {
		select {
		case lastAmpOut = <-clastamp:
			go func() { cfirstamp <- lastAmpOut }()
		case ampFinished := <-cfinished:
			fmt.Printf("Amp %d finished\n", ampFinished)
			if ampFinished == numAmps-1 {
				ccontroller <- lastAmpOut
				return
			}
		}
	}
}

func testPhaseSetting(phaseSetting *[]int, intcodes *[]int) int {
	// Create a goroutine and an input & output channel for each amplifier
	numAmps := len(*phaseSetting)
	var cin chan int = make(chan int)
	var cout chan int = cin
	var cfinished chan int = make(chan int)
	for i := 0; i < numAmps; i++ {
		cout = createAmplifier(intcodes, cout, cfinished, (*phaseSetting)[i], i)
	}
	var ccontroller chan int = make(chan int)
	go createController(numAmps, cin, cout, ccontroller, cfinished)
	//fmt.Printf("phaseSetting: %v | output: %d\n", *phaseSetting, output)
	return <-ccontroller
}

func removeFromSet(item int, setPtr *[]int) *[]int {
	set := *setPtr
	var newSet []int
	for i := 0; i < len(set); i++ {
		if set[i] != item {
			newSet = append(newSet, set[i])
		}
	}
	return &newSet
}

func inputLayer(inputSet *[]int, phaseSetting *[]int, phaseSettings *[][]int) *[][]int {
	for _, i := range *inputSet {
		newInputSet := removeFromSet(i, inputSet)
		newPhaseSetting := append(*phaseSetting, i)
		if len(*inputSet) == 1 { // at top of recursion tree
			newPhaseSettings := append(*phaseSettings, newPhaseSetting)
			phaseSettings = &newPhaseSettings
		} else {
			phaseSettings = inputLayer(newInputSet, &newPhaseSetting, phaseSettings)
		}
	}
	return phaseSettings
}

func createInputCombinations(input *[]int) *[][]int {
	var phaseSetting []int
	var phaseSettings [][]int
	newPhaseSettings := inputLayer(input, &phaseSetting, &phaseSettings)
	return newPhaseSettings
}

func main() {
	input, err := ioutil.ReadFile("day7/input.txt")
	if err != nil {
		panic(err)
	}
	intcodes := parseInput(&input)

	phaseSettingsSet := []int{5, 6, 7, 8, 9}
	phaseSettings := createInputCombinations(&phaseSettingsSet)
	// fmt.Printf("%v", phaseSettings)

	// Test all permutations of input
	var maximumOutput int
	var maximumPhaseSetting []int
	for i, phaseSetting := range *phaseSettings {
		output := testPhaseSetting(&phaseSetting, intcodes)
		if i == 0 || output > maximumOutput {
			maximumOutput = output
			maximumPhaseSetting = phaseSetting
		}
	}

	fmt.Printf("maxOutput: %d\nmaxPhaseSetting: %v", maximumOutput, maximumPhaseSetting)

}
