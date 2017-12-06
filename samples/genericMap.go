package main

import (
	"fmt"
	"github.com/ryannedolan/flim/fl"
	"github.com/ryannedolan/flim/lambda"
)

// demonstrates use of lambdas and closures to operate on a list
func main() {
	arr := []int{1, 2, 3, 4, 5}

	a := fl.Iter(arr).Map(`x + 1`).Force()
	fmt.Println(a)

	b := fl.Iter(a).Map(func(x int) int { return x + 1 }).Force()
	fmt.Println(b)

	c := fl.Iter(b).Map(func(x float64) float64 { return x + 1.0 }).Force()
	fmt.Println(c)

	d := fl.Iter(c).Map(lambda.X(`x + 1`)).Force()
	fmt.Println(d)
}
