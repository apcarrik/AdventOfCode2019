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

// func runTimeStep(moonsPtr *[]moon) *[]moon {
// 	moons := *moonsPtr
// 	// Apply gravity
// 	for _,moon1 := range moons {
// 		for _,moon2 := range moons {
// 			if moon1 != moon2 {
// 				// Update vx
// 				if moon1.x < moon2.x {
// 					moon1.vx +=1
// 				} else if moon1.x < moon2.x {
// 					moon1.vx -=1
// 				}
// 				// Update vy
// 				if moon1.y < moon2.y {
// 					moon1.vy +=1
// 				} else if moon1.y < moon2.y {
// 					moon1.vy -=1
// 				}
// 				// Update vz
// 				if moon1.z < moon2.z {
// 					moon1.vz +=1
// 				} else if moon1.z < moon2.z {
// 					moon1.vz -=1
// 				}
// 			} else {
// 				fmt.Printf("duplicate: %v == %v\n", &moon1, &moon2)
// 			}
// 		}
// 	}
// 	// TODO: Apply velocity
// 	return &moons
// }

func nBodyProblem(file string, numSteps int) int {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	// Get moons from input file
	moons := *parseInput(&input)
	fmt.Printf("moons: %v\n", moons)

	// TODO: Update moons for number of steps
	// for i:=0; i<numSteps; i++ {		
	// 	updatedMoonsPtr := runTimeStep(&moons)
	// 	moons := *updatedMoonsPtr
	// }

	// TODO: Calculate total energy of system
	return 0
}

func main() {
	numSteps := 1000
	totalEnergy := nBodyProblem("input.txt", numSteps)
	fmt.Printf("Total energy of system after %d steps: %d", numSteps, totalEnergy)
}
