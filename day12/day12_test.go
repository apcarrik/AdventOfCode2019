package main

import (
	"testing"
	"reflect"
)

func TestParseInput(t *testing.T) {
	expectedResult := []moon{
		moon{x:-3, y:15, z:-11, vx:0, vy:0, vz:0},
		moon{x:3, y:13, z:-19, vx:0, vy:0, vz:0},
		moon{x:-13, y:18, z:-2, vx:0, vy:0, vz:0},
		moon{x:6, y:0, z:-1, vx:0, vy:0, vz:0},
	}
	inputByteArray := []byte("<x=-3, y=15, z=-11>\n<x=3, y=13, z=-19>\n<x=-13, y=18, z=-2>\n<x=6, y=0, z=-1>")
	actualResult := parseInput(&inputByteArray)
	if !reflect.DeepEqual(*actualResult, expectedResult) {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func TestResultInput(t *testing.T) {
	inputFile := "input.txt"
	iterations := 1000
	expectedResult := 0
	actualResult := nBodyProblem(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func TestResultInput2(t *testing.T) {
	inputFile := "input2.txt"
	iterations := 1000
	expectedResult := 0
	actualResult := nBodyProblem(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}

func TestResultInput3(t *testing.T) {
	inputFile := "input3.txt"
	iterations := 1000
	expectedResult := 0
	actualResult := nBodyProblem(inputFile, iterations)
	if actualResult != expectedResult {
		t.Errorf("Expected result was %v, actual result was %v\n", expectedResult, actualResult)
	}
}