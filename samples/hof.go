package main

import (
	"fmt"
	"github.com/ryannedolan/flim/fl"
)

func main() {
	f := fl.Lambda(`x+1`)
	g := fl.Lambda(`x+2`)
	h := fl.Lambda(`x+y`)

	fog := fl.Compose(f, g)
	fmt.Println("fog(1) = ", fog(1))

	cur := fl.Curry(h, 1)
	fmt.Println("h(1)(2) = ", cur(2))

	cur2 := fl.Curry(h, 1, 2)
	fmt.Println("h(1,2)() = ", cur2())

	cur3 := fl.Curry(h)
	fmt.Println("h()(1,2) = ", cur3(1, 2))
}
