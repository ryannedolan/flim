package future

import (
	"fmt"
	"testing"
  "time"
  "github.com/ryannedolan/flim/fl"
)

func TestFutureResult(t *testing.T) {
	f := New()
	go f.Complete("hello!")
	res, err := f.Wait()
	if res != "hello!" {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}

func TestFutureFailure(t *testing.T) {
	f := New()
	go f.Fail(fmt.Errorf("boo"))
	res, err := f.Wait()
	if res != nil {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

func TestBroadcastPromise(t *testing.T) {
	f := Promise(Success("hello"))
	res, _ := f.Wait()
	res2, _ := f.Wait()
	if res != "hello" {
		t.Fail()
	}
	if res2 != "hello" {
		t.Fail()
	}
}

func TestPromiseLambdaResult(t *testing.T) {
	f := Promise(Success(fl.Lambda(`1 + 2`)(nil)))
	res, _ := f.Wait()
	if res != 3 {
		t.Fatalf("1 + 2 was %v", res)
	}
}

func TestPromiseLambda(t *testing.T) {
	f := Promise(Success(fl.Lambda(`x + y`)))
	res, _ := f.Wait()
	v := res.(func(...interface{}) interface{})(1, 2)
	if v != 3 {
		t.Fatalf("1 + 2 was %v", v)
	}
}

func TestAndThenOrElse(t *testing.T) {
	f1 := Promise(Success("one")).AndThen(Success("two")).OrElse(Success("three"))
	v1, e1 := f1.Wait()
	if v1 != "two" {
		t.Fatalf("expected two, got %v", v1)
	}
	if e1 != nil {
		t.Fatal(e1)
	}

	f2 := Promise(Failuref("one")).AndThen(Success("two")).OrElse(Success("three"))
	v2, e2 := f2.Wait()
	if v2 != "three" {
		t.Fatalf("expected three, got %v", v2)
	}
	if e2 != nil {
		t.Fatal(e2)
	}
}

func TestFutureIter(t *testing.T) {
	a := Promise(Success("foo")).Iter().Map(`x + "bar"`).Array()[0]
	if a != "foobar" {
		t.Fatalf("expected foobar, got %v", a)
	}
}

func TestFutureFromChan(t *testing.T) {
  ch := make(chan interface{}, 1)
  a := New()
  a.CompleteFrom(ch, Timeout(1))
  _, err := a.Wait()
  if err == nil {
    t.Fatalf("expected wait to timeout")
  }

  ch <- "foo"
  b := New()
  b.CompleteFrom(ch, Timeout(1*time.Second))
  s, err2 := b.Wait()
  if err2 != nil {
    t.Fatal(err2)
  }
  if s != "foo" {
    t.Fatalf("expected foo, got %v", s)
  }
}

func TestFirstOf(t *testing.T) {
  done := make(chan struct{}, 1)
  done <- struct{}{}
  ch := FirstOf(Timeout(time.Second), Timeout(time.Millisecond), Timeout(time.Nanosecond), done)
  <-ch 
}
