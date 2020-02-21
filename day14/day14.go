package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
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
	reactionYeild   int
	quantityCreated int
	quantityExtra   int
	reactants       []reactant
}

// ===========================================================

// ======================== Functions ========================

func addChemical(chemicalName string, chemicalsPtr *[]chemical) *[]chemical {
	chemicals := *chemicalsPtr
	cIdx := findChemicalIndex(&chemicals, chemicalName)
	if cIdx == -1 {
		// Chemical is new - need to add w/ default values
		chemicals = append(chemicals, chemical{name: chemicalName})
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
		outputChemicalYeild, err := strconv.Atoi(string(reactionChemicals[len(reactionChemicals)-1][1]))
		if err != nil {
			panic(err)
		}
		outputChemicalIdx := findChemicalIndex(&chemicals, outputChemicalName)
		if outputChemicalIdx == -1 {
			panic(fmt.Errorf("error: Chemical %s not found", outputChemicalName))
		}
		chemicals[outputChemicalIdx].reactionYeild = outputChemicalYeild
		// update reactants
		for i := range reactionChemicals {
			if i < (len(reactionChemicals) - 1) {
				chemicalName := string(reactionChemicals[i][2])
				chemicalQuantityUsed, err := strconv.Atoi(string(reactionChemicals[i][1]))
				if err != nil {
					panic(err)
				}
				chemicalIdx := findChemicalIndex(&chemicals, chemicalName)
				if chemicalIdx == -1 {
					panic(fmt.Errorf("error: Chemical %s not found", chemicalName))
				}
				chemicals[outputChemicalIdx].reactants = append(chemicals[outputChemicalIdx].reactants, reactant{
					chemicalIdx:  chemicalIdx,
					quantityUsed: chemicalQuantityUsed,
				})
			}
		}

	}

	return &chemicals

}

func getChemicalRequired(chemicalsPtr *[]chemical, chemicalIdx int, amountRequired int) *[]chemical {
	chemicals := *chemicalsPtr
	amountRequiredToMake := (amountRequired - chemicals[chemicalIdx].quantityExtra)
	yeild := chemicals[chemicalIdx].reactionYeild
	if yeild == 0 {
		yeild = 1
	}
	if amountRequiredToMake > 0 {
		reactionMultiplier := int(math.Ceil(float64(amountRequiredToMake) / float64(yeild)))
		// need to make more of output - look at reactants
		for _, reactant := range chemicals[chemicalIdx].reactants {
			// get amount of each reactant
			amountOfReactantRequired := reactant.quantityUsed * reactionMultiplier
			chemicalsPtr = getChemicalRequired(&chemicals, reactant.chemicalIdx, amountOfReactantRequired)
			chemicals = *chemicalsPtr
		}
		// need to update the amount of this chemical that is used & extra
		chemicals[chemicalIdx].quantityCreated += (yeild * reactionMultiplier)
		chemicals[chemicalIdx].quantityExtra = (yeild * reactionMultiplier) - amountRequiredToMake
	} else {
		chemicals[chemicalIdx].quantityExtra -= amountRequired
	}
	return &chemicals

}

func getOreUsedForFuel(chemicalsPtr *[]chemical, amountOfFuel int) int {
	chemicalsPtr = getChemicalRequired(chemicalsPtr, findChemicalIndex(chemicalsPtr, "FUEL"), amountOfFuel)
	return (*chemicalsPtr)[findChemicalIndex(chemicalsPtr, "ORE")].quantityCreated
}

func part1(inputFile string) int {
	// Get moons from input file
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	chemicals := *parseInput(&input)

	// Get total amount of ORE used to create FUEL
	return getOreUsedForFuel(&chemicals, 1)
}

func part2(inputFile string) int {
	// Get moons from input file
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	chemicals := *parseInput(&input)

	// TODO: figure out how to iterate the fuel use to find the right amount of ore
	amountOfFuel := 1
	// Get total amount of ORE used to create FUEL
	return getOreUsedForFuel(&chemicals, amountOfFuel)

}

func main() {
	start := time.Now()
	inputFile := "input/input.txt"

	// Part 1
	oreUsedForFuel := part1(inputFile)
	fmt.Printf("Ore used for fuel: %d\n", oreUsedForFuel)

	// Part 2
	fuelWithOneTrillionOre := part2(inputFile)
	fmt.Printf("Amount of fuel created with one trillion ore: %d\n", fuelForOneTrillionOre)

	elapsed := time.Since(start)
	fmt.Printf("Program took %s\n", elapsed)
}

// ===========================================================
