package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	// "github.com/jinzhu/copier"
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

func runOperations(orig_intcodes *[]int) *[]int {
	// Clone original intcodes
	var intcodes []int
	for _, i := range *orig_intcodes {
		intcodes = append(intcodes, i)
	}
	reader := bufio.NewReader(os.Stdin)

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
			fmt.Printf("Enter a number: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			in, err := strconv.Atoi(strings.Replace(text, "\n", "", -1))
			if err != nil {
				panic(err)
			}
			intcodes[intcodes[i+1]] = in
			i += 2
		case 4: // output value
			parametersPtr := GetParameters(1, i, &parameterMode, &intcodes)
			parameters := *parametersPtr
			fmt.Printf("%d\n", parameters[0])
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
			fmt.Println("exiting program")
			return &intcodes
		default:
			panic(fmt.Sprintf("Opcode %s is not recognized", opcode))
		}
	}
	panic("Never recieved program end opcode")
	return &intcodes

}

func main() {
	input, err := ioutil.ReadFile("day5/input.txt")
	if err != nil {
		panic(err)
	}
	intcodeStrings := strings.Split(string(input), ",")
	var intcodes []int
	for _, a := range intcodeStrings {
		j, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		intcodes = append(intcodes, j)
	}
	fmt.Println("Original:")
	fmt.Println(intcodes)

	newintcodes := *runOperations(&intcodes)

	fmt.Println("\n\nNew:")
	fmt.Println(newintcodes)

}
