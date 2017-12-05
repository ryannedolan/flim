package main

import (
  "github.com/ryannedolan/flim/lambda"
  "fmt"
)

func main() {
  f := lambda.X(` x*x `)
  fmt.Println(" f(x) = x*x ")
  fmt.Println(" f(1) =", f(1))
  fmt.Println(" f(2) =", f(2))
  fmt.Println(" f(3) =", f(3))

  g := lambda.XY(` x + y `)
  fmt.Println(" g(x,y) = x + y ")
  fmt.Println(" g(1,2) =", g(1,2))
  fmt.Println(" g(3,4) =", g(3,4))
  fmt.Println(" g(4,5) =", g(4,5))
}