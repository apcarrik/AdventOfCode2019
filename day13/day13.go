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

func pause(seconds int) {
	duration := time.Second * time.Duration(seconds)
	time.Sleep(duration)
}

func drawGameBoard(gameBoardPtr *[][]int) {
	gameBoard := *gameBoardPtr
	for y := range gameBoard {
		for x := range gameBoard[y] {
			switch tileID := gameBoard[y][x]; {
			case tileID == 0: // empty tile
				fmt.Printf(" ")
			case tileID == 1: // wall tile
				fmt.Printf("🍞")
			case tileID == 2: // block tile
				fmt.Printf("🥑")
			case tileID == 3: // horizontal paddle tile
				fmt.Printf("▂")
			case tileID == 4: // ball tile
				fmt.Printf("⚽")
			default:
				panic(fmt.Sprintf("Error: title id %d is not recognized\n", tileID))
			}
		}
		fmt.Printf("\n")
	}

}

func updateGameBoard(drawInstructionsPtr *[3]int64, gameBoardPtr *[][]int, score int) (*[][]int, int) {
	drawInstructions := *drawInstructionsPtr
	gameBoard := *gameBoardPtr
	tileX := int(drawInstructions[0])
	tileY := int(drawInstructions[1])
	tileID := int(drawInstructions[2])
	if tileX == -1 && tileY == 0 { // Output Score
		score = int(drawInstructions[2])
	} else {
		gameBoard[tileY][tileX] = tileID
	}
	return &gameBoard, score
}

func getObjectX(gameBoardPtr *[][]int, tileID int) (int, bool) {
	// Get Object's x coordinate
	gameBoard := *gameBoardPtr
	for y := range gameBoard {
		for x := range gameBoard[y] {
			if gameBoard[y][x] == tileID {
				return x, true
			}
		}
	}
	return 0, false
}

func getPaddleDirection(gameBoardPtr *[][]int) int {
	// TODO: decide where to move the paddle, return -1 if left, 0 if no movement, 1 if right
	ballX, ballFound := getObjectX(gameBoardPtr, 4)
	paddleX, paddleFound := getObjectX(gameBoardPtr, 3)
	if ballFound && paddleFound {
		if ballX > paddleX {
			return 1
		} else if ballX < paddleX {
			return -1
		}
	}
	return 0
}

func artificalInputController(cinput chan int64, gameBoardPtr *[][]int) {
	for { // TODO: figure out how to make this better (it is slow and laggy, don't know when it will ask for input but input may be stale at that point)
		in := int64(getPaddleDirection(gameBoardPtr))
		cinput <- in
		fmt.Printf("Inputted %d\n", in)
	}
}

func createController(cfirstinput, clastoutput, cfinished chan int64, ccontrollerFinished chan bool) {
	// cfirstinput <- 1 // TODO: comment out if not inputting to program
	var drawInstructions [3]int64
	drawInstructionsIdx := 0
	// Make game board
	maxX := 45
	maxY := 24
	gameBoard := make([][]int, maxY)
	for i := 0; i < maxY; i++ {
		gameBoard[i] = make([]int, maxX)
	}
	gameBoardPtr := &gameBoard
	score := 0
	go artificalInputController(cfirstinput, gameBoardPtr)
	for {
		select {
		case lastOut := <-clastoutput:
			drawInstructions[drawInstructionsIdx] = lastOut
			drawInstructionsIdx++
			if drawInstructionsIdx > 2 {
				gameBoardPtr, score = updateGameBoard(&drawInstructions, gameBoardPtr, score)
				// gameBoard = *gameBoardPtr
				drawGameBoard(gameBoardPtr)
				fmt.Printf("Score: %d\n", score)
				drawInstructionsIdx = 0
				time.Sleep(5 * time.Millisecond)
			}
		case computerFinished := <-cfinished:
			fmt.Printf("Computer %d finished\n", computerFinished)
			fmt.Printf("Score: %d\n", score)
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
