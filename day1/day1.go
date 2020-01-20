package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// func fuel(m int) int{
// 	return m/3 - 2
// }

func fuel(m int) int {
	if(m<9) {
		return 0
	}
	f :=  m/3 - 2
	return f+fuel(f)
}

func main() {
	input, err := ioutil.ReadFile("day1/day2.input.txt")
	if err != nil {
		panic(err)
	}
	strs := strings.Split(string(input), "\n")
	fmt.Println(strs)
	f := 0
	for _, i := range strs {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		fmt.Println(fuel(j))
		f = f + fuel(j)
		// panic(err)
	}
	fmt.Println(f)
}