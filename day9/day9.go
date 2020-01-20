package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func CountDigits(i int64) (count int64) {
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

func GetOpcodeAndParameterMode(opcode int64) (int64, *[]int64) {
	parameterMode := []int64{0, 0, 0}
	if CountDigits(opcode) > 2 {
		var opcodeArr []int64
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

func GetParameters(number_of_params int, instructionPtr int, parameterModePtr *[]int64, intcodesPtr *[]int64, relativeBase int64) *[]int64 {
	parameterMode := *parameterModePtr
	intcodes := *intcodesPtr
	var parameters []int64
	for i := 0; i < number_of_params; i++ {
		switch parameterMode[i] {
		case 1: // Immediate Mode
			parameters = append(parameters, int64(instructionPtr+i+1))
		case 2: // Relative Mode
			parameters = append(parameters, intcodes[instructionPtr+(i+1)]+relativeBase)
			// fmt.Printf("%d\n", intcodes[instructionPtr+(i+1)]+relativeBase)
		default: // Position Mode
			parameters = append(parameters, intcodes[instructionPtr+(i+1)])
		}

	}
	return &parameters
}

func runOperations(orig_intcodes *[]int64, c1 chan int64, c2 chan int64, cfinished chan int64, idx int, relativeBase int64) *[]int64 {
	// Clone original intcodes
	var intcodes []int64
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
			parametersPtr := GetParameters(3, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			intcodes[parameters[2]] = intcodes[parameters[0]] + intcodes[parameters[1]]
			i += 4
		case 2: // multiplication
			parametersPtr := GetParameters(3, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr

			intcodes[parameters[2]] = intcodes[parameters[0]] * intcodes[parameters[1]]
			i += 4
		case 3: // get input
			parametersPtr := GetParameters(1, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			in := <-c1
			intcodes[parameters[0]] = in
			i += 2
		case 4: // output value
			parametersPtr := GetParameters(1, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			c2 <- intcodes[parameters[0]]
			i += 2
		case 5: // jump if true
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			if intcodes[parameters[0]] != 0 {
				i = int(intcodes[parameters[1]])
			} else {
				i += 3
			}
		case 6: // jump if false
			parametersPtr := GetParameters(2, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			if intcodes[parameters[0]] == 0 {
				i = int(intcodes[parameters[1]])
			} else {
				i += 3
			}
		case 7: // less than
			parametersPtr := GetParameters(3, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			if intcodes[parameters[0]] < intcodes[parameters[1]] {
				intcodes[parameters[2]] = 1
			} else {
				intcodes[parameters[2]] = 0
			}
			i += 4
		case 8: // equals
			parametersPtr := GetParameters(3, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			if intcodes[parameters[0]] == intcodes[parameters[1]] {
				intcodes[parameters[2]] = 1
			} else {
				intcodes[parameters[2]] = 0
			}
			i += 4
		case 9: // adjust relative base
			parametersPtr := GetParameters(1, i, &parameterMode, &intcodes, relativeBase)
			parameters := *parametersPtr
			relativeBase += intcodes[parameters[0]]
			i += 2
		case 99:
			// fmt.Printf("modified: %v\n", intcodes)
			// fmt.Println("exiting program")
			cfinished <- int64(idx)
			return &intcodes
		default:
			panic(fmt.Sprintf("Opcode %s is not recognized", opcode))
		}
	}
	panic("Never recieved program end opcode")
	return &intcodes
}

func parseInput(input *[]byte) *[]int64 {
	intcodeStrings := strings.Split(string(*input), ",")
	intcodes := make([]int64, len(intcodeStrings)*10)
	for i, a := range intcodeStrings {
		j, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			panic(err)
		}
		intcodes[i] = j
	}
	return &intcodes
}

func createAmplifier(intcodes *[]int64, cin, cfinished chan int64, phaseSetting int64, idx int, relativeBase int64) chan int64 {
	var cout chan int64 = make(chan int64)
	go runOperations(intcodes, cin, cout, cfinished, idx, relativeBase)
	// cin <- phaseSetting // TODO: comment out if not inputting to program
	return cout
}

func createController(numAmps int, cfirstamp, clastamp, ccontroller, cfinished chan int64) {
	cfirstamp <- 2 // TODO: comment out if not inputting to program
	var lastAmpOut int64
	for {
		select {
		case lastAmpOut = <-clastamp:
			// go func() { cfirstamp <- lastAmpOut }()
			fmt.Printf("lastAmpOut: %d\n", lastAmpOut)
		case ampFinished := <-cfinished:
			fmt.Printf("Amp %d finished\n", ampFinished)
			if ampFinished == int64(numAmps-1) {
				ccontroller <- lastAmpOut
				return
			}
		}
	}
}

func testPhaseSetting(phaseSetting *[]int64, intcodes *[]int64, relativeBase int64) int64 {
	// Create a goroutine and an input & output channel for each amplifier
	numAmps := len(*phaseSetting)
	var cin chan int64 = make(chan int64)
	var cout chan int64 = cin
	var cfinished chan int64 = make(chan int64)
	for i := 0; i < numAmps; i++ {
		cout = createAmplifier(intcodes, cout, cfinished, (*phaseSetting)[i], i, relativeBase)
	}
	var ccontroller chan int64 = make(chan int64)
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
	input, err := ioutil.ReadFile("day9/input.txt")
	if err != nil {
		panic(err)
	}
	intcodes := parseInput(&input)
	// fmt.Printf("len: %d\n", len(*intcodes))
	// fmt.Println("Original:")
	// fmt.Println(intcodes)

	phaseSetting := []int64{1}
	output := testPhaseSetting(&phaseSetting, intcodes, 0)

	fmt.Printf("\n\nOutput:%d", output)

}
