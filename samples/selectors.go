package main

import (
	"fmt"
	"github.com/ryannedolan/flim/lambda"
)

type Person struct {
	name string
	Age  int
	F    func() string
}

func (p Person) Name() string {
	return p.name
}

func main() {
	ry := Person{"ryanne", 33, func() string { return "F!" }}
	f := lambda.X(`x.Age`)
	g := lambda.X(`x.F()`)
	h := lambda.X(`x.Name()`)
	fmt.Println(f(ry))
	fmt.Println(g(ry))
	fmt.Println(h(ry))
}
