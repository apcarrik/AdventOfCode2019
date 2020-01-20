package main

import (
	"fmt"
	"io/ioutil"
	"math"
)

type Asteroid struct {
	x                        int
	y                        int
	numAsteroidsDetectable   int
	asteroidsDetectable      []*Asteroid
	asteroidsDetectableTheta []float64
}

func findAsteroids(inputPtr *[]byte) (*[]Asteroid, int, int) {
	var asteroids []Asteroid
	input := *inputPtr
	y := 0
	x := 0
	linewidth := 0
	for _, c := range input {
		if c == 10 { // newline char
			linewidth = x + 1
			y++
			x = 0
		} else {
			if c == 35 { // '#' char
				asteroids = append(asteroids, Asteroid{x: x, y: y})
			}
			x++
		}
	}
	return &asteroids, linewidth, y + 1
}

func getSlope(a1 *Asteroid, a2 *Asteroid) float32 {
	return float32((*a2).y-(*a1).y) / float32((*a2).x-(*a1).x)
}

func getDistanceToA(a *Asteroid, otherA *Asteroid) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(a.x-otherA.x)), 2) + math.Pow(math.Abs(float64(a.y-otherA.y)), 2))
}

func getThetaToA(a *Asteroid, otherA *Asteroid) float64 {
	return math.Atan2(float64(otherA.y-a.y), float64(otherA.x-a.x))
}

func getAsteroidsDetectable(asteroidsPtr *[]Asteroid) *[]Asteroid {
	asteroids := *asteroidsPtr
	for i := range asteroids {
		a := &asteroids[i]
		for j := range asteroids {
			otherA := &asteroids[j]
			if a != otherA {
				otherATheta := getThetaToA(a, otherA)
				thetaMatch := false
				for k, t := range a.asteroidsDetectableTheta {
					if t == otherATheta {
						if getDistanceToA(a, otherA) < getDistanceToA(a, a.asteroidsDetectable[k]) { // otherA is at same angle as another asteroid and is closer to a
							a.asteroidsDetectable[k] = otherA
						}
						thetaMatch = true
					}
				}
				if !thetaMatch {
					a.asteroidsDetectable = append(a.asteroidsDetectable, otherA)
					a.asteroidsDetectableTheta = append(a.asteroidsDetectableTheta, otherATheta)
					a.numAsteroidsDetectable++
				}

			}
		}
	}
	return &asteroids
}

func getMaxAsteroidsDetectableCount(asteroidsPtr *[]Asteroid) *Asteroid {
	asteroids := *asteroidsPtr
	max := 0
	maxAsteroid := asteroids[0]
	for _, a := range asteroids {
		if a.numAsteroidsDetectable > max {
			max = a.numAsteroidsDetectable
			maxAsteroid = a
			// fmt.Printf("New max detected - Max %d from asteroid %v\n", a.numAsteroidsDetectable, a)
		}
	}
	return &maxAsteroid
}

func printAsteroidMap(asteroidsPtr *[]Asteroid, linewidth int, mapheight int) {
	asteroids := *asteroidsPtr
	// make linewidth x mapheight matrix of strings
	printMap := [][]string{{}}
	aCount := 0
	for i := 0; i < mapheight; i++ {
		for j := 0; j < linewidth; j++ {
			if asteroids[aCount].x == j && asteroids[aCount].y == i {
				printMap[i] = append(printMap[i], fmt.Sprintf(" %03d ", asteroids[aCount].numAsteroidsDetectable))
				if aCount < len(asteroids)-1 {
					aCount++
				}
			} else {
				printMap[i] = append(printMap[i], " ... ")
			}
		}
		printMap = append(printMap, []string{})
	}

	outputStr := ""
	for i := 0; i < mapheight; i++ {
		for j := 0; j < linewidth; j++ {
			outputStr += printMap[i][j]
		}
		outputStr += "\n"
	}

	fmt.Println(outputStr)
}

func printSpecificAsteroidMap(asteroid *Asteroid, linewidth int, mapheight int) {
	// make linewidth x mapheight matrix of strings
	printMap := [][]string{{}}
	acount := 0
	for i := 0; i < mapheight; i++ {
		for j := 0; j < linewidth; j++ {
			asteroidFound := false
			if asteroid.x == j && asteroid.y == i {
				asteroidFound = true
				printMap[i] = append(printMap[i], " ### ")
			} else {
				for _, a := range asteroid.asteroidsDetectable {
					if a.x == j && a.y == i {
						asteroidFound = true
						printMap[i] = append(printMap[i], fmt.Sprintf(" %03d ", acount))
						acount++
					}
				}
			}
			if !asteroidFound {
				printMap[i] = append(printMap[i], " ... ")
			}
		}
		printMap = append(printMap, []string{})
	}

	outputStr := ""
	for i := 0; i < mapheight; i++ {
		for j := 0; j < linewidth; j++ {
			outputStr += printMap[i][j]
		}
		outputStr += "\n"
	}

	fmt.Println(outputStr)
}

func printVaporizedAsteroidMap(asteroid *Asteroid, vaporizedAsteroidsPtr *[]*Asteroid, linewidth int, mapheight int) {
	// make linewidth x mapheight matrix of strings
	vaporizedAsteroids := *vaporizedAsteroidsPtr
	printMap := [][]string{{}}
	for i := 0; i < mapheight; i++ {
		for j := 0; j < linewidth; j++ {
			asteroidFound := false
			if asteroid.x == j && asteroid.y == i {
				asteroidFound = true
				printMap[i] = append(printMap[i], " ### ")
			} else {
				for k, a := range vaporizedAsteroids {
					if a.x == j && a.y == i {
						asteroidFound = true
						printMap[i] = append(printMap[i], fmt.Sprintf(" %03d ", k))
					}
				}
			}
			if !asteroidFound {
				printMap[i] = append(printMap[i], " ... ")
			}
		}
		printMap = append(printMap, []string{})
	}

	outputStr := ""
	for i := 0; i < mapheight; i++ {
		for j := 0; j < linewidth; j++ {
			outputStr += printMap[i][j]
		}
		outputStr += "\n"
	}

	fmt.Println(outputStr)
}

func removeAsteroids(a *Asteroid, removedAsteroidsPtr *[]*Asteroid, lowerBound float64, upperBound float64) []*Asteroid {
	removedAsteroids := *removedAsteroidsPtr
	for true{
		var localOtherA *Asteroid
		var localOtherATheta float64
		var localOtherAIndex int
		for i,otherA := range a.asteroidsDetectable {
			otherATheta := getThetaToA(a,otherA)
			if otherATheta <= upperBound && otherATheta >= lowerBound {
				if localOtherA == nil || otherATheta < localOtherATheta {
					localOtherA = otherA
					localOtherATheta = otherATheta
					localOtherAIndex = i
				}
			}		
		}
		if localOtherA == nil {
			break
		}
		removedAsteroids = append(removedAsteroids, localOtherA)
		// remove element from a.asteroidsDetectable
		a.asteroidsDetectable[localOtherAIndex] = a.asteroidsDetectable[len(a.asteroidsDetectable)-1]
		a.asteroidsDetectable[len(a.asteroidsDetectable)-1] = nil
		a.asteroidsDetectable = a.asteroidsDetectable[:len(a.asteroidsDetectable)-1]

	}
	return removedAsteroids
}

func getVaporizedAsteroidsAroundA(asteroidsPtr *[]Asteroid, a *Asteroid) []*Asteroid {
	asteroids := *asteroidsPtr
	// fmt.Printf("%v\n",asteroids[0])
	asteroids = asteroids
	// var newAsteroids []*Asteroid
	// for _,a := range origAsteroids{
	// 	newAsteroids = append(newAsteroids, &a)
	// }
	removedAsteroids := []*Asteroid{}

	// Do 1 rotation starting from striaght up
	// Find quadrant 1 detectable asteroids
	removedAsteroids = removeAsteroids(a, &removedAsteroids, -math.Pi/2, -0)
	// Find quadrant 2 detectable asteroids
	removedAsteroids = removeAsteroids(a, &removedAsteroids, +0, math.Pi/2)
	// Find quadrant 3 detectable asteroids
	removedAsteroids = removeAsteroids(a, &removedAsteroids, math.Pi/2, math.Pi)
	// Find quadrant 4 detectable asteroids
	removedAsteroids = removeAsteroids(a, &removedAsteroids, -math.Pi, -math.Pi/2)


	// for true{
	// 	var localOtherA *Asteroid
	// 	var localOtherATheta float64
	// 	var localOtherAIndex int
	// 	for i,otherA := range a.asteroidsDetectable {
	// 		otherATheta := getThetaToA(a,otherA)
	// 		if otherATheta <= math.Pi/2 && otherATheta <= +0 {
	// 			if localOtherA == nil || otherATheta > localOtherATheta {
	// 				localOtherA = otherA
	// 				localOtherATheta = otherATheta
	// 				localOtherAIndex = i
	// 			}
	// 		}		
	// 	}
	// 	if localOtherA == nil {
	// 		break
	// 	}
	// 	removedAsteroids = append(removedAsteroids, localOtherA)
	// 	// remove element from a.asteroidsDetectable
	// 	a.asteroidsDetectable[localOtherAIndex] = a.asteroidsDetectable[len(a.asteroidsDetectable)-1]
	// 	a.asteroidsDetectable[len(a.asteroidsDetectable)-1] = nil
	// 	a.asteroidsDetectable = a.asteroidsDetectable[:len(a.asteroidsDetectable)-1]

	// }
	return removedAsteroids


	// Get new detectable asteroids
	// Repeat
}

func main() {
	input, err := ioutil.ReadFile("day10/input.txt")
	if err != nil {
		panic(err)
	}
	// Create Asteriods
	asteroidsPtr, mapwidth, mapheight := findAsteroids(&input)
	// asteroidsPtr, _, _ := findAsteroids(&input)
	asteroids := *asteroidsPtr

	// Find asteroids detectable for each Asteroid
	asteroids = *getAsteroidsDetectable(&asteroids)

	// Print max detectable count for Asteroids
	maxAsteroid := *getMaxAsteroidsDetectableCount(&asteroids)
	fmt.Printf("Max Asteroids Detected: %d\nAsteroid coordinates: (%d,%d)\n\n", maxAsteroid.numAsteroidsDetectable, maxAsteroid.x, maxAsteroid.y)

	// Vaporize asteroids surrounding maxAsteroid
	asteroidsVaporized := getVaporizedAsteroidsAroundA(&asteroids, &maxAsteroid)
	fmt.Printf("Removed %d Asteroids\n", len(asteroidsVaporized))

	twoHundoVapoX := asteroidsVaporized[199].x
	twoHundoVapoY := asteroidsVaporized[199].y
	fmt.Printf("200th Vaporized Asteroid Coordinates: (%d,%d)\nAnswer: %d\n", twoHundoVapoX, twoHundoVapoY, 100*twoHundoVapoX+twoHundoVapoY)
	printVaporizedAsteroidMap(&maxAsteroid, &asteroidsVaporized, mapwidth, mapheight)


	// Print asteroid map with detectable counts
	// printAsteroidMap(&asteroids, mapwidth, mapheight)
	// printSpecificAsteroidMap(&asteroids[4], mapwidth, mapheight)
	// fmt.Printf("%v\nasteroids: %v", input, asteroids)

}
