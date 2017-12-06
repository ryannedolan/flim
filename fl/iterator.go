package fl

import (
	"fmt"
)

type Iterable interface {
	Next() bool
	Pos() interface{}
}

type Iterator struct {
	Iterable
}

type chanIterable struct {
	*Iterator
	ch  chan interface{}
	pos interface{}
}

func (it *Iterator) Map(arg interface{}) *Iterator {
	switch arg.(type) {
	case string:
		return it.mapf(f1(arg.(string)))
	default:
		return it.mapf(eraseF1(arg))
	}
}

func (it *Iterator) Filter(arg interface{}) *Iterator {
	switch arg.(type) {
	case string:
		return it.filter(f1(arg.(string)))
	default:
		return it.filter(eraseF1(arg))
	}
}

func (it *Iterator) Fold(a interface{}, arg interface{}) *Iterator {
	switch arg.(type) {
	case string:
		return it.fold(a, f2(arg.(string)))
	default:
		return it.fold(a, eraseF2(arg))
	}
}

func (it *Iterator) Floats() []float64 {
	res := make([]float64, 0)
	for it.Next() {
		res = append(res, asFloat64(it.Pos()))
	}
	return res
}

func (it *Iterator) Force() []interface{} {
	res := make([]interface{}, 0)
	for it.Next() {
		res = append(res, it.Pos())
	}
	return res
}

func (it *Iterator) Chan() chan interface{} {
	ch := make(chan interface{}, 0)
	go func() {
		for it.Next() {
			ch <- it.Pos()
		}
		close(ch)
	}()
	return ch
}

func (it *chanIterable) Next() bool {
	a, ok := <-it.ch
	it.pos = a
	return ok
}

func (it chanIterable) Pos() interface{} {
	return it.pos
}

func Range(i int, j int) *Iterator {
	ch := make(chan interface{}, j-i+1)
	for ; i <= j; i++ {
		ch <- i
	}
	close(ch)
	return &Iterator{&chanIterable{ch: ch}}
}

func newChanIterator(ch chan interface{}) *Iterator {
	return &Iterator{&chanIterable{ch: ch}}
}

func Iter(arr interface{}) *Iterator {
	switch arr.(type) {
	case []interface{}:
		ch := make(chan interface{}, len(arr.([]interface{})))
		for _, e := range arr.([]interface{}) {
			ch <- e
		}
		close(ch)
		return newChanIterator(ch)
	case []float64:
		ch := make(chan interface{}, len(arr.([]float64)))
		for _, e := range arr.([]float64) {
			ch <- e
		}
		close(ch)
		return newChanIterator(ch)
	case []int:
		ch := make(chan interface{}, len(arr.([]int)))
		for _, e := range arr.([]int) {
			ch <- e
		}
		close(ch)
		return newChanIterator(ch)
	case []string:
		ch := make(chan interface{}, len(arr.([]string)))
		for _, e := range arr.([]string) {
			ch <- e
		}
		close(ch)
		return newChanIterator(ch)
	default:
		panic(fmt.Errorf("dunno how to iterate over %v", arr))
	}
}

func (it *Iterator) mapf(f func(interface{}) interface{}) *Iterator {
	out := make(chan interface{})
	go func() {
		for it.Next() {
			out <- f(it.Pos())
		}
		close(out)
	}()
	return newChanIterator(out)
}

func f1(expr string) func(interface{}) interface{} {
	f := Lambda(expr)
	return func(x interface{}) interface{} {
		return f(X(x))
	}
}

func f2(expr string) func(interface{}, interface{}) interface{} {
	f := Lambda(expr)
	return func(x interface{}, y interface{}) interface{} {
		return f(XY(x, y))
	}
}

func eraseF1(f interface{}) func(interface{}) interface{} {
	switch f.(type) {
	case func(interface{}) interface{}:
		return f.(func(interface{}) interface{})
	case func(float64) float64:
		return func(x interface{}) interface{} {
			return f.(func(float64) float64)(asFloat64(x))
		}
	case func(int) int:
		return func(x interface{}) interface{} {
			return f.(func(int) int)(asInt(x))
		}
	case func(float64) bool:
		return func(x interface{}) interface{} {
			return f.(func(float64) bool)(asFloat64(x))
		}
	case func(int) bool:
		return func(x interface{}) interface{} {
			return f.(func(int) bool)(asInt(x))
		}
	case func(string) bool:
		return func(x interface{}) interface{} {
			return f.(func(string) bool)(x.(string))
		}
	case func(interface{}) bool:
		return func(x interface{}) interface{} {
			return f.(func(interface{}) bool)(x)
		}
	default:
		panic(fmt.Errorf("dunno how to erase %v", f))
	}
}

func eraseF2(f interface{}) func(interface{}, interface{}) interface{} {
	switch f.(type) {
	case func(interface{}, interface{}) interface{}:
		return f.(func(interface{}, interface{}) interface{})
	case func(float64, float64) float64:
		return func(x interface{}, y interface{}) interface{} {
			return f.(func(float64, float64) float64)(asFloat64(x), asFloat64(y))
		}
	case func(int, int) int:
		return func(x interface{}, y interface{}) interface{} {
			return f.(func(int, int) int)(asInt(x), asInt(y))
		}
	case func(float64, float64) bool:
		return func(x interface{}, y interface{}) interface{} {
			return f.(func(float64, float64) bool)(asFloat64(x), asFloat64(y))
		}
	case func(int, int) bool:
		return func(x interface{}, y interface{}) interface{} {
			return f.(func(int, int) bool)(asInt(x), asInt(y))
		}
	case func(string, string) bool:
		return func(x interface{}, y interface{}) interface{} {
			return f.(func(string, string) bool)(x.(string), y.(string))
		}
	case func(interface{}, interface{}) bool:
		return func(x interface{}, y interface{}) interface{} {
			return f.(func(interface{}, interface{}) bool)(x, y)
		}
	default:
		panic(fmt.Errorf("dunno how to erase %v", f))
	}
}

func (it *Iterator) filter(f func(interface{}) interface{}) *Iterator {
	out := make(chan interface{})
	go func() {
		for it.Next() {
			x := it.Pos()
			if f(x).(bool) {
				out <- x
			}
		}
		close(out)
	}()
	return newChanIterator(out)
}

func (it *Iterator) fold(z interface{}, f func(x interface{}, y interface{}) interface{}) *Iterator {
	out := make(chan interface{})
	go func() {
		for it.Next() {
			x := it.Pos()
			z = f(z, x)
		}
		out <- z
		close(out)
	}()
	return newChanIterator(out)
}
