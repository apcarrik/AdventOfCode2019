package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

func parseTestInput(inputPtr *[]byte) *[]moon {
	input := *inputPtr
	moons := []moon{}
	re, err := regexp.Compile(`pos=<x=\s*(?P<x>-?\d+), y=\s*(?P<y>-?\d+), z=\s*(?P<z>-?\d+)>, vel=<x=\s*(?P<vx>-?\d+), y=\s*(?P<vy>-?\d+), z=\s*(?P<vz>-?\d+)>`)
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
		vxint, err := strconv.Atoi(string(j[4]))
		if err != nil {
			panic(err)
		}
		vyint, err := strconv.Atoi(string(j[5]))
		if err != nil {
			panic(err)
		}
		vzint, err := strconv.Atoi(string(j[6]))
		if err != nil {
			panic(err)
		}
		newMoon := moon{x: xint, y: yint, z: zint, vx: vxint, vy: vyint, vz: vzint}
		moons = append(moons, newMoon)
	}
	return &moons
}

func Test_UnitTest_parseInput(t *testing.T) {
	expectedResult := []moon{
		moon{x: -3, y: 15, z: -11, vx: 0, vy: 0, vz: 0},
		moon{x: 3, y: 13, z: -19, vx: 0, vy: 0, vz: 0},
		moon{x: -13, y: 18, z: -2, vx: 0, vy: 0, vz: 0},
		moon{x: 6, y: 0, z: -1, vx: 0, vy: 0, vz: 0},
	}
	inputByteArray := []byte("<x=-3, y=15, z=-11>\n<x=3, y=13, z=-19>\n<x=-13, y=18, z=-2>\n<x=6, y=0, z=-1>")
	actualResult := parseInput(&inputByteArray)
	if !reflect.DeepEqual(*actualResult, expectedResult) {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_applyGravity(t *testing.T) {
	expectedResult := []moon{
		moon{x: -1, y: 0, z: 2, vx: 3, vy: -1, vz: -1},
		moon{x: 2, y: -10, z: -7, vx: 1, vy: 3, vz: 3},
		moon{x: 4, y: -8, z: 8, vx: -3, vy: 1, vz: -3},
		moon{x: 3, y: 5, z: -1, vx: -1, vy: -3, vz: 1},
	}
	inputMoons := []moon{
		moon{x: -1, y: 0, z: 2, vx: 0, vy: 0, vz: 0},
		moon{x: 2, y: -10, z: -7, vx: 0, vy: 0, vz: 0},
		moon{x: 4, y: -8, z: 8, vx: 0, vy: 0, vz: 0},
		moon{x: 3, y: 5, z: -1, vx: 0, vy: 0, vz: 0},
	}
	actualResult := applyGravity(&inputMoons)
	if !reflect.DeepEqual(*actualResult, expectedResult) {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_applyVelocity(t *testing.T) {
	expectedResult := []moon{
		moon{x: 2, y: -1, z: 1, vx: 3, vy: -1, vz: -1},
		moon{x: 3, y: -7, z: -4, vx: 1, vy: 3, vz: 3},
		moon{x: 1, y: -7, z: 5, vx: -3, vy: 1, vz: -3},
		moon{x: 2, y: 2, z: 0, vx: -1, vy: -3, vz: 1},
	}
	inputMoons := []moon{
		moon{x: -1, y: 0, z: 2, vx: 3, vy: -1, vz: -1},
		moon{x: 2, y: -10, z: -7, vx: 1, vy: 3, vz: 3},
		moon{x: 4, y: -8, z: 8, vx: -3, vy: 1, vz: -3},
		moon{x: 3, y: 5, z: -1, vx: -1, vy: -3, vz: 1},
	}
	actualResult := applyVelocity(&inputMoons)
	if !reflect.DeepEqual(*actualResult, expectedResult) {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_runTimeStep(t *testing.T) {
	expectedResult := []moon{
		moon{x: 2, y: -1, z: 1, vx: 3, vy: -1, vz: -1},
		moon{x: 3, y: -7, z: -4, vx: 1, vy: 3, vz: 3},
		moon{x: 1, y: -7, z: 5, vx: -3, vy: 1, vz: -3},
		moon{x: 2, y: 2, z: 0, vx: -1, vy: -3, vz: 1},
	}
	inputMoons := []moon{
		moon{x: -1, y: 0, z: 2, vx: 0, vy: 0, vz: 0},
		moon{x: 2, y: -10, z: -7, vx: 0, vy: 0, vz: 0},
		moon{x: 4, y: -8, z: 8, vx: 0, vy: 0, vz: 0},
		moon{x: 3, y: 5, z: -1, vx: 0, vy: 0, vz: 0},
	}
	actualResult := runTimeStep(&inputMoons)
	if !reflect.DeepEqual(*actualResult, expectedResult) {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

// TODO: add unit tests for makeHashes, deepCopyMoons, gcd

func Test_IntegrationTest_runTimeStep(t *testing.T) {

	// Test input2 steps 0-10
	for i := 0; i < 10; i++ {
		// Read input2 - step i
		input, err := ioutil.ReadFile(fmt.Sprintf("input2.step%d.txt", i))
		if err != nil {
			panic(err)
		}
		moonsInitialStep := *parseTestInput(&input)

		// Read input2 - step i+1
		input, err = ioutil.ReadFile(fmt.Sprintf("input2.step%d.txt", i+1))
		if err != nil {
			panic(err)
		}

		// Compare expected moons to actual
		moonsFinalStepExpected := *parseTestInput(&input)
		moonsFinalStepActual := runTimeStep(&moonsInitialStep)
		if !reflect.DeepEqual(*moonsFinalStepActual, moonsFinalStepExpected) {
			t.Errorf("Input2, run %d - Expected result was %v, actual result was %v\n", i, moonsFinalStepExpected, moonsFinalStepActual)
		}

	}

	// Test input3 steps 0-100 (increments of 10)
	for i := 0; i < 100; i += 10 {
		// Read input2 - step i
		input, err := ioutil.ReadFile(fmt.Sprintf("input3.step%d.txt", i))
		if err != nil {
			panic(err)
		}
		moonsInitialStep := *parseTestInput(&input)

		// Read input2 - step i+1
		input, err = ioutil.ReadFile(fmt.Sprintf("input3.step%d.txt", i+10))
		if err != nil {
			panic(err)
		}

		// Compare expected moons to actual
		moonsFinalStepExpected := *parseTestInput(&input)
		moonsFinalStepActualPtr := &moonsInitialStep
		for j := 0; j < 10; j++ {
			moonsFinalStepActualPtr = runTimeStep(moonsFinalStepActualPtr)
		}
		if !reflect.DeepEqual(*moonsFinalStepActualPtr, moonsFinalStepExpected) {
			t.Errorf("Input3, run %d - Expected result was %v, actual result was %v\n", i, moonsFinalStepExpected, *moonsFinalStepActualPtr)
		}

	}
}

func Test_EndToEnd_Input(t *testing.T) {
	inputFile := "input.txt"
	iterations := 1000
	expectedResult := 12070
	actualResult := nBodyProblemPart1(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_EndToEnd_Input2(t *testing.T) {
	inputFile := "input2.txt"
	iterations := 10
	expectedResult := 179
	actualResult := nBodyProblemPart1(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_EndToEnd_Input3(t *testing.T) {
	inputFile := "input3.txt"
	iterations := 100
	expectedResult := 1940
	actualResult := nBodyProblemPart1(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_EndToEnd_Part2_Input2(t *testing.T) {
	inputFile := "input2.txt"
	iterations := 2773
	expectedResult := 2772
	actualResult := nBodyProblemPart2(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func Test_EndToEnd_Part2_Input3(t *testing.T) {
	inputFile := "input3.txt"
	iterations := 4686774925
	expectedResult := 4686774924
	actualResult := nBodyProblemPart2(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}
