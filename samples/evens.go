package main

import (
	"fmt"
	"github.com/ryannedolan/flim/fl"
	"math"
	"time"
)

// slowly prints even numbers
func main() {

	evens := fl.Range(1, 100).Filter(
		func(x float64) bool {
			time.Sleep(100000000)
			return math.Remainder(x, 2) == 0
		}).Chan()

	for i := range evens {
		fmt.Println(i)
	}
}
