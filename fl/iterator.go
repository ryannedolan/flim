package fl

type Iterable interface {
	Next() bool
	Pos() interface{}
	MapF(func(x interface{}) interface{}) *Iterator
	Force() []interface{}
	Floats() []float64
}

type Iterator struct {
	Iterable
}

type index struct {
	i int
	j int
}

type chanIterable struct {
	*Iterator
	ch  chan interface{}
	pos interface{}
}

func (it *Iterator) Map(expr string) *Iterator {
	f := Lambda(expr)
	return it.MapF(func(x interface{}) interface{} {
		return f(X(x))
	})
}

func (it *Iterator) Filter(expr string) *Iterator {
	f := Lambda(expr)
	return it.FilterF(func(x interface{}) bool {
		return asBool(f(X(x)))
	})
}

func (it *Iterator) Fold(z interface{}, expr string) *Iterator {
	f := Lambda(expr)
	return it.FoldF(z, func(x interface{}, y interface{}) interface{} {
		return f(XY(x, y))
	})
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

func (i *index) Next() bool {
	i.i += 1
	return i.i <= i.j
}

func (i *index) Pos() interface{} {
	return i.i
}

func newChanIterator(ch chan interface{}) *Iterator {
	return &Iterator{&chanIterable{ch: ch}}
}

func Iter(arr []interface{}) *Iterator {
	ch := make(chan interface{}, len(arr))
	for _, e := range arr {
		ch <- e
	}
	close(ch)
	return newChanIterator(ch)
}

func Floats(arr []float64) *Iterator {
	ch := make(chan interface{}, len(arr))
	for _, e := range arr {
		ch <- e
	}
	close(ch)
	return newChanIterator(ch)
}

func (it *Iterator) MapF(f func(interface{}) interface{}) *Iterator {
	out := make(chan interface{})
	go func() {
		for it.Next() {
			out <- f(it.Pos())
		}
		close(out)
	}()
	return newChanIterator(out)
}

func (it *Iterator) FilterF(f func(interface{}) bool) *Iterator {
	out := make(chan interface{})
	go func() {
		for it.Next() {
			x := it.Pos()
			if f(x) {
				out <- x
			}
		}
		close(out)
	}()
	return newChanIterator(out)
}

func (it *Iterator) FoldF(z interface{}, f func(x interface{}, y interface{}) interface{}) *Iterator {
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
