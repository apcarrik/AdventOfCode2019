package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type OrbitNode struct {
	parent   *OrbitNode
	children []*OrbitNode
	name     string
}

func createNode(l string) *OrbitNode {
	newNode := OrbitNode{
		name: l,
	}
	return &newNode
}

func countOrbits(totalOrbits int, thisOrbitCount int, thisNode *OrbitNode) int {
	totalOrbits += thisOrbitCount
	thisOrbitCount++
	for _, c := range (*thisNode).children {
		totalOrbits = countOrbits(totalOrbits, thisOrbitCount, c)
	}
	return totalOrbits
}

func getNodeAddress(matchName string, thisNode *OrbitNode) *OrbitNode {
	if (*thisNode).name == matchName {
		return thisNode
	}
	for _, c := range (*thisNode).children {
		result := getNodeAddress(matchName, c)
		if result != nil {
			return result
		}
	}
	return nil
}

func getParentTraceToNode(thisNode *OrbitNode, toNode *OrbitNode) *[]string {
	var parentTrace []string
	for (*thisNode).parent != toNode {
		parentTrace = append(parentTrace, (*thisNode).parent.name)
		thisNode = (*thisNode).parent
	}
	parentTrace = append(parentTrace, (*thisNode).parent.name)
	return &parentTrace
}

func findFirstIntersection(YOUParentTrace *[]string, SANParentTrace *[]string) string {
	for _, y := range *YOUParentTrace {
		for _, s := range *SANParentTrace {
			if y == s {
				return y
			}
		}
	}
	return ""
}

func main() {
	input, err := ioutil.ReadFile("day6/input.txt")
	if err != nil {
		panic(err)
	}
	orbitLines := strings.Split(string(input), "\n")

	fmt.Printf("# lines: %d\n", len(orbitLines))

	var nodes []*OrbitNode
	nodes = append(nodes, &OrbitNode{name: "COM"})
	// loop through lines and add new node
	for _, l := range orbitLines[0:] {
		nodeNames := strings.Split(l, ")")
		parentFound := false
		childFound := false
		var parentNode *OrbitNode
		var childNode *OrbitNode
		// loop through existing nodes to see if parent or child exists. if not, add them to list of nodes
		for _, n := range nodes {
			if n.name == strings.TrimSpace(nodeNames[0]) {
				parentFound = true
				parentNode = n
			} else if n.name == strings.TrimSpace(nodeNames[1]) {
				childFound = true
				childNode = n
			}
		}
		if !parentFound {
			parentNode = createNode(nodeNames[0])
			nodes = append(nodes, parentNode)
		}
		if !childFound {
			childNode = createNode(nodeNames[1])
			nodes = append(nodes, childNode)
		}
		parentNode.children = append((*parentNode).children, childNode)
		childNode.parent = parentNode

		// fmt.Printf("Parent Node: %v\n", parentNode)
	}

	fmt.Printf("nodes: %d\n", len(nodes))

	// Part 1: Count direct and indirect orbits
	orbitCount := countOrbits(0, 0, nodes[0])
	fmt.Printf("total orbits: %d\n", orbitCount)

	// Part 2: Find minimum number of orbital transfers
	// Get parent trace for both YOU and SAN and find duplicate node
	YOUNode := getNodeAddress("YOU", nodes[0])
	SANNode := getNodeAddress("SAN", nodes[0])

	YOUParentTrace := getParentTraceToNode(YOUNode, nodes[0])
	SANParentTrace := getParentTraceToNode(SANNode, nodes[0])
	commonNode := getNodeAddress(findFirstIntersection(YOUParentTrace, SANParentTrace), nodes[0])
	YOUTrace := *getParentTraceToNode(YOUNode, commonNode)
	SANTrace := *getParentTraceToNode(SANNode, commonNode)

	fmt.Printf("Minimum orbital transfers: %d\n", len(YOUTrace)+len(SANTrace)-2)

}
