package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
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
	if i < 0 {
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
	for i, moon1 := range moons {
		for _, moon2 := range moons {
			if moon1 != moon2 {
				// Update vx
				if moon1.x < moon2.x {
					moons[i].vx += 1
				} else if moon1.x > moon2.x {
					moons[i].vx -= 1
				}
				// Update vy
				if moon1.y < moon2.y {
					moons[i].vy += 1
				} else if moon1.y > moon2.y {
					moons[i].vy -= 1
				}
				// Update vz
				if moon1.z < moon2.z {
					moons[i].vz += 1
				} else if moon1.z > moon2.z {
					moons[i].vz -= 1
				}
			}
		}
	}
	return &moons
}

func applyVelocity(moonsPtr *[]moon) *[]moon {
	moons := *moonsPtr
	for i, _ := range moons {
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
	return moonsPtr
}

func calculateSystemEnergy(moonsPtr *[]moon) int {
	moons := *moonsPtr
	totalSystemEnergy := 0
	for _, moon := range moons { // [						Potential Energy						  ]	  [								Kinetic energy 								]
		totalSystemEnergy += int(absoluteValue(moon.x)+absoluteValue(moon.y)+absoluteValue(moon.z)) * (absoluteValue(moon.vx) + absoluteValue(moon.vy) + absoluteValue(moon.vz))
	}
	return totalSystemEnergy
}

func nBodyProblemPart1(file string, numSteps int) int {

	// Get moons from input file
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	moons := *parseInput(&input)

	// Update moons for number of steps
	for i := 0; i < numSteps; i++ {
		updatedMoonsPtr := runTimeStep(&moons)
		moons = *updatedMoonsPtr
	}

	// Calculate total energy of system
	totalSystemEnergy := calculateSystemEnergy(&moons)

	return totalSystemEnergy
}

// Part 2
type moonset struct {
	moons     *[]moon
	timeStamp int
}

func makeHashes(moonsPtr *[]moon) (int, int, int) {
	moons := *moonsPtr
	xHash := 0
	yHash := 0
	zHash := 0
	for _, moon := range moons {
		xHash += absoluteValue(moon.x) * absoluteValue(moon.vx)
		yHash += absoluteValue(moon.y) * absoluteValue(moon.vy)
		zHash += absoluteValue(moon.z) * absoluteValue(moon.vz)
	}
	return xHash, yHash, zHash
}

func deepCopyMoons(moonsPtr *[]moon) *[]moon {
	moons := *moonsPtr
	newMoons := []moon{}
	for _, m := range moons {
		newMoons = append(newMoons, moon{
			x:  m.x,
			y:  m.y,
			z:  m.z,
			vx: m.vx,
			vy: m.vy,
			vz: m.vz,
		})
	}
	return &newMoons
}

func gcd(a int, b int) int {
	var t int
	for b != 0 {
		t = b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int) int {
	return (a * b / gcd(a, b))
}

func nBodyProblemPart2(file string, maxSteps int) int {

	// Get moons from input file
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	moons := *parseInput(&input)
	dimensionPeriods := []int{0, 0, 0}
	xMap := make(map[int][]moonset)
	yMap := make(map[int][]moonset)
	zMap := make(map[int][]moonset)
	for i := 0; i < maxSteps; i++ {
		updatedMoonsPtr := runTimeStep(&moons)
		moons = *updatedMoonsPtr
		xHash, yHash, zHash := makeHashes(&moons)
		newMoonsPtr := deepCopyMoons(&moons)
		xMoonSets, xMoonSetsExist := xMap[xHash]
		if xMoonSetsExist {
			for _, moonSet := range xMoonSets {
				correctCoordinatesCounter := 0
				for j := range *moonSet.moons {
					if (*moonSet.moons)[j].x == (*newMoonsPtr)[j].x && (*moonSet.moons)[j].vx == (*newMoonsPtr)[j].vx {
						correctCoordinatesCounter++
					}
				}
				if correctCoordinatesCounter == 4 {
					dimensionPeriods[0] = i - moonSet.timeStamp
				}
			}
			// place into xmap
			xMap[xHash] = append(xMap[xHash], moonset{
				moons:     newMoonsPtr,
				timeStamp: i,
			})
		} else {
			// place into xmap
			xMap[xHash] = []moonset{
				moonset{
					moons:     newMoonsPtr,
					timeStamp: i,
				},
			}
		}

		yMoonSets, yMoonSetsExist := yMap[yHash]
		if yMoonSetsExist {
			for _, moonSet := range yMoonSets {
				correctCoordinatesCounter := 0
				for j := range *moonSet.moons {
					if (*moonSet.moons)[j].y == (*newMoonsPtr)[j].y && (*moonSet.moons)[j].vy == (*newMoonsPtr)[j].vy {
						correctCoordinatesCounter++
					}
				}
				if correctCoordinatesCounter == 4 {
					dimensionPeriods[1] = i - moonSet.timeStamp
				}
			}
			// place into ymap
			yMap[yHash] = append(yMap[yHash], moonset{
				moons:     newMoonsPtr,
				timeStamp: i,
			})
		} else {
			// place into ymap
			yMap[yHash] = []moonset{
				moonset{
					moons:     newMoonsPtr,
					timeStamp: i,
				},
			}
		}

		zMoonSets, zMoonSetsExist := zMap[zHash]
		if zMoonSetsExist {
			for _, moonSet := range zMoonSets {
				correctCoordinatesCounter := 0
				for j := range *moonSet.moons {
					if (*moonSet.moons)[j].z == (*newMoonsPtr)[j].z && (*moonSet.moons)[j].vz == (*newMoonsPtr)[j].vz {
						correctCoordinatesCounter++
					}
				}
				if correctCoordinatesCounter == 4 {
					dimensionPeriods[2] = i - moonSet.timeStamp
				}
			}
			// place into zmap
			zMap[zHash] = append(zMap[zHash], moonset{
				moons:     newMoonsPtr,
				timeStamp: i,
			})
		} else {
			// place into zmap
			zMap[zHash] = []moonset{
				moonset{
					moons:     newMoonsPtr,
					timeStamp: i,
				},
			}
		}

		if dimensionPeriods[0] != 0 && dimensionPeriods[1] != 0 && dimensionPeriods[2] != 0 {
			return lcm(dimensionPeriods[0], lcm(dimensionPeriods[1], dimensionPeriods[2]))
		}
	}

	return -1
}

func main() {
	start := time.Now()
	inputFile := "input.txt"
	iterations := 4686774925
	xResult := nBodyProblemPart2(inputFile, iterations)
	fmt.Printf("xResult: %d/n", xResult)
	elapsed := time.Since(start)
	fmt.Printf("Program took %s", elapsed)
}
