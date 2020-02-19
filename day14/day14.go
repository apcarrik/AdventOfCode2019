package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

// ========================= Structs =========================

type reactant struct {
	chemicalIdx  int // Index of chemical in chemicals array
	quantityUsed int
}

type chemical struct { // TODO: should this be a map?
	name            string
	quantityCreated int
	reactants       []reactant
}

// ===========================================================

// ======================== Functions ========================

func addChemical(chemicalName string, chemicalsPtr *[]chemical) *[]chemical {
	chemicals := *chemicalsPtr
	cIdx := findChemicalIndex(&chemicals, chemicalName)
	if cIdx == -1 {
		// Chemical is new - need to add w/ default values
		chemicals = append(chemicals, chemical{name: chemicalName, quantityCreated: 0, reactants: nil})
	}
	return &chemicals
}

func findChemicalIndex(chemicalsPtr *[]chemical, chemicalName string) int {
	chemicals := *chemicalsPtr
	for i := 0; i < len(chemicals); i++ {
		if chemicals[i].name == chemicalName {
			return i
		}
	}
	return -1
}

func parseInput(inputPtr *[]byte) *[]chemical {
	input := *inputPtr
	chemicals := []chemical{}

	// Store all chemicals
	inputLines := bytes.SplitN(input, []byte("\n"), -1)
	re, err := regexp.Compile(`(\d+) ([A-Z]+)`) //7 A, 1 B => 1 C
	if err != nil {
		panic(err)
	}
	for _, line := range inputLines {
		reactionChemicals := re.FindAllSubmatch(line, -1)
		for i := 0; i < len(reactionChemicals); i++ {
			chemicalsPtr := addChemical(string(reactionChemicals[i][2]), &chemicals)
			chemicals = *chemicalsPtr
		}
	}

	// Add reactant information to chemicals
	for _, line := range inputLines {
		reactionChemicals := re.FindAllSubmatch(line, -1)
		outputChemicalName := string(reactionChemicals[len(reactionChemicals)-1][2])
		outputChemicalQuantityCreated, err := strconv.Atoi(string(reactionChemicals[len(reactionChemicals)-1][1]))
		if err != nil {
			panic(err)
		}
		outputChemicalIdx := findChemicalIndex(&chemicals, outputChemicalName)
		if outputChemicalIdx == -1 {
			panic(fmt.Errorf("error: Chemical %s not found", outputChemicalName))
		}
		chemicals[outputChemicalIdx].quantityCreated = outputChemicalQuantityCreated
		// update reactants
		for i := range reactionChemicals {
			if i < (len(reactionChemicals) - 1) {
				chemicalName := string(reactionChemicals[i][2])
				chemicalQuantityCreated, err := strconv.Atoi(string(reactionChemicals[i][1]))
				if err != nil {
					panic(err)
				}
				chemicalIdx := findChemicalIndex(&chemicals, chemicalName)
				if chemicalIdx == -1 {
					panic(fmt.Errorf("error: Chemical %s not found", chemicalName))
				}
				chemicals[outputChemicalIdx].reactants = append(chemicals[outputChemicalIdx].reactants, reactant{
					chemicalIdx:  chemicalIdx,
					quantityUsed: chemicalQuantityCreated,
				})
			}
		}

	}

	return &chemicals

}

func main() {

	start := time.Now()
	inputFile := "input/input2.txt"

	// Get moons from input file
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	chemicals := parseInput(&input)
	fmt.Printf("Chemicals: %v\n", chemicals)

	elapsed := time.Since(start)
	fmt.Printf("Program took %s\n", elapsed)
}

// ===========================================================
