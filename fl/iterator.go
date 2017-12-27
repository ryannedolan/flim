package fl

import (
	"fmt"
	"reflect"
	"strings"
)

type Iterable interface {
	Next() bool
	Pos() interface{}
}

type Iterator struct {
	Iterable
}

type chanIterable struct {
	ch  chan interface{}
	pos interface{}
}

type node struct {
	next *node
	v    interface{}
}

type listIterable struct {
	head *node
	v    interface{}
}

func (it Iterator) Map(arg interface{}) *Iterator {
	switch arg.(type) {
	case string:
		return it.mapf(f1(arg.(string)))
	default:
		return it.mapf(eraseF1(arg))
	}
}

func (it Iterator) Filter(arg interface{}) *Iterator {
	switch arg.(type) {
	case string:
		return it.filter(f1(arg.(string)))
	default:
		return it.filter(eraseF1(arg))
	}
}

func (it Iterator) Fold(a interface{}, arg interface{}) *Iterator {
	switch arg.(type) {
	case string:
		return it.fold(a, f2(arg.(string)))
	default:
		return it.fold(a, eraseF2(arg))
	}
}

func (it Iterator) Reverse() *Iterator {
	list := EmptyList()
	for it.Next() {
		list.Push(it.Pos())
	}
	return &Iterator{list}
}

func (it Iterator) Floats() []float64 {
	res := make([]float64, 0)
	for it.Next() {
		res = append(res, asFloat64(it.Pos()))
	}
	return res
}

func (it Iterator) Ints() []int {
	res := make([]int, 0)
	for it.Next() {
		res = append(res, asInt(it.Pos()))
	}
	return res
}

func (it Iterator) Strings() []string {
	res := make([]string, 0)
	for it.Next() {
		res = append(res, asString(it.Pos()))
	}
	return res
}

func (it Iterator) Array() []interface{} {
	res := make([]interface{}, 0)
	for it.Next() {
		res = append(res, it.Pos())
	}
	return res
}

func (it Iterator) Chan() chan interface{} {
	ch := make(chan interface{}, 0)
	go func() {
		for it.Next() {
			ch <- it.Pos()
		}
		close(ch)
	}()
	return ch
}

func (it Iterator) String() string {
	switch it.Iterable.(type) {
	case *chanIterable:
		return it.Iterable.(*chanIterable).String()
	case *listIterable:
		return it.Iterable.(*listIterable).String()
	default:
		panic(fmt.Errorf("dunno how to stringify %v", it))
	}
}

func (it *chanIterable) String() string {
	return "Stream(...)"
}

func (it *listIterable) String() string {
	it2 := &Iterator{it}
	return fmt.Sprintf("List(%s)", strings.Join(it2.Strings(), ", "))
}

func (it Iterator) List() *Iterator {
	list := EmptyList()
	it2 := it.Reverse()
	for it2.Next() {
		list.Push(it2.Pos())
	}
	return &Iterator{list}
}

func (it *chanIterable) Next() bool {
	a, ok := <-it.ch
	it.pos = a
	return ok
}

func (it chanIterable) Pos() interface{} {
	return it.pos
}

func (it *listIterable) Next() bool {
	if it.head == nil {
		return false
	}
	it.v = it.head.v
	it.head = it.head.next
	return true
}

func (it *listIterable) Push(v interface{}) {
	n := &node{next: it.head, v: v}
	it.head = n
}

func (it listIterable) Pos() interface{} {
	return it.v
}

func EmptyList() *listIterable {
	return &listIterable{}
}

func Range(i int, j int) *Iterator {
	ch := make(chan interface{}, j-i+1)
	for ; i <= j; i++ {
		ch <- i
	}
	close(ch)
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
		return &Iterator{&chanIterable{ch: ch}}
	case []float64:
		ch := make(chan interface{}, len(arr.([]float64)))
		for _, e := range arr.([]float64) {
			ch <- e
		}
		close(ch)
		return &Iterator{&chanIterable{ch: ch}}
	case []int:
		ch := make(chan interface{}, len(arr.([]int)))
		for _, e := range arr.([]int) {
			ch <- e
		}
		close(ch)
		return &Iterator{&chanIterable{ch: ch}}
	case []string:
		ch := make(chan interface{}, len(arr.([]string)))
		for _, e := range arr.([]string) {
			ch <- e
		}
		close(ch)
		return &Iterator{&chanIterable{ch: ch}}
	case chan interface{}:
		return &Iterator{&chanIterable{ch: arr.(chan interface{})}}
	case chan float64:
		ch := make(chan interface{})
		go func() {
			for x := range arr.(chan float64) {
				ch <- x
			}
			close(ch)
		}()
		return &Iterator{&chanIterable{ch: ch}}
	case chan int:
		ch := make(chan interface{})
		go func() {
			for x := range arr.(chan int) {
				ch <- x
			}
			close(ch)
		}()
		return &Iterator{&chanIterable{ch: ch}}
	case chan string:
		ch := make(chan interface{})
		go func() {
			for x := range arr.(chan string) {
				ch <- x
			}
			close(ch)
		}()
		return &Iterator{&chanIterable{ch: ch}}
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
	return &Iterator{&chanIterable{ch: out}}
}

func f1(expr string) func(interface{}) interface{} {
	f := Lambda(expr)
	return func(x interface{}) interface{} {
		return f(x)
	}
}

func f2(expr string) func(interface{}, interface{}) interface{} {
	f := Lambda(expr)
	return func(x interface{}, y interface{}) interface{} {
		return f(x, y)
	}
}

func eraseF1(f interface{}) func(interface{}) interface{} {
	switch f.(type) {
	case func(interface{}) interface{}:
		return f.(func(interface{}) interface{})
	case func(...interface{}) interface{}:
		return func(x interface{}) interface{} {
			return f.(func(...interface{}) interface{})(x)
		}
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
		panic(fmt.Errorf("dunno how to erase %v", reflect.TypeOf(f)))
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
	return &Iterator{&chanIterable{ch: out}}
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
	return &Iterator{&chanIterable{ch: out}}
}
