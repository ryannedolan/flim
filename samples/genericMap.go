package main

import (
  "github.com/ryannedolan/flim/fl"
  "github.com/ryannedolan/flim/lambda"
  "fmt"
)

func main() {
  arr := []int{1, 2, 3, 4, 5}

  a := fl.Iter(arr).Map(`x + 1`).Force()
  fmt.Println(a)

  b := fl.Iter(a).Map(func (x int) int { return x + 1 }).Force()
  fmt.Println(b)

  c := fl.Iter(b).Map(lambda.X(`x + 1`)).Force()
  fmt.Println(c)
}
