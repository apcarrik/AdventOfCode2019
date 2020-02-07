package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
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

func createComputer(intcodes *[]int64, cin, cfinished chan int64, idx int, relativeBase int64) chan int64 {
	var cout chan int64 = make(chan int64)
	go runOperations(intcodes, cin, cout, cfinished, idx, relativeBase)
	return cout
}

// func panelNotChangedBefore(panelsChangedPtr *[][]int, thisPanelPtr *[]int) bool{
// 	panelsChanged := *panelsChangedPtr
// 	thisPanel := *thisPanelPtr
// 	for _,panel := range panelsChanged {
// 		if panel[0] == thisPanel[0] && panel[1] == thisPanel[1] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func setPanelColor(panelsChangedPtr *[][]int, panelGridPtr *[][]int, robotLocationPtr *[]int, color int) (*[][]int, *[][]int) {
// 	panelsChanged := *panelsChangedPtr
// 	robotLocation := *robotLocationPtr
// 	panelGrid := *panelGridPtr
// 	panelGrid[robotLocation[1]][robotLocation[0]] = color
// 	if panelNotChangedBefore(panelsChangedPtr, robotLocationPtr){
// 		panelsChanged = append(panelsChanged, []int{robotLocation[0], robotLocation[1]})
// 	}
// 	return &panelGrid, &panelsChanged
// }

// func moveRobot(robotLocationPtr *[]int, currentDirection int, newDirection int) (int, *[]int) {
// 	robotLocation := *robotLocationPtr
// 	robotDirection := currentDirection
// 	switch newDirection {
// 	case 0: // left
// 		robotDirection--
// 	case 1: // right
// 		robotDirection++
// 	default:
// 		panic(fmt.Sprintf("Direction %d is not recognized", newDirection))
// 	}
// 	if robotDirection < 0 {
// 		robotDirection = 3
// 	}
// 	if robotDirection > 3 {
// 		robotDirection = 0
// 	}

// 	switch robotDirection {
// 	case 0: // 0 = up
// 		robotLocation[1]--
// 	case 1: // 1 = right
// 		robotLocation[0]++
// 	case 2: // 2 = down
// 		robotLocation[1]++
// 	case 3: // 3 = left
// 		robotLocation[0]--
// 	default:
// 		panic(fmt.Sprintf("robotDirection %d is not recognized", newDirection))
// 	}

// 	return robotDirection, &robotLocation

// }

// func checkColor(panelGridPtr *[][]int, robotLocationPtr *[]int) int {
// 	robotLocation := *robotLocationPtr
// 	panelGrid := *panelGridPtr
// 	retVal := panelGrid[robotLocation[1]][robotLocation[0]]
// 	return retVal
// }

// func printPanels(panelGridPtr *[][]int, robotLocationPtr *[]int, robotDirection int) { // TODO: put this somewhere and add pause between painting
// 	robotLocation := *robotLocationPtr
// 	panelGrid := *panelGridPtr
// 	robotDirectionChar := 0x52 // "R"
// 	switch robotDirection {
// 		case 0: // "Up"
// 			robotDirectionChar = 0x5E // "^"
// 		case 1: // "Right"
// 			robotDirectionChar = 0x3E // ">"
// 		case 2: // "Down"
// 			robotDirectionChar = 0x76 // "v"
// 		case 3: // "Left"
// 			robotDirectionChar = 0x3C // "<"
// 	}
// 	for y, panelRow := range panelGrid {
// 		for x, panel := range panelRow {
// 			if x == robotLocation[0] && y == robotLocation[1] {
// 				fmt.Printf("%c", robotDirectionChar)
// 			} else {
// 				if panel == 0 {
// 					fmt.Printf(".")
// 				}else{
// 					fmt.Printf("#")
// 				}
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

func pause(seconds int) {
	duration := time.Second * time.Duration(seconds)
	time.Sleep(duration)
}

func checkIfBlockTile(drawInstructionsPtr *[3]int64) bool {
	drawInstructions := *drawInstructionsPtr
	switch titleID := drawInstructions[2]; {
	case titleID == 0: // empty tile
	case titleID == 1: // wall tile
	case titleID == 2: // block tile
		return true
	case titleID == 3: // horizontal paddle tile
	case titleID == 4: // ball tile
	default:
		panic(fmt.Sprintf("Error: title id %d is not recognized\n", titleID))
	}
	return false
}

func createController(cfirstinput, clastoutput, cfinished chan int64, ccontrollerFinished chan bool) {
	// cfirstinput <- 1 // TODO: comment out if not inputting to program
	var drawInstructions [3]int64
	drawInstructionsIdx := 0
	blockTileCounter := 0
	for {
		select {
		case lastOut := <-clastoutput:
			drawInstructions[drawInstructionsIdx] = lastOut
			drawInstructionsIdx++
			if drawInstructionsIdx > 2 {
				if checkIfBlockTile(&drawInstructions) {
					blockTileCounter++
				}
				drawInstructionsIdx = 0
			}
		case computerFinished := <-cfinished:
			fmt.Printf("Computer %d finished\n", computerFinished)
			fmt.Printf("Number of Block Tiles: %d\n", blockTileCounter)
			ccontrollerFinished <- true
			return
		}
	}
}

func runProgram(numComputers int, intcodes *[]int64, relativeBase int64) bool {
	// Create a goroutine and an input & output chanel for each amplifier
	cin := make(chan int64)
	cout := cin
	cfinished := make(chan int64)
	ccontrollerFinished := make(chan bool)
	for i := 0; i < numComputers; i++ {
		cout = createComputer(intcodes, cout, cfinished, i, relativeBase)
	}
	go createController(cin, cout, cfinished, ccontrollerFinished)
	return <-ccontrollerFinished
}

func main() {
	input, err := ioutil.ReadFile("input/input.txt")
	if err != nil {
		panic(err)
	}
	intcodes := parseInput(&input)
	// fmt.Printf("len: %d\n", len(*intcodes))
	// fmt.Println("Original:")
	// fmt.Println(intcodes)

	numComputers := 1
	output := runProgram(numComputers, intcodes, 0)

	fmt.Printf("\n\nOutput:%t", output)

}
