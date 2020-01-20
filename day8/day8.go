package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func splitIntoLayers(inPtr *[]byte, layerSize int) *[][]byte {
	var ret [][]byte
	in := *inPtr
	for i := 0; i < len(in); i += layerSize {
		ret = append(ret, in[i:i+layerSize])
	}
	return &ret
}

func countDigits(layers *[][]byte) *[][]int {
	var ret [][]int
	for _, layer := range *layers {
		digitArr := make([]int, 3)
		for _, c := range layer {
			cInt, err := strconv.Atoi(string([]byte{c}))
			if err != nil {
				panic(err)
			}
			// fmt.Printf("%d", cInt)
			digitArr[cInt]++
		}
		ret = append(ret, digitArr)
	}
	return &ret
}

func makeFinalLayer(layers *[][]byte, layerSize int) *[]int {
	filledCounter := 0
	filled := make([]int, layerSize)
	for _, l := range *layers {
		for ic, c := range l {
			cInt, err := strconv.Atoi(string([]byte{c}))
			if err != nil {
				panic(err)
			}
			if (cInt+1 == 1 || cInt+1 == 2) && filled[ic] == 0 {
				filled[ic] = cInt + 1
				filledCounter++
			}
		}
		if filledCounter == layerSize {
			return &filled
		}
	}
	return &filled
}

func printLayer(layerPtr *[]int, pixelWidth int) {
	layer := *layerPtr
	for i := 0; i < len(layer); i += pixelWidth {
		fmt.Println(layer[i : i+pixelWidth])
	}
}

func main() {
	input, err := ioutil.ReadFile("day8/input.txt")
	if err != nil {
		panic(err)
	}
	const pixelWidth = 25
	const pixelHeight = 6

	// Part 1
	layers := splitIntoLayers(&input, pixelWidth*pixelHeight)
	digitSums := countDigits(layers)
	minLayer := -1
	minZeros := 0
	for i, l := range *digitSums {
		if (minLayer == -1) || (l[0] < minZeros) {
			minZeros = l[0]
			minLayer = i
		}
	}
	output := (*digitSums)[minLayer][1] * (*digitSums)[minLayer][2]
	fmt.Printf("minZeros:%d\nminLayer:%d\noutput:%d\n", minZeros, minLayer, output)

	// Part 2
	// go over every layer and make final layer until no transparent spaces left
	finalLayer := makeFinalLayer(layers, pixelWidth*pixelHeight)
	printLayer(finalLayer, pixelWidth)

}
