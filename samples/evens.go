package main

import (
	"fmt"
	"github.com/ryannedolan/flim/fl"
	"math"
	"time"
)

// slowly prints even numbers
func main() {

	evens := fl.Range(1, 100).FilterF(
		func(x interface{}) bool {
			time.Sleep(100000000)
			return math.Remainder(float64(x.(int)), float64(2)) == 0
		}).Chan()

	for i := range evens {
		fmt.Println(i)
	}
}
