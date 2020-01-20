package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)


// Part 1
// func main() {
// 	input, err := ioutil.ReadFile("day2/input.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	intcodeStrings := strings.Split(string(input), ",")
// 	var intcodes []int
// 	for _,a := range intcodeStrings{
// 		j, err := strconv.Atoi(a)
// 		if err != nil {
// 			panic(err)
// 		}
// 		intcodes = append(intcodes,j)
// 	}
// 	fmt.Println("Original:")
// 	fmt.Println(intcodes)
// 	for i:=0; i < len(intcodes); i+=4 {
// 		opcode := intcodes[i]
// 		switch opcode {
// 		case 1:
// 			intcodes[intcodes[i+3]] = intcodes[intcodes[i+1]] + intcodes[intcodes[i+2]] 
// 		case 2:
// 			intcodes[intcodes[i+3]] = intcodes[intcodes[i+1]] * intcodes[intcodes[i+2]] 
// 		case 99:
// 			fmt.Println("\n\nNew:")
// 			fmt.Println(intcodes)
// 			return;
// 		default:
// 			panic(fmt.Sprintf("Opcode %s is not recognized", opcode))
// 		}
// 	}
// }


// Part 2
func main() {
	input, err := ioutil.ReadFile("day2/input.txt")
	if err != nil {
		panic(err)
	}
	intcodeStrings := strings.Split(string(input), ",")
	var intcodesOriginal []int
	for _,a := range intcodeStrings{
		j, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		intcodesOriginal = append(intcodesOriginal,j)
	}
	min := -1

	answerfound := false
	for i:=0; i<=99; i++{
		for j:=0; j<=99; j++{
			intcodesNew := make([]int, len(intcodesOriginal))
			copy(intcodesNew, intcodesOriginal)
			intcodesNew[1] = i
			intcodesNew[2] = j		
			newval := runComputer(&intcodesNew)
			localmin := Abs(newval - 19690720)
			// fmt.Printf("Trying pos 1 = %d and pos 2 = %d , got %d \n", i, j, newval)
			if min == -1 {
				min = localmin
			} else if localmin < min {
				min = localmin	
			}
			if localmin == 0{
				fmt.Printf("Solution found with pos 1 = %d and pos 2 = %d\n", i, j)
				fmt.Printf("Answer is %d", 100*i+j)
				return;
			}
			if (answerfound) { break}
		}
		if (answerfound) { break}
	}
	fmt.Printf("No solution found, min was %d\n", min)

}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runComputer(intcodesptr *[]int) int {
	intcodes := *intcodesptr
	for i:=0; i < len(intcodes); i+=4 {
		opcode := intcodes[i]
		switch opcode {
		case 1:
			intcodes[intcodes[i+3]] = intcodes[intcodes[i+1]] + intcodes[intcodes[i+2]] 
		case 2:
			intcodes[intcodes[i+3]] = intcodes[intcodes[i+1]] * intcodes[intcodes[i+2]] 
		case 99:
			return intcodes[0];
		default:
			panic(fmt.Sprintf("Opcode %s is not recognized", opcode))
		}
	}
	return -1;
}