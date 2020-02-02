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

func absoluteValue(i int) int {
	if i<0 {
		return -i
	}
	return i
}

func parseInput(inputPtr *[]byte) *[]moon {
	input := *inputPtr
	moons := []moon{}
	re, err := regexp.Compile(`<x=(?P<x>-?\d+), y=(?P<y>-?\d+), z=(?P<z>-?\d+)>`) //<x=-3, y=15, z=-11>
	if err != nil {
		panic(err)
	}
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

func applyGravity(moonsPtr *[]moon) *[]moon {
	moons := *moonsPtr
	for i,moon1 := range moons {
		for _,moon2 := range moons {
			if moon1 != moon2 {
				// Update vx
				if moon1.x < moon2.x {
					moons[i].vx +=1
				} else if moon1.x > moon2.x {
					moons[i].vx -=1
				}
				// Update vy
				if moon1.y < moon2.y {
					moons[i].vy +=1
				} else if moon1.y > moon2.y {
					moons[i].vy -=1
				}
				// Update vz
				if moon1.z < moon2.z {
					moons[i].vz +=1
				} else if moon1.z > moon2.z {
					moons[i].vz -=1
				}
			}
		}
	}
	return &moons
}

func applyVelocity(moonsPtr *[]moon) *[]moon {
	moons := *moonsPtr
	for i,_ := range moons {
		// Update x
		moons[i].x += moons[i].vx
		// Update y
		moons[i].y += moons[i].vy
		// Update z
		moons[i].z += moons[i].vz
	}
	return &moons
}

func runTimeStep(moonsPtr *[]moon) *[]moon {
	moonsPtr = applyGravity(moonsPtr)
	moonsPtr = applyVelocity(moonsPtr)
	return 	moonsPtr
}

func calculateSystemEnergy(moonsPtr *[]moon) int {
	moons := *moonsPtr
	totalSystemEnergy := 0
	for _,moon := range moons { // [						Potential Energy						  ]	  [								Kinetic energy 								]
		totalSystemEnergy += int(absoluteValue(moon.x) + absoluteValue(moon.y) + absoluteValue(moon.z)) * (absoluteValue(moon.vx) + absoluteValue(moon.vy) + absoluteValue(moon.vz))
	}
	return totalSystemEnergy
}

func nBodyProblem(file string, numSteps int) int {

	// Get moons from input file
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	moons := *parseInput(&input)

	// Update moons for number of steps
	for i:=0; i<numSteps; i++ {		
		updatedMoonsPtr := runTimeStep(&moons)
		moons = *updatedMoonsPtr		
	}

	// Calculate total energy of system
	totalSystemEnergy := calculateSystemEnergy(&moons)

	return totalSystemEnergy
}

func main() {
	numSteps := 1000
	totalEnergy := nBodyProblem("input.txt", numSteps)
	fmt.Printf("Total energy of system after %d steps: %d", numSteps, totalEnergy)
}
