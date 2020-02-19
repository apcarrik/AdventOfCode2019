package main

import (
	"reflect"
	"testing"
)

func Test_UnitTest_addChemical(t *testing.T) {
	testChemicals := []chemical{
		chemical{name: "KJY", quantityCreated: 0, reactants: nil},
		chemical{name: "SNT", quantityCreated: 0, reactants: nil},
		chemical{name: "WER", quantityCreated: 0, reactants: nil},
	}
	newChemicalName := "QID"
	expectedResult := &[]chemical{
		chemical{name: "KJY", quantityCreated: 0, reactants: nil},
		chemical{name: "SNT", quantityCreated: 0, reactants: nil},
		chemical{name: "WER", quantityCreated: 0, reactants: nil},
		chemical{name: "QID", quantityCreated: 0, reactants: nil},
	}
	actualResult := addChemical(newChemicalName, &testChemicals)
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}

}

func Test_UnitTest_findChemicalIndex(t *testing.T) {
	testChemicals := []chemical{
		chemical{name: "KJY", quantityCreated: 0, reactants: nil},
		chemical{name: "SNT", quantityCreated: 0, reactants: nil},
		chemical{name: "WER", quantityCreated: 0, reactants: nil},
	}
	expectedResult := 1
	actualResult := findChemicalIndex(&testChemicals, "SNT")
	if expectedResult != actualResult {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}
}

func Test_UnitTest_parseInput(t *testing.T) {
	// Creating expectedResult chemicals
	chemORE := chemical{name: "ORE", quantityCreated: 0, reactants: nil}
	chemA := chemical{name: "A", quantityCreated: 10, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 10}}}
	chemB := chemical{name: "B", quantityCreated: 1, reactants: []reactant{reactant{chemicalIdx: 0, quantityUsed: 1}}}
	chemC := chemical{name: "C", quantityCreated: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 2, quantityUsed: 1}}}
	chemD := chemical{name: "D", quantityCreated: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 3, quantityUsed: 1}}}
	chemE := chemical{name: "E", quantityCreated: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 4, quantityUsed: 1}}}
	chemFUEL := chemical{name: "FUEL", quantityCreated: 1, reactants: []reactant{reactant{chemicalIdx: 1, quantityUsed: 7}, reactant{chemicalIdx: 5, quantityUsed: 1}}}

	expectedResult := &[]chemical{chemORE, chemA, chemB, chemC, chemD, chemE, chemFUEL}
	inputByteArray := []byte("10 ORE => 10 A\n1 ORE => 1 B\n7 A, 1 B => 1 C\n7 A, 1 C => 1 D\n7 A, 1 D => 1 E\n7 A, 1 E => 1 FUEL") // input 2
	actualResult := parseInput(&inputByteArray)
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Errorf("Expected result was:\n%+v\n actual result was:\n%+v\n", expectedResult, actualResult)
	}
}
