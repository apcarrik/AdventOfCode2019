package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"sort"
)

// type LineSegment struct {
// 	start Point
// 	end Point
// }

// func isHorizontal(line *LineSegment) bool {
// 	if line.start.y == line.end.y {
// 		return true
// 	} else if line.start.x == line.end.x {
// 		return false
// 	} else {
// 		panic("Line is neither horizontal nor vertical")
// 	}
// }

// func createWirePath(wireptr *[]string) *[]LineSegment {
// 	var newWire []LineSegment 
// 	wire := *wireptr
// 	lastPt := Point{0,0}
// 	for _,c := range wire{
// 		direction := c[0]
// 		length, err := strconv.Atoi(c[1:])
// 		if err != nil {
// 			panic(err)
// 		}
// 		newLine := LineSegment{start: lastPt, end: lastPt}
// 		switch direction {
// 		case 0x55: // "U"
// 			newLine.end = Point{newLine.start.x, newLine.start.y + length}
// 		case 0x44: // "D"
// 			newLine.end = Point{newLine.start.x, newLine.start.y - length}
// 		case 0x4C: // "L"
// 			newLine.end = Point{newLine.start.x - length, newLine.start.y}
// 		case 0x52: // "R"
// 			newLine.end = Point{newLine.start.x + length, newLine.start.y}
// 		default:
// 			panic(fmt.Sprintf("Direction %s is not recognized", direction))
// 		}
// 		lastPt = newLine.end
// 		newWire = append(newWire, newLine)
// 	}
// 	return &newWire
// }

type Point struct {
	x int
	y int
}

type WirePath struct {
	h []HorizontalLineSegment
	v []VerticalLineSegment
}

type HorizontalLineSegment struct {
	y int
	xmin int
	xmax int
}

type VerticalLineSegment struct {
	x int
	ymin int
	ymax int
}

func createWirePath(wireptr *[]string) *WirePath {
	wire := *wireptr
	var horizontalLines []HorizontalLineSegment
	var verticalLines []VerticalLineSegment
	lastPt := Point{0,0}
	var newPt Point
	for _,c := range wire{
		direction := c[0]
		length, err := strconv.Atoi(c[1:])
		if err != nil {
			panic(err)
		}
		switch direction {
		case 0x55: // "U"
			newPt = Point{lastPt.x, lastPt.y + length}
			verticalLines = append(verticalLines, VerticalLineSegment{x: newPt.x, ymin: lastPt.y, ymax: newPt.y})
		case 0x44: // "D"
			newPt = Point{lastPt.x, lastPt.y - length}
			verticalLines = append(verticalLines, VerticalLineSegment{x: newPt.x, ymin: newPt.y, ymax:lastPt.y})
		case 0x4C: // "L"
			newPt = Point{lastPt.x - length, lastPt.y}
			horizontalLines = append(horizontalLines, HorizontalLineSegment{y: newPt.y, xmin: newPt.x, xmax:lastPt.x})
		case 0x52: // "R"
			newPt = Point{lastPt.x + length, lastPt.y}
			horizontalLines = append(horizontalLines, HorizontalLineSegment{y: newPt.y, xmin: lastPt.x, xmax: newPt.x})
		default:
			panic(fmt.Sprintf("Direction %s is not recognized", direction))
		}
		lastPt = newPt
	}
	newWirePath := WirePath{h: horizontalLines, v: verticalLines}
	return &newWirePath
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanDistance(p *Point) int {
	return Abs(p.x) + Abs(p.y)
}

func pointOnHorizontalLineSegment(p *Point, s *HorizontalLineSegment) bool {
	return p.y == s.y && s.xmin <= p.x && p.x <= s.xmax 
}

func pointOnVerticalLineSegment(p *Point, s *VerticalLineSegment) bool {
	return p.x == s.x && s.ymin <= p.y && p.y <= s.ymax 
}

func getPathDistanceToPoint(i *Point, wireptr *[]string) int {
	wire := *wireptr
	// fmt.Printf("%v", wire)
	pathCounter := 0
	lastPt := Point{0,0}
	var newPt Point
	for _,c := range wire{
		// fmt.Println(c)
		direction := c[0]
		length, err := strconv.Atoi(c[1:])
		if err != nil {
			panic(err)
		}
		switch direction {
		case 0x55: // "U"
			newPt = Point{lastPt.x, lastPt.y + length}			
			if pointOnVerticalLineSegment(i, &VerticalLineSegment{x: newPt.x, ymin: lastPt.y, ymax: newPt.y}) {
				return pathCounter + (i.y - lastPt.y)
			}
		case 0x44: // "D"
			newPt = Point{lastPt.x, lastPt.y - length}		
			if pointOnVerticalLineSegment(i, &VerticalLineSegment{x: newPt.x, ymin: newPt.y, ymax:lastPt.y}) {
				return pathCounter + (lastPt.y - i.y)
			}
		case 0x4C: // "L"
			newPt = Point{lastPt.x - length, lastPt.y}		
			if pointOnHorizontalLineSegment(i, &HorizontalLineSegment{y: newPt.y, xmin: newPt.x, xmax:lastPt.x}) {
				return pathCounter + (lastPt.x - i.x)
			}
		case 0x52: // "R"
			newPt = Point{lastPt.x + length, lastPt.y}	
			if pointOnHorizontalLineSegment(i, &HorizontalLineSegment{y: newPt.y, xmin: lastPt.x, xmax: newPt.x}) {
				return pathCounter + (i.x - lastPt.x)
			}
		default:
			panic(fmt.Sprintf("Direction %s is not recognized", direction))
		}
		pathCounter += length
		// fmt.Printf("%d\n",pathCounter)
		lastPt = newPt
	}
	return -1
}

func main() {
	input, err := ioutil.ReadFile("day3/input.txt")
	if err != nil {
		panic(err)
	}
	in := strings.Split(string(input), "\n")
	line1 := strings.Split(string(in[0]), ",")
	line2 := strings.Split(string(in[1]), ",")

	// Create Path of Wire with Line Segments
	wire1 := *createWirePath(&line1)
	wire2 := *createWirePath(&line2)

	fmt.Printf("Wire1:\nHorizontal:\n%v\n\nVertical:\n%v\n\n", wire1.h, wire1.v)
	fmt.Printf("Wire2:\nHorizontal:\n%v\n\nVertical:\n%v\n\n", wire2.h, wire2.v)
	
	// Sort wire1 horizontal lines by y coordinate, then x max
	sort.SliceStable(wire1.h, func(i, j int) bool {
		if wire1.h[i].y != wire1.h[j].y {
			return wire1.h[i].y < wire1.h[j].y
		} else {
			return wire1.h[i].xmax < wire2.h[j].xmax
		}
	})
	sort.SliceStable(wire1.v, func(i, j int) bool {
		if wire1.v[i].x != wire1.v[j].x {
			return wire1.v[i].x < wire1.v[j].x
		} else {
			return wire1.v[i].ymax < wire1.v[j].ymax
		}
	})

	// Sort wire2 vertical lines by x coordinate, then y max
	sort.SliceStable(wire2.h, func(i, j int) bool {
		if wire2.h[i].y != wire2.h[j].y {
			return wire2.h[i].y < wire2.h[j].y
		} else {
			return wire2.h[i].xmax < wire2.h[j].xmax
		}
	})
	sort.SliceStable(wire2.v, func(i, j int) bool {
		if wire2.v[i].x != wire2.v[j].x {
			return wire2.v[i].x < wire2.v[j].x
		} else {
			return wire2.v[i].ymax < wire2.v[j].ymax
		}
	})
	fmt.Println("Sorted:\n")
	fmt.Printf("Wire1:\nHorizontal:\n%v\n\nVertical:\n%v\n\n", wire1.h, wire1.v)
	fmt.Printf("Wire2:\nHorizontal:\n%v\n\nVertical:\n%v\n\n", wire2.h, wire2.v)

	// For each vertical line in wire1, check horizontal lines of wire2 to see if there is an intersection
	var intersections []Point
	for _,vline := range wire1.v {
		for _,hline := range wire2.h {
			// if vline's x coordinate is within hline's xmin & xmax and hline's y coordinate is within vlines' ymin & ymax, there is intersection
			if ( hline.xmin < vline.x && hline.xmax > vline.x && vline.ymin < hline.y && vline.ymax > hline.y ){
				intersections = append(intersections, Point{x: vline.x, y: hline.y})

			}
		}
	}

	// For each horizontal line in wire1, check vertical lines of wire2 to see if there is an intersection
	for _,hline := range wire1.h {
		for _,vline := range wire2.v {
			// if hline's y coordinate is within vlines's ymin & ymax and vline's x coordinate is within hlines' xmin & xmax, there is intersection
			if ( vline.ymin < hline.y && vline.ymax > hline.y && hline.xmin < vline.x && hline.xmax > vline.x ){
				intersections = append(intersections, Point{x: vline.x, y: hline.y})

			}
		}
	}
	

	fmt.Printf("\nIntersections:\n%v\n", intersections)

	// Sort intersections by manhattan distance
	sort.SliceStable(intersections, func(i, j int) bool {
			return manhattanDistance(&intersections[i]) < manhattanDistance(&intersections[j])
	})

	fmt.Printf("\nIntersections Sorted:\n%v\n", intersections)
	fmt.Printf("\nSmallest Manhattan Distance: %d\n", manhattanDistance(&intersections[0]))


	// Part 2
	fmt.Println("\n\n\n------\nPart2\n\n")
	// fmt.Printf("%s" ,line1)
	// for each intersection, trace the path of the wires and find the combined distance to the intersection for the wires
	min_distance := -1
	for _,i := range intersections {
		distance := getPathDistanceToPoint(&i, &line1)
		// fmt.Printf("%d, ", distance)
		distance += getPathDistanceToPoint(&i, &line2)
		// fmt.Printf("%d\n", distance)
		if (min_distance == -1 || distance < min_distance) {
			min_distance = distance
		}
	}

	fmt.Printf("\nMin distance: %d\n", min_distance)




}