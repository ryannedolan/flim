package fl

import (
	"testing"
)

func TestLiteralAddition(t *testing.T) {
	a := Eval(` 1 + 2 `, nil)
	if a.(int) != 3 {
		t.Fatal("1 + 2 should be 3")
	}

	b := Eval(` 1.0 + 2.0 `, nil)
	if b.(float64) != 3.0 {
		t.Fatal("1.0 + 2.0 should be 3.0")
	}
}

func TestLiteralConcatination(t *testing.T) {
	a := Eval(` "hello" + "world" `, nil)
	if a.(string) != "helloworld" {
		t.Fatalf("hello + world should be helloworld; got %v", a)
	}
}

func TestXYAddition(t *testing.T) {
	a := Eval(` x + y `, XY(1, 2))
	if a.(int) != 3 {
		t.Fatal("x + y should be 3")
	}
}
