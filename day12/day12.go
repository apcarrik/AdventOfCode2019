package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

type moon struct {
	x  int
	y  int
	z  int
	vx int
	vy int
	vz int
}

func parseInput(inputPtr *[]byte) *[]moon {
	input := *inputPtr
	moons := []moon{}
	re := regexp.MustCompile(`<x=(?P<x>-?\d+), y=(?P<y>-?\d+), z=(?P<z>-?\d+)>`) //<x=-3, y=15, z=-11>
	for _, j := range re.FindAllSubmatch(input, -1) {
		xint, err := strconv.Atoi(string(j[1]))
		if err != nil {
			panic(err)
		}
		yint, err := strconv.Atoi(string(j[2]))
		if err != nil {
			panic(err)
		}
		zint, err := strconv.Atoi(string(j[3]))
		if err != nil {
			panic(err)
		}
		newMoon := moon{x: xint, y: yint, z: zint}
		moons = append(moons, newMoon)
	}
	return &moons
}

func main() {
	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	moons := *parseInput(&input)
	fmt.Printf("%v", moons)
}
