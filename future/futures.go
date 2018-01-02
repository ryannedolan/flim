package future

import (
	"fmt"
	"time"
  "github.com/ryannedolan/flim/fl"
)

type T struct {
	ch  chan struct{}
	err error
	v   interface{}
}

func New() *T {
	return &T{ch: make(chan struct{}, 1)}
}

func Promise(fun func() (interface{}, error)) *T {
	f := New()
	f.CompleteWith(fun)
	return f
}

func Success(a interface{}) func() (interface{}, error) {
	return func() (interface{}, error) { return a, nil }
}

func Failure(err error) func() (interface{}, error) {
	return func() (interface{}, error) { return nil, err }
}

func Failuref(s string, args ...interface{}) func() (interface{}, error) {
	return Failure(fmt.Errorf(s, args...))
}

func Timeout(dt time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
    <-time.After(dt)
    close(ch)
  } ()
	return ch
}

func FirstOf(chs ...<-chan struct{}) <-chan struct{} {
  // we can't do a select on an arbitrary number of chans,
  // but we can complete a future once any one returns
  defer func() {
    recover() // swallow panics from multiple Completes
  }()
  f := New()
  for _, ch := range chs {
    f.CompleteWith(func() (interface{}, error) {
      return <-ch, nil
    })
  }
  _, err := f.Wait()
  if err != nil {
    panic(err)
  }
  done := make(chan struct{})
  close(done)
  return done 
}

func (f *T) finish() {
	close(f.ch) // broadcast that v or err is ready
}

func (f *T) Cancel() {
	f.Fail(fmt.Errorf("cancelled"))
}

func (f *T) Complete(v interface{}) {
	if f.v != nil {
		panic("Future already completed")
	}
	f.v = v
	f.finish()
}

func (f *T) CompleteWith(fun func() (interface{}, error)) {
	go func() {
		res, err := fun()
		if err != nil {
			f.Fail(err)
		} else {
			f.Complete(res)
		}
	}()
}

func (f *T) CompleteFrom(ch <-chan interface{}, done <-chan struct{}) {
  go func () {
    select {
    case res := <-ch:
      f.Complete(res)
    case <-done:
      f.Cancel()
    }
  } ()
}

func (f *T) Fail(err error) {
	f.err = err
	f.finish()
}

func (f *T) Wait() (interface{}, error) {
	for _ = range f.ch { // continues when ch is closed
	}
	return f.v, f.err
}

func (f *T) Error() error {
	return f.err
}

func (f *T) Get() (<-chan interface{}, <-chan error) {
	return f.get()
}

func (f *T) get() (chan interface{}, chan error) {
	v := make(chan interface{}, 1)
	err := make(chan error, 1)
	go func() {
		a, b := f.Wait()
		if a != nil {
			v <- a
		}
		if b != nil {
			err <- b
		}
		close(err)
		close(v)
	}()
	return v, err
}

func (f *T) Iter() *fl.Iterator {
	ch, err := f.get()
	go func() {
		for e := range err {
			panic(e)
		}
	}()
	return fl.Iter(ch)
}

func (f *T) AndThen(fun func() (interface{}, error)) *T {
	f2 := New()
	go func() {
		_, err := f.Wait()
		if err != nil {
			f2.Fail(err)
		} else {
			f2.CompleteWith(fun)
		}
	}()
	return f2
}

func (f *T) OrElse(fun func() (interface{}, error)) *T {
	f2 := New()
	go func() {
		v, err := f.Wait()
		if err != nil {
			f2.CompleteWith(fun)
		} else {
			f2.Complete(v)
		}
	}()
	return f2
}
