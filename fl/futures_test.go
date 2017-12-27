package fl

import (
  "testing"
  "fmt"
)

func TestFutureResult(t *testing.T) {
  f := Future()
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
  f := Future()
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
  f := Promise(Success(Lambda(`1 + 2`)(nil)))
  res, _ := f.Wait()
  if res != 3 {
    t.Fatalf("1 + 2 was %v", res)
  }
}

func TestPromiseLambda(t *testing.T) {
  f := Promise(Success(Lambda(`x + y`)))
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

  b := Promise(Failuref("foo")).Iter().Map(`x + "bar"`).Array()
  if len(b) != 0 {
    t.Fatalf("expected empty list, got %v", b)
  } 
}
