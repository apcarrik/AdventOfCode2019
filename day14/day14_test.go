package main

import (
	"reflect"
	"testing"
)

func Test_UnitTest_addChemical(t *testing.T) {
	testChemicals := []chemical{
		chemical{name: "KJY", reactionYeild: 0, reactants: nil},
		chemical{name: "SNT", reactionYeild: 0, reactants: nil},
		chemical{name: "WER", reactionYeild: 0, reactants: nil},
	}
	newChemicalName := "QID"
	expectedResult := &[]chemical{
		chemical{name: "KJY", reactionYeild: 0, reactants: nil},
		chemical{name: "SNT", reactionYeild: 0, reactants: nil},
		chemical{name: "WER", reactionYeild: 0, reactants: nil},
		chemical{name: "QID", reactionYeild: 0, reactants: nil},
	}
	actualResult := addChemical(newChemicalName, &testChemicals)
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}

}

func Test_UnitTest_findChemicalIndex(t *testing.T) {
	testChemicals := []chemical{
		chemical{name: "KJY", reactionYeild: 0, reactants: nil},
		chemical{name: "SNT", reactionYeild: 0, reactants: nil},
		chemical{name: "WER", reactionYeild: 0, reactants: nil},
	}
	expectedResult := 1
	actualResult := findChemicalIndex(&testChemicals, "SNT")
	if expectedResult != actualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_parseInput(t *testing.T) {
	// Creating expectedResult chemicals
	chemORE := chemical{name: "ORE", reactionYeild: 0, reactants: nil}
	chemA := chemical{name: "A", reactionYeild: 10, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 10}}}
	chemB := chemical{name: "B", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 1}}}
	chemC := chemical{name: "C", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 2, quantityUsed: 1}}}
	chemD := chemical{name: "D", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 3, quantityUsed: 1}}}
	chemE := chemical{name: "E", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 4, quantityUsed: 1}}}
	chemFUEL := chemical{name: "FUEL", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 5, quantityUsed: 1}}}

	expectedResult := &[]chemical{chemORE, chemA, chemB, chemC, chemD, chemE, chemFUEL}
	inputByteArray := []byte("10 ORE => 10 A\n1 ORE => 1 B\n7 A, 1 B => 1 C\n7 A, 1 C => 1 D\n7 A, 1 D => 1 E\n7 A, 1 E => 1 FUEL") // input 2
	actualResult := parseInput(&inputByteArray)
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_getChemicalRequired(t *testing.T) {
	inputChemORE := chemical{name: "ORE", reactionYeild: 0, reactants: nil}
	inputChemA := chemical{name: "A", reactionYeild: 10, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 10}}}
	inputChemB := chemical{name: "B", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 1}}}
	inputChemC := chemical{name: "C", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 2, quantityUsed: 1}}}
	inputChemD := chemical{name: "D", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 3, quantityUsed: 1}}}
	inputChemE := chemical{name: "E", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 4, quantityUsed: 1}}}
	inputChemFUEL := chemical{name: "FUEL", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 5, quantityUsed: 1}}}
	inputChemicals := []chemical{inputChemORE, inputChemA, inputChemB, inputChemC, inputChemD, inputChemE, inputChemFUEL}

	expectedChemORE := chemical{name: "ORE", reactionYeild: 0, quantityCreated: 31, quantityExtra: 0, reactants: nil}
	expectedChemA := chemical{name: "A", reactionYeild: 10, quantityCreated: 30, quantityExtra: 2, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 10}}}
	expectedChemB := chemical{name: "B", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 1}}}
	expectedChemC := chemical{name: "C", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 2, quantityUsed: 1}}}
	expectedChemD := chemical{name: "D", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 3, quantityUsed: 1}}}
	expectedChemE := chemical{name: "E", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 4, quantityUsed: 1}}}
	expectedChemFUEL := chemical{name: "FUEL", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 5, quantityUsed: 1}}}
	expectedChemicals := []chemical{expectedChemORE, expectedChemA, expectedChemB, expectedChemC, expectedChemD, expectedChemE, expectedChemFUEL}

	expectedResult := expectedChemicals
	actualResult := *getChemicalRequired(&inputChemicals, 6, 1)
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_getOreUsedForFuel(t *testing.T) {
	chemORE := chemical{name: "ORE", reactionYeild: 0, reactants: nil}
	chemA := chemical{name: "A", reactionYeild: 10, quantityCreated: 30, quantityExtra: 2, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 10}}}
	chemB := chemical{name: "B", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 1}}}
	chemC := chemical{name: "C", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 2, quantityUsed: 1}}}
	chemD := chemical{name: "D", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 3, quantityUsed: 1}}}
	chemE := chemical{name: "E", reactionYeild: 1, quantityCreated: 1, quantityExtra: 0, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 4, quantityUsed: 1}}}
	chemFUEL := chemical{name: "FUEL", reactionYeild: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 5, quantityUsed: 1}}}
	testChemicals := []chemical{chemORE, chemA, chemB, chemC, chemD, chemE, chemFUEL}

	expectedResult := 31
	actualResult := getOreUsedForFuel(&testChemicals, 1)
	if expectedResult != actualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}
}

func Test_EndToEnd_input(t *testing.T) {
	part1ExpectedResult := 654909
	part1ActualResult := part1("input/input.txt")
	if part1ExpectedResult != part1ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part1ExpectedResult, part1ActualResult)
	}
	part2ExpectedResult := 2876992
	part2ActualResult := part2("input/input.txt")
	if part2ExpectedResult != part2ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part2ExpectedResult, part2ActualResult)
	}
}

func Test_EndToEnd_input2(t *testing.T) {
	part1ExpectedResult := 31
	part1ActualResult := part1("input/input2.txt")
	if part1ExpectedResult != part1ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part1ExpectedResult, part1ActualResult)
	}
}

func Test_EndToEnd_input3(t *testing.T) {
	part1ExpectedResult := 165
	part1ActualResult := part1("input/input3.txt")
	if part1ExpectedResult != part1ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part1ExpectedResult, part1ActualResult)
	}
}

func Test_EndToEnd_input4(t *testing.T) {
	part1ExpectedResult := 13312
	part1ActualResult := part1("input/input4.txt")
	if part1ExpectedResult != part1ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part1ExpectedResult, part1ActualResult)
	}
	part2ExpectedResult := 82892753
	part2ActualResult := part2("input/input4.txt")
	if part2ExpectedResult != part2ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part2ExpectedResult, part2ActualResult)
	}
}

func Test_EndToEnd_input5(t *testing.T) {
	part1ExpectedResult := 180697
	part1ActualResult := part1("input/input5.txt")
	if part1ExpectedResult != part1ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part1ExpectedResult, part1ActualResult)
	}
	part2ExpectedResult := 5586022
	part2ActualResult := part2("input/input5.txt")
	if part2ExpectedResult != part2ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part2ExpectedResult, part2ActualResult)
	}
}

func Test_EndToEnd_input6(t *testing.T) {
	part1ExpectedResult := 2210736
	part1ActualResult := part1("input/input6.txt")
	if part1ExpectedResult != part1ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part1ExpectedResult, part1ActualResult)
	}
	part2ExpectedResult := 460664
	part2ActualResult := part2("input/input6.txt")
	if part2ExpectedResult != part2ActualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", part2ExpectedResult, part2ActualResult)
	}
}
