package fl

import (
  "fmt"
  "time"
)

type future struct {
  ch chan struct{}
  err error
  v interface{} 
}

func Future() *future {
  return &future{ch: make(chan struct{}, 1)}
}

func Promise(fun func() (interface{}, error)) *future {
  f := Future()
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
  ch := make(chan struct{}, 1)
  <-time.After(dt)
  ch <- struct{}{}
  return ch
}

func (f *future) finish() {
  close(f.ch) // broadcast that v or err is ready
}

func (f *future) Cancel() {
  f.Fail(fmt.Errorf("cancelled"))
}

func (f *future) Complete(v interface{}) {
  if f.v != nil {
    panic("Future already completed")
  }
  f.v = v
  f.finish()
}

func (f *future) CompleteWith(fun func() (interface{}, error)) {
  go func () {
    res, err := fun()
    if err != nil {
      f.Fail(err)
    } else {
      f.Complete(res)
    }
  } ()
}

func (f *future) Fail(err error) {
  f.err = err
  f.finish()
}

func (f *future) Wait() (interface{}, error) {
  for a := range f.ch { // continues when ch is closed
    _ = a
  }
  return f.v, f.err 
}

func (f *future) Error() error {
  return f.err
}

func (f *future) Get() (<-chan interface{}, <-chan error) {
  return f.get()
}

func (f *future) get() (chan interface{}, chan error) {
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
    close(v)
    close(err)
  } ()
  return v, err
}

func (f *future) Iter() *Iterator {
  ch, err := f.get()
  go func () {
    for e := range err {
      panic(e)
    }
  } ()
  return &Iterator{&chanIterable{ch: ch}}
}

func (f *future) AndThen(fun func() (interface{}, error)) *future {
  f2 := Future()
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

func (f *future) OrElse(fun func() (interface{}, error)) *future {
  f2 := Future()
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
