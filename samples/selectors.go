package main

import (
	"fmt"
	"github.com/ryannedolan/flim/fl"
)

type Person struct {
	name string
	Age  int
	F    func() string
}

func (p Person) Name() string {
	return p.name
}

func foo() string {
	return "ryanne"
}

func main() {
	ry := Person{"ryanne", 33, func() string { return "F!" }}
	f := fl.Lambda(`x.Age`)
	g := fl.Lambda(`x.F()`)
	h := fl.Lambda(`x.Name()`)
	b := fl.Lambda(`x.Name() == "foo"`)
	c := fl.Lambda(`x.Name() == y()`)
	fmt.Println(f(ry))
	fmt.Println(g(ry))
	fmt.Println(h(ry))
	fmt.Println(b(ry))
	fmt.Println(c(ry, foo))
}
