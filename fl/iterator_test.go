package fl

import (
	"strings"
	"testing"
)

func assertEquals(t *testing.T, a []interface{}, b []interface{}) {
	if len(a) != len(b) {
		t.Fatalf("%v should have length %d", a, len(b))
	}
	for i, e := range a {
		if e != b[i] {
			t.Fatalf("%v should be %v", a, b)
		}
	}
}

func TestRange(t *testing.T) {
	iter := Range(0, 5)
	elems := iter.Force()
	assertEquals(t, elems, []interface{}{0, 1, 2, 3, 4, 5})
}

func TestFloatMap(t *testing.T) {
	arr := []float64{1.0, 2.0, 3.0, 4.0}
	res := Iter(arr).Map(`x + 1`).Force()
	assertEquals(t, res, []interface{}{2.0, 3.0, 4.0, 5.0})
}

func TestStringFilter(t *testing.T) {
	arr := []interface{}{"foo", "bar", "car", "far"}
	res := Iter(arr).Filter(func(x string) bool {
		return strings.HasPrefix(x, "f")
	}).Force()
	assertEquals(t, res, []interface{}{"foo", "far"})
}

func TestFloatFilter(t *testing.T) {
	arr := []float64{1.0, 2.0, 3.0, 4.0}
	res := Iter(arr).Filter(` x >= 3.0 `).Force()
	assertEquals(t, res, []interface{}{3.0, 4.0})
}

/*
func TestFloatMapFilterFold(t *testing.T) {
	res := Range(1, 4).Filter(` x <= 3 `).Map(` x*x `).Fold(0, ` x + y `).Force()
	assertEquals(t, res, []interface{}{14})
}*/
