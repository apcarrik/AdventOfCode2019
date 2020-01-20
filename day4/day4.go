package main

import (
	"fmt"
)



func twoDigitsAdjacent(i int) bool {
	d1 := i%10
	i/=10
	d2 := i%10
	i/=10
	d3 := i%10
	i/=10
	d4 := i%10
	i/=10
	d5 := i%10
	i/=10
	d6 := i%10
	i/=10

	if ( d1 == d2 && d2 != d3) {
		return true
	} else if (d1 != d2 && d2 == d3 && d3 != d4) {
		return true
	} else if (d2 != d3 && d3 == d4 && d4 != d5) {
		return true
	} else if (d3 != d4 && d4 == d5 && d5 != d6) {
		return true
	} else if (d4 != d5 && d5 == d6) {
		return true
	}
	return false
}


// func twoDigitsAdjacent(n int) bool {
// 	max:=4
// 	original := n
// 	last_last_digit := n%10
// 	n /= 10
// 	last_digit := n%10
// 	n /= 10
// 	for i := 0; i<max; i++{
// 		if i == 0 {
// 			if last_last_digit == last_digit && last_digit != n%10 {
// 				return true
// 			}
// 		} else if i == max-1 {
// 			tmp := original /1000
// 			if (tmp != last_last_digit && last_digit == last_last_digit) {
// 				return true
// 			}			
// 		} else {
			
// 		}
// 		last_last_digit = last_digit
// 		last_digit = n%10
// 		n /= 10
// 	}
// 	return false
// }

func digitsNeverDecrease(n int) bool {
	last_digit := n%10
	n /= 10
	for i := 0; i<5; i++{
		if last_digit < n%10 {
			return false
		}
		last_digit = n%10
		n /= 10
	}
	return true
}


func main() {
	range_start := 197487
	range_end := 673251
	count := 0

	for i:= range_start; i<= range_end; i++ {
		if twoDigitsAdjacent(i) {
			if digitsNeverDecrease(i) {
				count += 1
			}
		}
	}

	fmt.Printf("%d found out of %d", count, range_end-range_start)
}
